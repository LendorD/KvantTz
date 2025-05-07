package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Логирование перед обработкой запроса
		start := time.Now()

		// Обработка запроса
		c.Next()

		// Логирование после обработки
		duration := time.Since(start)
		log.Printf(
			"[REQUEST] %s %s - %d (%v)",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
		)
	}
}
