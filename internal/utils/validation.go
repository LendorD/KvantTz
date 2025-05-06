// internal/utils/validation.go
package utils

import (
	"golang.org/x/crypto/bcrypt"
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
