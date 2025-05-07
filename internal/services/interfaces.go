package services

import (
	"KvantTZ/internal/models"
)

// UserService интерфейс для работы с пользователями
type UserService interface {
	CreateUser(req *models.CreateUserRequest) (*models.UserResponse, error)
	GetAllUsers(page, limit, minAge, maxAge int) ([]models.UserResponse, int64, error)
	GetUserByID(id int) (*models.UserResponse, error)
	UpdateUser(id int, req *models.UpdateUserRequest) (*models.UserResponse, error)
	DeleteUser(id int) error
}

// OrderService интерфейс для работы с заказами
type OrderService interface {
	CreateOrder(userID int, req models.OrderRequest) (*models.OrderResponse, error)
	GetOrders(userID int) ([]models.OrderResponse, error)
}

type AuthService interface {
	Login(email, password string) (string, error)
}
