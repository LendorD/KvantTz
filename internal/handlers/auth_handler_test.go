package handlers_test

import (
	"KvantTZ/internal/handlers"
	"KvantTZ/internal/services/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Login(t *testing.T) {
	mockService := new(mocks.AuthService)
	handler := handlers.NewAuthHandler(mockService)
	router := gin.Default()
	router.POST("/auth/login", handler.Login)

	t.Run("Успешная авторизация", func(t *testing.T) {
		credentials := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(credentials)

		mockService.On("Login", "test@example.com", "password123").Return("valid_token", nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "valid_token")
	})

	t.Run("Неверные учетные данные", func(t *testing.T) {
		credentials := map[string]string{
			"email":    "wrong@example.com",
			"password": "invalid",
		}
		jsonBody, _ := json.Marshal(credentials)

		mockService.On("Login", "wrong@example.com", "invalid").Return("", errors.New("ошибка"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Отсутствует email или пароль", func(t *testing.T) {
		testCases := []struct {
			name string
			body map[string]string
		}{
			{"Нет email", map[string]string{"password": "123"}},
			{"Нет пароля", map[string]string{"email": "test@test.com"}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				jsonBody, _ := json.Marshal(tc.body)

				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))

				router.ServeHTTP(w, req)
				assert.Equal(t, http.StatusBadRequest, w.Code)
			})
		}
	})
}
