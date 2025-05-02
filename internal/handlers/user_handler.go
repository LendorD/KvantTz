package handlers

import (
	"KvantTZ/internal/models"
	"KvantTZ/internal/repository"
	"errors"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	var user models.User

	// Прочитать JSON из запроса
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Проверить, что email уникален
	var existing models.User
	if err := repository.DB.Where("email = ?", user.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь с таким email уже существует"})
		return
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хэшировании пароля"})
		return
	}
	user.PasswordHash = string(hashedPassword)

	// Сохраняем в базу
	if err := repository.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении в базу"})
		return
	}

	// Логирование
	log.Printf("Создан пользователь: %s <%s>\n", user.Name, user.Email)

	// Ответ без пароля
	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"age":   user.Age,
	})
}

func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	minAge, _ := strconv.Atoi(c.Query("min_age"))
	maxAge, _ := strconv.Atoi(c.Query("max_age"))

	offset := (page - 1) * limit
	query := repository.DB.Model(&models.User{})

	if minAge > 0 {
		query = query.Where("age >= ?", minAge)
	}
	if maxAge > 0 {
		query = query.Where("age <= ?", maxAge)
	}

	var total int64
	query.Count(&total)

	var users []models.User
	query.Offset(offset).Limit(limit).Find(&users)

	response := gin.H{
		"page":  page,
		"limit": limit,
		"total": total,
		"users": users,
	}

	c.JSON(http.StatusOK, response)
}

func GetUserByID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID пользователя"})
		return
	}

	var user models.User
	result := repository.DB.First(&user, userID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	response := models.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}

	c.JSON(http.StatusOK, response)
}

func UpdateUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID пользователя"})
		return
	}

	var existingUser models.User
	err = repository.DB.First(&existingUser, userID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	var updateData struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
		Age   int    `json:"age" binding:"required,gte=0,lte=150"`
	}

	err = c.ShouldBindJSON(&updateData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	if updateData.Email != existingUser.Email {
		var confilctUser models.User
		err = repository.DB.Where("email = ?", updateData.Email).First(&confilctUser).Error
		//Если запрос без ошибки, значит email занят
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email уже используется"})
			return
		}
		//Если ошибка НЕ связана с отсутствием записи - возвращаем 500
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
			return
		}
	}

	existingUser.Name = updateData.Name
	existingUser.Email = updateData.Email
	existingUser.Age = updateData.Age

	err = repository.DB.Save(&existingUser).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении пользователя"})
		return
	}

	response := models.User{
		ID:    existingUser.ID,
		Name:  existingUser.Name,
		Email: existingUser.Email,
		Age:   existingUser.Age,
	}
	c.JSON(http.StatusOK, response)
}

func DeleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID пользователя"})
		return
	}

	var user models.User
	result := repository.DB.First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	err = repository.DB.Delete(&models.User{}, userID).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при удалении пользователя"})
		return
	}

	c.Status(http.StatusNoContent)
}
