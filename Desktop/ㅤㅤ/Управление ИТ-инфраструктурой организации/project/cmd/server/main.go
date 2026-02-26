package main

import (
	"log"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lalka1231/mirea-service-desk/internal/handlers"
	"github.com/lalka1231/mirea-service-desk/internal/middleware"
	"github.com/lalka1231/mirea-service-desk/internal/repository"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		Username: getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "mirea_service"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	userRepo := repository.NewUserRepository(db)
	ticketRepo := repository.NewTicketRepository(db)

	authHandler := handlers.NewAuthHandler(userRepo)
	ticketHandler := handlers.NewTicketHandler(ticketRepo)

	router := gin.Default()
	
	// Публичные маршруты
	auth := router.Group("/api")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Защищенные маршруты
	api := router.Group("/api", middleware.AuthMiddleware())
	{
		api.GET("/tickets", ticketHandler.GetAll)
		api.POST("/tickets", ticketHandler.Create)
		api.GET("/tickets/:id", ticketHandler.GetByID)
		api.PUT("/tickets/:id/status", ticketHandler.UpdateStatus)
	}

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	router.Run(":" + port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
