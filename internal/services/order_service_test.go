package services_test

import (
	"KvantTZ/internal/models"
	"KvantTZ/internal/repository/mocks"
	"KvantTZ/internal/services"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestOrderService_CreateOrder(t *testing.T) {
	mockOrderRepo := new(mocks.OrderRepository)
	mockUserRepo := new(mocks.UserRepository)
	service := services.NewOrderService(mockOrderRepo, mockUserRepo, testLogger)

	t.Run("Успешное создание заказа", func(t *testing.T) {
		req := models.OrderRequest{
			Product:  "Laptop",
			Quantity: 1,
			Price:    1200.50,
		}

		mockUserRepo.On("FindByID", 1).Return(&models.User{ID: 1}, nil)
		mockOrderRepo.On("Create", mock.AnythingOfType("*models.Order")).Return(nil)

		response, err := service.CreateOrder(1, req)
		assert.NoError(t, err)
		assert.Equal(t, "Laptop", response.Product)
		mockUserRepo.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("Пользователь не найден", func(t *testing.T) {
		req := models.OrderRequest{Product: "Phone"}

		mockUserRepo.On("FindByID", 999).Return(nil, gorm.ErrRecordNotFound)

		_, err := service.CreateOrder(999, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("Невалидные данные заказа", func(t *testing.T) {
		testCases := []struct {
			name string
			req  models.OrderRequest
		}{
			{"Отрицательное количество", models.OrderRequest{Quantity: -1, Price: 100}},
			{"Нулевая цена", models.OrderRequest{Quantity: 1, Price: 0}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				mockUserRepo.On("FindByID", 1).Return(&models.User{ID: 1}, nil)
				_, err := service.CreateOrder(1, tc.req)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid quantity or price")
			})
		}
	})
}

func TestOrderService_GetOrders(t *testing.T) {
	mockOrderRepo := new(mocks.OrderRepository)
	mockUserRepo := new(mocks.UserRepository)
	service := services.NewOrderService(mockOrderRepo, mockUserRepo, testLogger)

	t.Run("Успешное получение заказов", func(t *testing.T) {
		mockUserRepo.On("FindByID", 1).Return(&models.User{ID: 1}, nil)
		mockOrderRepo.On("GetByUserID", 1).Return([]models.Order{
			{ID: 1, Product: "Laptop", CreatedAt: time.Now()},
		}, nil)

		responses, err := service.GetOrders(1)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(responses))
		assert.Equal(t, "Laptop", responses[0].Product)
	})

	t.Run("Ошибка репозитория", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepository)
		mockOrderRepo := new(mocks.OrderRepository)
		service := services.NewOrderService(mockOrderRepo, mockUserRepo, testLogger)
		
		// Настройка моков
		mockUserRepo.On("FindByID", 1).Return(&models.User{ID: 1}, nil)
		mockOrderRepo.On("GetByUserID", 1).Return([]models.Order{}, errors.New("database error"))

		// Вызов метода
		_, err := service.GetOrders(1)

		// Проверки
		assert.Error(t, err)
		assert.EqualError(t, err, "database error") // Проверяем точное совпадение

		mockUserRepo.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})
}
