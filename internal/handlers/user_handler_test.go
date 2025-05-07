package handlers_test

import (
	"KvantTZ/internal/handlers"
	"KvantTZ/internal/models"
	"KvantTZ/internal/services/mocks"
	"bytes"
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_CreateUser(t *testing.T) {
	mockService := new(mocks.UserService)
	userHandler := handlers.NewUserHandler(mockService)

	t.Run("Успешное создание", func(t *testing.T) {
		reqBody := models.CreateUserRequest{
			Name:     "John",
			Email:    "john@example.com",
			Password: "password",
			Age:      30,
		}
		jsonBody, _ := json.Marshal(reqBody)

		mockService.On("CreateUser", &reqBody).Return(&models.UserResponse{ID: 1}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		userHandler.CreateUser(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetAllUsers(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := handlers.NewUserHandler(mockService)

	t.Run("Успешное получение списка", func(t *testing.T) {
		mockService.On("GetAllUsers", 1, 10, 0, 150).Return(
			[]models.UserResponse{{ID: 1}, {ID: 2}},
			int64(100),
			nil,
		)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/users?page=-1&limit=abc", nil)

		handler.GetAllUsers(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetUserByID(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := handlers.NewUserHandler(mockService)

	t.Run("Успешное получение пользователя", func(t *testing.T) {
		mockService.On("GetUserByID", 1).Return(&models.UserResponse{ID: 1, Name: "John"}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/users/1", nil)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		handler.GetUserByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Неверный ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/users/invalid", nil)
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		handler.GetUserByID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserHandler_UpdateUser(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := handlers.NewUserHandler(mockService)

	t.Run("Успешное обновление", func(t *testing.T) {
		reqBody := models.UpdateUserRequest{
			Name:  "New Name",
			Email: "new@example.com",
			Age:   35,
		}
		jsonBody, _ := json.Marshal(reqBody)

		mockService.On("UpdateUser", 1, &reqBody).Return(&models.UserResponse{ID: 1}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		handler.UpdateUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Неверный формат данных", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PUT", "/users/1", bytes.NewBufferString("{invalid}"))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		handler.UpdateUser(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserHandler_DeleteUser(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := handlers.NewUserHandler(mockService)
	router := gin.Default()
	router.DELETE("/users/:id", handler.DeleteUser)

	t.Run("Успешное удаление", func(t *testing.T) {
		mockService.On("DeleteUser", 1).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/users/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockService.AssertExpectations(t)
	})
	t.Run("Несуществующий пользователь", func(t *testing.T) {
		mockService.On("DeleteUser", 999).Return(gorm.ErrRecordNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("DELETE", "/users/999", nil)
		c.Params = []gin.Param{{Key: "id", Value: "999"}}

		handler.DeleteUser(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
