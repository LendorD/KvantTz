package handlers

import (
	"KvantTZ/internal/models"
	"KvantTZ/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login godoc
// @Summary Аутентификация пользователя
// @Description Выполняет вход пользователя и возвращает JWT-токен
// @Tags auth
// @Accept json
// @Produce json
// @Param  input body models.Credentials true "Данные пользователя"
// @Success 200 {object} map[string]string "{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...}"
// @Failure 400 {object} map[string]string "{"error": "Email и пароль обязательны"}"
// @Failure 401 {object} map[string]string "{"error": "Неверные учетные данные"}"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var credentials models.Credentials
	if err := c.ShouldBindJSON(&credentials); err != nil || credentials.Email == "" || credentials.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email и пароль обязательны"})
		return
	}

	token, err := h.authService.Login(credentials.Email, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
