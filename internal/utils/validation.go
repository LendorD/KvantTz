// internal/utils/validation.go
package utils

import (
	"KvantTZ/internal/models"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strconv"
)

// Хеширование пароля
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Проверка пароля
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Валидация возраста
func ValidateAge(age int, min, max int) bool {
	return age >= min && age <= max
}

func IsValidEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ValidateCreateRequest(req *models.CreateUserRequest) error {
	if !IsValidEmail(req.Email) {
		return errors.New("invalid email")
	}
	if !ValidateAge(req.Age, 0, 150) {
		return errors.New("invalid age")
	}
	if req.Email == "" || req.Name == "" || req.Password == "" {
		return errors.New("email, name, and password cannot be empty")
	}
	return nil
}

func ValidateListUsersParams(c *gin.Context) (page, limit, minAge, maxAge int) {
	page = 1
	limit = 10
	minAge = 0
	maxAge = 150

	if pageStr := c.Query("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 100 {
			limit = 10
		}
	}

	if minAgeStr := c.Query("minAge"); minAgeStr != "" {
		minAge, err := strconv.Atoi(minAgeStr)
		if err != nil || minAge < 0 {
			minAge = 0
		}
	}

	if maxAgeStr := c.Query("maxAge"); maxAgeStr != "" {
		maxAge, err := strconv.Atoi(maxAgeStr)
		if err != nil || maxAge < 0 {
			maxAge = 150
		}
	}

	if minAge > maxAge {
		minAge = 0
		maxAge = 100
	}

	return page, limit, minAge, maxAge
}
