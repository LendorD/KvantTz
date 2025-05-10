package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

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
