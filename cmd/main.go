package main

import (
	"KvantTZ/internal/handlers"
	"KvantTZ/internal/middleware"
	"KvantTZ/internal/repository"
	"KvantTZ/internal/services"
	"KvantTZ/migrations"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	db, err := repository.InitDB()
	if err != nil {
		log.Fatal(err.Error())
	}
	//err := runMigrations()
	err = migrations.Migrate(db)
	if err != nil {
		log.Fatal("Ошибка миграций: ", err.Error())
	}
	err = godotenv.Load()

	// Репозитории
	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Сервисы
	userService := services.NewUserService(userRepo)
	orderService := services.NewOrderService(orderRepo, userRepo)
	authService := services.NewAuthService(userRepo)

	// Хендлеры
	userHandler := handlers.NewUserHandler(userService)
	orderHandler := handlers.NewOrderHandler(orderService)
	authHandler := handlers.NewAuthHandler(authService)

	router := gin.Default()

	// Public routes
	router.POST("/auth/login", authHandler.Login)
	router.POST("/users", userHandler.CreateUser)

	// Protected routes
	authGroup := router.Group("/")
	authGroup.Use(middleware.JWTAuth())
	{

		authGroup.GET("/users", userHandler.GetAllUsers)
		authGroup.GET("/users/:id", userHandler.GetUserByID)
		authGroup.PUT("/users/:id", userHandler.UpdateUser)
		authGroup.DELETE("/users/:id", userHandler.DeleteUser)

		authGroup.POST("/users/:id/orders", orderHandler.CreateOrder)
		authGroup.GET("/users/:id/orders", orderHandler.GetOrders)

	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}
