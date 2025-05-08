package handlers

import (
	"KvantTZ/internal/models"
	"KvantTZ/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder godoc
// @Summary Создать заказ
// @Description Создает новый заказ для указанного пользователя
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param input body models.OrderRequest true "Данные заказа"
// @Success 201 {object} models.OrderResponse
// @Failure 400 {object} map[string]string "Неверный ID пользователя или данные заказа"
// @Failure 404 {object} map[string]string "Пользователь не найден"
// @Router /users/{id}/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var orderRequest models.OrderRequest
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order request"})
		return
	}

	order, err := h.orderService.CreateOrder(userID, orderRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetOrders godoc
// @Summary Получить заказы пользователя
// @Description Возвращает список всех заказов указанного пользователя
// @Tags orders
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {array} models.OrderResponse
// @Failure 400 {object} map[string]string "Неверный ID пользователя"
// @Failure 404 {object} map[string]string "Пользователь не найден"
// @Router /users/{id}/orders [get]
func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	orders, err := h.orderService.GetOrders(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}
