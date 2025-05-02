package main

import (
	"KvantTZ/internal/handlers"
	"KvantTZ/internal/middleware"
	"KvantTZ/internal/repository"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func main() {
	repository.InitDB()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Ошибка загрузки ENV")
	}
	err = runMigrations()
	if err != nil {
		log.Fatal("Ошибка миграций: ", err)
	}

	router := gin.Default()

	// Public routes
	router.POST("/auth/login", handlers.Login)
	router.POST("/users", handlers.CreateUser)

	// Protected routes
	authGroup := router.Group("/")
	authGroup.Use(middleware.JWTAuth())
	{

		authGroup.GET("/users", handlers.GetUsers)
		authGroup.GET("/users/:id", handlers.GetUserByID)
		authGroup.PUT("/users/:id", handlers.UpdateUser)
		authGroup.DELETE("/users/:id", handlers.DeleteUser)

		authGroup.POST("/users/:id/orders", handlers.CreateOrder)
		authGroup.GET("/users/:id/orders", handlers.GetOrders)

	}

	router.Run(":8080")
}

func runMigrations() error {
	// Команда для применения миграций через psql
	cmd := exec.Command("psql",
		"-h", os.Getenv("DB_HOST"),
		"-U", os.Getenv("DB_USER"),
		"-d", os.Getenv("DB_NAME"),
		"-f", "/app/migrations/001_create_tables.sql",
	)

	// Установка пароля
	cmd.Env = append(os.Environ(), "PGPASSWORD="+os.Getenv("DB_PASSWORD"))

	// Запуск
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Println("Вывод команды: ", string(output))
		return err
	}

	return nil
}
