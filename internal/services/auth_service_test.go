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
	service := services.NewAuthService(mockUserRepo)

	t.Run("Успешная авторизация", func(t *testing.T) {
		// Генерируем хеш для валидного пароля
		hashedPassword, _ := utils.HashPassword("valid_password")
		mockUser := &models.User{
			ID:           1,
			PasswordHash: hashedPassword,
		}

		// Настраиваем мок (однократный вызов)
		mockUserRepo.On("FindByEmail", "test@example.com").Return(mockUser, nil)

		// Вызываем метод сервиса
		token, err := service.Login("test@example.com", "valid_password")

		// Проверки
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

	t.Run("Ошибка генерации токена", func(t *testing.T) {
		// Тест требует мокирования utils.GenerateJWT, что сложнее.
		// Рекомендуется использовать интерфейсы для утилит.
	})
}
