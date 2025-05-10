package services_test

import (
	"KvantTZ/internal/models"
	"KvantTZ/internal/repository/mocks"
	"KvantTZ/internal/services"
	"errors"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var testLogger = log.New(io.Discard, "", 0)

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	service := services.NewUserService(mockRepo, testLogger)

	t.Run("Успешное создание пользователя", func(t *testing.T) {
		req := &models.CreateUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Age:      30,
			Password: "password123",
		}

		mockRepo.On("FindByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound)
		mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

		user, err := service.CreateUser(req)
		assert.NoError(t, err)
		assert.Equal(t, req.Name, user.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка: email уже существует", func(t *testing.T) {
		req := &models.CreateUserRequest{
			Name:     "Existing User",
			Email:    "existing@example.com",
			Password: "password123",
			Age:      25,
		}

		mockRepo.On("FindByEmail", req.Email).Return(&models.User{}, nil)

		_, err := service.CreateUser(req)
		assert.Error(t, err)
		assert.Equal(t, "user already exists", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetAllUsers(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	service := services.NewUserService(mockRepo, testLogger)

	t.Run("Успешное получение списка", func(t *testing.T) {
		mockRepo.On("GetAll", 0, 10, 0, 0).Return([]models.User{
			{ID: 1, Name: "User1"},
			{ID: 2, Name: "User2"},
		}, int64(2), nil)

		users, total, err := service.GetAllUsers(1, 10, 0, 0)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(users))
		assert.Equal(t, int64(2), total)
	})
}

func TestUserService_GetUserByID(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	service := services.NewUserService(mockRepo, testLogger)

	t.Run("Успешное получение пользователя", func(t *testing.T) {
		mockRepo.On("FindByID", 1).Return(&models.User{
			ID:    1,
			Name:  "John",
			Email: "john@example.com",
			Age:   30,
		}, nil)

		user, err := service.GetUserByID(1)
		assert.NoError(t, err)
		assert.Equal(t, "John", user.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Пользователь не найден", func(t *testing.T) {
		mockRepo.On("FindByID", 999).Return(nil, gorm.ErrRecordNotFound)

		_, err := service.GetUserByID(999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	service := services.NewUserService(mockRepo, testLogger)

	t.Run("Успешное обновление", func(t *testing.T) {
		mockRepo.On("FindByID", 1).Return(&models.User{
			ID:    1,
			Email: "old@example.com",
		}, nil)
		mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil)

		req := &models.UpdateUserRequest{
			Name:  "New Name",
			Email: "new@example.com",
			Age:   35,
		}

		_, err := service.UpdateUser(1, req)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка валидации email", func(t *testing.T) {
		req := &models.UpdateUserRequest{Email: "invalid-email"}
		_, err := service.UpdateUser(1, req)
		assert.Error(t, err)
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	service := services.NewUserService(mockRepo, testLogger)

	t.Run("Успешное удаление", func(t *testing.T) {
		mockRepo.On("Delete", 1).Return(nil)
		err := service.DeleteUser(1)
		assert.NoError(t, err)
	})

	t.Run("Ошибка при удалении", func(t *testing.T) {
		mockRepo.On("Delete", 2).Return(errors.New("database error"))
		err := service.DeleteUser(2)
		assert.Error(t, err)
	})
}

func TestUserService_EdgeCases(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	service := services.NewUserService(mockRepo, testLogger)

	t.Run("Некорректная пагинация", func(t *testing.T) {
		_, _, err := service.GetAllUsers(-1, 1000, -5, 200)
		assert.Error(t, err)
	})

	t.Run("Фильтрация по возрасту", func(t *testing.T) {
		mockRepo.On("GetAll", 0, 10, 20, 30).Return([]models.User{}, int64(0), nil)
		_, _, err := service.GetAllUsers(1, 10, 20, 30)
		assert.NoError(t, err)
	})
}
