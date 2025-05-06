package repository

import "KvantTZ/internal/models"

type UserRepository interface {
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id int) error
	GetAll(offset, limit, minAge, maxAge int) ([]models.User, int64, error)
	FindByEmail(email string) (*models.User, error)
	FindByID(id int) (*models.User, error)
}

type OrderRepository interface {
	Create(order *models.Order) error
	GetByUserID(userID int) ([]models.Order, error)
}
