package mocks

import (
	"KvantTZ/internal/models"
	"github.com/stretchr/testify/mock"
)

type OrderService struct {
	mock.Mock
}

func (m *OrderService) CreateOrder(userID int, req models.OrderRequest) (*models.OrderResponse, error) {
	args := m.Called(userID, req)
	return args.Get(0).(*models.OrderResponse), args.Error(1)
}

func (m *OrderService) GetOrders(userID int) ([]models.OrderResponse, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.OrderResponse), args.Error(1)
}
