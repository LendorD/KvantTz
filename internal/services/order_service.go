package services

import (
	"KvantTZ/internal/models"
	"KvantTZ/internal/repository"
	"errors"
	"log"
	"time"
)

type orderService struct {
	orderRepo repository.OrderRepository
	userRepo  repository.UserRepository
	logger    *log.Logger
}

func NewOrderService(orderRepo repository.OrderRepository, userRepo repository.UserRepository, logger *log.Logger) OrderService {
	return &orderService{
		orderRepo: orderRepo,
		userRepo:  userRepo,
		logger:    logger,
	}
}

func (s *orderService) CreateOrder(userID int, req models.OrderRequest) (*models.OrderResponse, error) {
	// Проверяем, существует ли пользователь
	s.logger.Printf("[INFO] Создание заказа у пользователя с ID: %d", userID)
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		s.logger.Println("[ERROR] Ошибка при создании заказа, пользователь не найден")
		return nil, errors.New("user not found")
	}

	// Валидация
	if req.Quantity <= 0 || req.Price <= 0 {
		s.logger.Println("[ERROR] Ошибка параметров. Цена или колечество неверны")
		return nil, errors.New("invalid quantity or price")
	}

	order := models.Order{
		UserID:   uint(userID),
		Product:  req.Product,
		Quantity: req.Quantity,
		Price:    req.Price,
	}

	if err := s.orderRepo.Create(&order); err != nil {
		s.logger.Printf("[ERROR] Ошибка при создании заказа: %v", err)
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
	s.logger.Printf("[INFO] Создан заказ у пользователся с ID: %d", userID)
	return response, nil
}

func (s *orderService) GetOrders(userID int) ([]models.OrderResponse, error) {
	s.logger.Printf("[INFO] Получение заказа у пользователя с ID: %d", userID)
	// Проверяем, существует ли пользователь
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		s.logger.Println("[ERROR] Ошибка при получении заказа. Пользователь не найден")
		return nil, errors.New("user not found")
	}

	orders, err := s.orderRepo.GetByUserID(userID)
	if err != nil {
		s.logger.Printf("[ERROR] Ошибка при получении пользователся по ID: %v", err)
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
	s.logger.Printf("[INFO] Список заказов у пользователя с ID: %d возвращён", userID)
	return responses, nil
}

var _ OrderService = (*orderService)(nil)
