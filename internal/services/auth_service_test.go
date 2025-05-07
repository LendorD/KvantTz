package services_test

import (
	"KvantTZ/internal/models"
	"KvantTZ/internal/repository/mocks"
	"KvantTZ/internal/services"
	"KvantTZ/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAuthService_Login(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	service := services.NewAuthService(mockUserRepo, testLogger)

	t.Run("Успешная авторизация", func(t *testing.T) {
		hashedPassword, _ := utils.HashPassword("valid_password")
		mockUser := &models.User{
			ID:           1,
			PasswordHash: hashedPassword,
		}

		mockUserRepo.On("FindByEmail", "test@example.com").Return(mockUser, nil)

		token, err := service.Login("test@example.com", "valid_password")

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Пользователь не найден", func(t *testing.T) {
		mockUserRepo.On("FindByEmail", "invalid@example.com").Return(nil, gorm.ErrRecordNotFound)

		_, err := service.Login("invalid@example.com", "any_password")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("Неверный пароль", func(t *testing.T) {
		mockUser := &models.User{
			ID:           1,
			PasswordHash: "$2a$10$fakehash",
		}

		mockUserRepo.On("FindByEmail", "test@example.com").Return(mockUser, nil)

		_, err := service.Login("test@example.com", "wrong_password")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid password")
	})

}
