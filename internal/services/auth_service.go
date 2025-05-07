package services

import (
	"KvantTZ/internal/repository"
	"KvantTZ/internal/utils"
	"errors"
	"gorm.io/gorm"
)

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", errors.New("user not found")
	}
	if err != nil {
		return "", errors.New("invalid email")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid password")
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", errors.New("error generating token")
	}

	return token, nil
}

var _ AuthService = (*authService)(nil)
