package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func GenerateJWT(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateJWT(tokenString string) (uint, error) {
	if tokenString == "" {
		return 0, errors.New("токен отсутствует")
	}

	// Парсим токен с секретным ключом
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("невалидный токен")
	}

	// Проверяем claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Проверяем срок действия
		exp, err := claims.GetExpirationTime()
		if err != nil || exp.Before(time.Now()) {
			return 0, errors.New("токен истек")
		}

		// Извлекаем user_id
		userID, ok := claims["user_id"]
		if !ok {
			return 0, errors.New("токен не содержит user_id")
		}

		// Конвертируем в uint
		switch v := userID.(type) {
		case float64:
			return uint(v), nil
		case int:
			return uint(v), nil
		default:
			return 0, fmt.Errorf("некорректный тип user_id: %T", userID)
		}
	}

	return 0, errors.New("не удалось распарсить claims")
}
