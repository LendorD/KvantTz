package handlers

import (
	"KvantTZ/internal/models"
	"KvantTZ/internal/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	// Получаем user_id из URL
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID пользователя"})
		return
	}

	// Проверяем существование пользователя
	var user models.User
	err = repository.DB.First(&user, userID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	var orderRequest models.OrderRequest
	err = c.ShouldBindJSON(&orderRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Валидация данных
	if orderRequest.Quantity <= 0 || orderRequest.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Количество и цена должны быть положительными"})
		return
	}

	// Создаем заказ
	order := models.Order{
		UserID:   uint(userID),
		Product:  orderRequest.Product,
		Quantity: orderRequest.Quantity,
		Price:    orderRequest.Price,
	}

	// Сохраняем в базу
	if err := repository.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании заказа"})
		return
	}

	// Формируем ответ
	response := models.OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		Product:   order.Product,
		Quantity:  order.Quantity,
		Price:     order.Price,
		CreatedAt: order.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, response)
}

func GetOrders(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID пользователя"})
		return
	}

	// Проверяем существование пользователя
	var user models.User
	err = repository.DB.First(&user, userID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	// Получаем заказы пользователя
	var orders []models.Order
	err = repository.DB.Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении заказов"})
		return
	}

	// Формируем ответ
	var response []models.OrderResponse
	for _, order := range orders {
		response = append(response, models.OrderResponse{
			ID:        order.ID,
			UserID:    order.UserID,
			Product:   order.Product,
			Quantity:  order.Quantity,
			Price:     order.Price,
			CreatedAt: order.CreatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, response)
}
