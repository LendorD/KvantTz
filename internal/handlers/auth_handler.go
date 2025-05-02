package handlers

import (
	"KvantTZ/internal/models"
	"KvantTZ/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

func Login(c *gin.Context) {
	// Структура для приёма учётных данных из запроса
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// Парсинг JSON тела запроса
	err := c.ShouldBindJSON(&credentials)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}
	if credentials.Email == "" || credentials.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email и пароль обязательны"})
		return
	}
	var user models.User
	// Поиск пользователя в базе по email

	err = repository.DB.Where("email = ?", credentials.Email).First(&user).Error
	if err != nil {
		// Если пользователь не найден, возвращаем 401 Unauthorized
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные учетные данные"})
		return
	}
	// Сравнение хэша пароля из БД с введённым паролем
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные учетные данные"})
		return
	}
	// Создание JWT токена      Алгоритм подписи       Полезная нагрузка токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	// Подписание токена секретным ключом из переменных окружения
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
