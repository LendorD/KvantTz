package repository

import (
	"KvantTZ/internal/models"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}

}

func (r orderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r orderRepository) GetByUserID(userID int) ([]models.Order, error) {
	var orders []models.Order
	// SELECT * FROM orders WHERE user_id = ?
	err := r.db.Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}
