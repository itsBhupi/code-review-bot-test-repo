package main

import (
	"fmt"
	"log"
	"os"

	"code-review-bot-test-repo/controllers"
	"code-review-bot-test-repo/services"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	dsn := "host=localhost user=postgres password=postgres dbname=test_db port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Database connected successfully")
}

func setupRouter() *gin.Engine {
	// Initialize services
	helloService := services.NewHelloService()

	// Initialize controllers
	helloController := controllers.NewHelloController(helloService)

	// Create Gin router
	r := gin.Default()

	// Set up routes
	helloController.RegisterRoutes(r)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	return r
}

func main() {
	// Initialize database
	initDB()

	// Set Gin mode
	if os.Getenv("GIN_MODE") != "release" {
		gin.SetMode(gin.DebugMode)
	}

	// Setup router
	r := setupRouter()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
