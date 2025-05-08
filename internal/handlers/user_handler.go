package handlers

import (
	_ "KvantTZ/docs"
	"KvantTZ/internal/models"
	"KvantTZ/internal/services"
	"KvantTZ/internal/utils"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser godoc
// @Summary      Создать пользователя
// @Description  Создает нового пользователя
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        input body models.CreateUserRequest true "Данные пользователя"
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.CreateUserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := h.userService.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

// GetAllUsers godoc
// @Summary      Получить список пользователей
// @Description  Возвращает список пользователей с пагинацией и фильтрацией по возрасту
// @Tags         users
// @Produce      json
// @Param        page     query int false "Номер страницы (по умолчанию 1)"
// @Param        limit    query int false "Лимит элементов на странице (по умолчанию 10)"
// @Param        min_age  query int false "Минимальный возраст"
// @Param        max_age  query int false "Максимальный возраст"
// @Success      200 {object} map[string]interface{} "{"users": [], "total": 0, "page": 0, "limit": 0}"
// @Failure      500 {object} map[string]string "{"error": "service error"}"
// @Router       /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	page, limit, minAge, maxAge := utils.ValidateListUsersParams(c)
	users, total, err := h.userService.GetAllUsers(page, limit, minAge, maxAge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetUserByID godoc
// @Summary      Получить пользователя по ID
// @Description  Возвращает данные пользователя по указанному ID
// @Tags         users
// @Produce      json
// @Param        id path int true "ID пользователя"
// @Success      200 {object} models.User
// @Failure      400 {object} map[string]string "{"error": "invalid user ID"}"
// @Failure      404 {object} map[string]string "{"error": "user not found"}"
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary      Обновить данные пользователя
// @Description  Обновляет данные пользователя по указанному ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path int true "ID пользователя"
// @Param        input body models.UpdateUserRequest true "Новые данные пользователя"
// @Success      200 {object} models.User
// @Failure      400 {object} map[string]string "{"error": "invalid user ID"}"
// @Failure      404 {object} map[string]string "{"error": "user not found"}"
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.userService.UpdateUser(id, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser godoc
// @Summary      Удалить пользователя
// @Description  Удаляет пользователя по указанному ID
// @Tags         users
// @Param        id path int true "ID пользователя"
// @Success      204 "No Content"
// @Failure      400 {object} map[string]string "{"error": "invalid user ID"}"
// @Failure      404 {object} map[string]string "{"error": "user not found"}"
// @Failure      500 {object} map[string]string "{"error": "internal error message"}"
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	err = h.userService.DeleteUser(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
