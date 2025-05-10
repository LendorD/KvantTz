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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

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
	// Значения по умолчанию
	page = 1
	limit = 10
	minAge = 0
	maxAge = 150

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if minAgeStr := c.Query("min_age"); minAgeStr != "" {
		if ma, err := strconv.Atoi(minAgeStr); err == nil && ma >= 0 {
			minAge = ma
		}
	}

	if maxAgeStr := c.Query("max_age"); maxAgeStr != "" {
		if ma, err := strconv.Atoi(maxAgeStr); err == nil && ma >= 0 {
			maxAge = ma
		}
	}

	if minAge > maxAge {
		minAge, maxAge = 0, 150
	}

	return
}
