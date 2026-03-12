package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
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

	// Используем SQLite вместо PostgreSQL
	dbPath := getEnv("DB_PATH", "./mirea_service.db")
	db, err := repository.NewSQLiteDB(dbPath)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	userRepo := repository.NewUserRepository(db)
	ticketRepo := repository.NewTicketRepository(db)

	authHandler := handlers.NewAuthHandler(userRepo)
	ticketHandler := handlers.NewTicketHandler(ticketRepo)

	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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
