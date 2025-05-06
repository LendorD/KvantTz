package services

import (
	"KvantTZ/internal/models"
	"KvantTZ/internal/repository"
	"errors"
	"time"
)

type OrderService struct {
	orderRepo repository.OrderRepository
	userRepo  repository.UserRepository
}

func NewOrderService(orderRepo repository.OrderRepository, userRepo repository.UserRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		userRepo:  userRepo,
	}
}

func (s *OrderService) CreateOrder(userID int, req models.OrderRequest) (*models.OrderResponse, error) {
	// Проверяем, существует ли пользователь
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Валидация
	if req.Quantity <= 0 || req.Price <= 0 {
		return nil, errors.New("invalid quantity or price")
	}

	order := models.Order{
		UserID:   uint(userID),
		Product:  req.Product,
		Quantity: req.Quantity,
		Price:    req.Price,
	}

	if err := s.orderRepo.Create(&order); err != nil {
		return nil, err
	}

	response := &models.OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		Product:   order.Product,
		Quantity:  order.Quantity,
		Price:     order.Price,
		CreatedAt: order.CreatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (s *OrderService) GetOrders(userID int) ([]models.OrderResponse, error) {
	// Проверяем, существует ли пользователь
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	orders, err := s.orderRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []models.OrderResponse
	for _, order := range orders {
		responses = append(responses, models.OrderResponse{
			ID:        order.ID,
			UserID:    order.UserID,
			Product:   order.Product,
			Quantity:  order.Quantity,
			Price:     order.Price,
			CreatedAt: order.CreatedAt.Format(time.RFC3339),
		})
	}

	return responses, nil
}
