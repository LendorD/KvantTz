package mocks

import (
	"KvantTZ/internal/models"
	"github.com/stretchr/testify/mock"
)

type OrderRepository struct {
	mock.Mock
}

func (o OrderRepository) Create(order *models.Order) error {
	args := o.Called(order)
	return args.Error(0)
}

func (o OrderRepository) GetByUserID(userID int) ([]models.Order, error) {
	args := o.Called(userID)
	return args.Get(0).([]models.Order),
		args.Error(1)
}
