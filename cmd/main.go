package main

import (
	_ "KvantTZ/docs"
	"KvantTZ/internal/handlers"
	"KvantTZ/internal/middleware"
	"KvantTZ/internal/repository"
	"KvantTZ/internal/services"
	"KvantTZ/migrations"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
)

// @title KvantTZ API
// @version 1.0
// @description API для управления пользователями и заказами
// @contact.name API Support
// @contact.email support@kvanttz.ru
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("[INFO] Не найден .env файл, используются системные переменные")
	}

	db, err := repository.InitDB()
	if err != nil {
		log.Fatal("[ERROR] Ошибка инициализации БД: ", err)
	}

	if err := migrations.Migrate(db); err != nil {
		log.Fatal("[ERROR] Ошибка миграций: ", err)
	}
	logger := log.New(os.Stdout, "", log.LstdFlags)
	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	userService := services.NewUserService(userRepo, logger)
	orderService := services.NewOrderService(orderRepo, userRepo, logger)
	authService := services.NewAuthService(userRepo, logger)

	router := gin.Default()
	router.Use(middleware.RequestLogger())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	setupPublicRoutes(router, authService, userService)

	setupProtectedRoutes(router, userService, orderService)

	port := getPort()
	log.Printf("[INFO] Сервер запущен на порту :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("[ERROR] Ошибка запуска сервера: ", err)
	}
}

func getPort() string {
	if port := os.Getenv("APP_PORT"); port != "" {
		return port
	}
	return "8080"
}

func setupPublicRoutes(router *gin.Engine, authService services.AuthService, userService services.UserService) {
	router.POST("/auth/login", handlers.NewAuthHandler(authService).Login)
	router.POST("/users", handlers.NewUserHandler(userService).CreateUser)
}

func setupProtectedRoutes(router *gin.Engine, userService services.UserService, orderService services.OrderService) {
	authGroup := router.Group("/")
	authGroup.Use(middleware.JWTAuth())
	{
		// Пользователи
		authGroup.GET("/users", handlers.NewUserHandler(userService).GetAllUsers)
		authGroup.GET("/users/:id", handlers.NewUserHandler(userService).GetUserByID)
		authGroup.PUT("/users/:id", handlers.NewUserHandler(userService).UpdateUser)
		authGroup.DELETE("/users/:id", handlers.NewUserHandler(userService).DeleteUser)

		// Заказы
		authGroup.POST("/users/:id/orders", handlers.NewOrderHandler(orderService).CreateOrder)
		authGroup.GET("/users/:id/orders", handlers.NewOrderHandler(orderService).GetOrders)
	}
}
