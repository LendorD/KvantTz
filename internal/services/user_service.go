package services

import (
	"KvantTZ/internal/utils"
	"errors"

	"KvantTZ/internal/models"
	"KvantTZ/internal/repository"

	"gorm.io/gorm"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// Создание нового пользователя
func (s *UserService) CreateUser(req *models.CreateUserRequest) (*models.UserResponse, error) {
	// Проверка уникальности email
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if existingUser == nil {
		return nil, errors.New("email already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if !utils.ValidateAge(req.Age, 0, 150) {
		return nil, errors.New("invalid age")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("error hashing password")
	}

	// Создание объекта User из запроса
	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		Age:          req.Age,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("ошибка при создании пользователя")
	}

	return &models.UserResponse{
		ID:    user.ID, // ID автоматически генерируется при сохранении
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}, nil
}

// Получение всех пользователей с пагинацией и фильтрацией
func (s *UserService) GetAllUsers(page, limit, minAge, maxAge int) ([]models.UserResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	if minAge < 0 || maxAge < 0 || minAge > maxAge {
		return nil, 0, errors.New("incorrect data parameters")
	}

	offset := (page - 1) * limit
	users, total, err := s.userRepo.GetAll(offset, limit, minAge, maxAge)
	if err != nil {
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

// Получение пользователя по ID
func (s *UserService) GetUserByID(id int) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}, nil
}

// Обновление пользователя
func (s *UserService) UpdateUser(id int, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Age = req.Age

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}, nil
}

// Удаление пользователя
func (s *UserService) DeleteUser(id int) error {
	return s.userRepo.Delete(id)
}
