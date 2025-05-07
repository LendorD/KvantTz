package handlers_test

import (
	"KvantTZ/internal/handlers"
	"KvantTZ/internal/models"
	"KvantTZ/internal/services/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrderHandler_CreateOrder(t *testing.T) {
	mockService := new(mocks.OrderService)
	handler := handlers.NewOrderHandler(mockService)
	router := gin.Default()
	router.POST("/users/:id/orders", handler.CreateOrder)

	t.Run("Успешное создание заказа", func(t *testing.T) {
		orderRequest := models.OrderRequest{
			Product:  "Laptop",
			Quantity: 1,
			Price:    1200.50,
		}
		jsonBody, _ := json.Marshal(orderRequest)

		mockService.On("CreateOrder", 1, orderRequest).Return(&models.OrderResponse{
			ID:       1,
			UserID:   1,
			Product:  "Laptop",
			Quantity: 1,
			Price:    1200.50,
		}, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users/1/orders", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Неверный user_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users/invalid/orders", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Ошибка сервиса", func(t *testing.T) {
		orderRequest := models.OrderRequest{Product: "Phone"}
		jsonBody, _ := json.Marshal(orderRequest)

		mockService.On("CreateOrder", 2, mock.Anything).Return((*models.Order)(nil), errors.New("ошибка"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users/2/orders", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestOrderHandler_GetOrders(t *testing.T) {
	mockService := new(mocks.OrderService)
	handler := handlers.NewOrderHandler(mockService)
	router := gin.Default()
	router.GET("/users/:id/orders", handler.GetOrders)

	t.Run("Успешное получение заказов", func(t *testing.T) {
		mockService.On("GetOrders", 1).Return([]models.OrderResponse{
			{ID: 1, Product: "Laptop"},
			{ID: 2, Product: "Phone"},
		}, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/1/orders", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Пользователь без заказов", func(t *testing.T) {
		mockService.On("GetOrders", 2).Return([]models.OrderResponse{}, errors.New("нет заказов"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/2/orders", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
