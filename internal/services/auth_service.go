package services

import (
	"KvantTZ/internal/repository"
	"KvantTZ/internal/utils"
	"errors"
	"gorm.io/gorm"
	"log"
)

type authService struct {
	userRepo repository.UserRepository
	logger   *log.Logger
}

func NewAuthService(userRepo repository.UserRepository, logger *log.Logger) AuthService {
	return &authService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *authService) Login(email, password string) (string, error) {
	s.logger.Printf("[INFO] Авторизация пользователя: %s", email)
	user, err := s.userRepo.FindByEmail(email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Println("[ERROR] Пользователь не найден")
		return "", errors.New("user not found")
	}
	if err != nil {
		s.logger.Println("[ERROR] Некорректный email")
		return "", errors.New("invalid email")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		s.logger.Println("[ERROR] Некорректный пароль")
		return "", errors.New("invalid password")
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		s.logger.Println("[ERROR] Ошибка при генерации JWT")
		return "", errors.New("error generating token")
	}

	return token, nil
}

var _ AuthService = (*authService)(nil)
