package services

import (
	"KvantTZ/internal/models"
	"KvantTZ/internal/repository"
	"KvantTZ/internal/utils"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type userService struct {
	userRepo repository.UserRepository
	logger   *log.Logger
}

func NewUserService(userRepo repository.UserRepository, logger *log.Logger) UserService {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *userService) CreateUser(req *models.CreateUserRequest) (*models.UserResponse, error) {
	s.logger.Printf("[INFO] Создание пользователя: %s", req.Email)
	err := utils.ValidateCreateRequest(req)
	if err != nil {
		s.logger.Printf("[ERROR] Ошибка параметров: %v", err)
		return nil, errors.New(err.Error())
	}
	_, err = s.userRepo.FindByEmail(req.Email)
	if err == nil {
		s.logger.Printf("[ERROR] Ошибка поиска пользователя: %v", err)
		return nil, errors.New("user already exists")
	}
	if err == nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		Age:          req.Age,
		PasswordHash: hashedPassword,
	}
	err = s.userRepo.Create(user)
	if err != nil {
		s.logger.Printf("[ERROR] Ошибка создания пользователя: %v", err)
		return nil, fmt.Errorf("error creating user: %w", err)
	}
	s.logger.Printf("[INFO] Пользователь создан : %s", user.Email)

	return &models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}, nil
}

func (s *userService) GetAllUsers(page, limit, minAge, maxAge int) ([]models.UserResponse, int64, error) {
	s.logger.Println("[INFO] Получение всех пользователей")
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	if minAge < 0 || maxAge < 0 || minAge > maxAge {
		s.logger.Printf("[ERROR] Ошибка параметров при получении пользователй")
		return nil, 0, errors.New("incorrect data parameters")
	}

	offset := (page - 1) * limit
	users, total, err := s.userRepo.GetAll(offset, limit, minAge, maxAge)
	if err != nil {
		s.logger.Printf("[ERROR] Ошибка при получении пользователей: %v", err)
		return nil, 0, err
	}

	var response = make([]models.UserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, models.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
		})
	}
	return response, total, nil
}

func (s *userService) GetUserByID(id int) (*models.UserResponse, error) {
	s.logger.Printf("[INFO] Получение пользователя по ID: %d ", id)

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		s.logger.Printf("[ERROR] Ошибка при получении пользователся по ID: %v ", err)
		return nil, err
	}
	return &models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}, nil
}

func (s *userService) UpdateUser(id int, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	s.logger.Printf("[INFO] Обновление данных пользователя c ID: %d ", id)
	if !utils.IsValidEmail(req.Email) {
		s.logger.Println("[ERROR] Ошибка при валидации email")
		return nil, errors.New("invalid email")
	}
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		s.logger.Printf("[ERROR] Пользователь не найден: %v ", err)
		return nil, err
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Age = req.Age

	if err := s.userRepo.Update(user); err != nil {
		s.logger.Printf("[ERROR] Не удалость обновить данные пользователя: %v ", err)
		return nil, err
	}

	return &models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}, nil
}

func (s *userService) DeleteUser(id int) error {
	s.logger.Printf("[INFO] Удаление пользователя с ID: %d ", id)
	err := s.userRepo.Delete(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Println("[ERROR] Не удалось найти и удалить пользователя")
		return err
	}
	return err
}

var _ UserService = (*userService)(nil)
