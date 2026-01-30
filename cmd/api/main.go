package main

import (
	"log"

	"taskflow-api/internal/config"
	"taskflow-api/internal/database"
	"taskflow-api/internal/handlers"
	"taskflow-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()
	log.Println("✅ Config loaded")

	// Connect to database
	db := database.Connect(cfg.DatabaseURL)
	defer db.Close()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)
	taskHandler := handlers.NewTaskHandler(db)

	// Setup Gin router
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "TaskFlow API is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/profile", middleware.AuthMiddleware(cfg.JWTSecret, db), authHandler.GetProfile)
		}

		// Task routes (protected)
		tasks := v1.Group("/tasks")
		tasks.Use(middleware.AuthMiddleware(cfg.JWTSecret, db))
		{
			tasks.POST("", taskHandler.CreateTask)
			tasks.GET("", taskHandler.GetTasks)
			tasks.GET("/stats", taskHandler.GetStats)
			tasks.GET("/:id", taskHandler.GetTask)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
			tasks.PATCH("/:id/complete", taskHandler.ToggleComplete)
		}
	}

	// Start server
	log.Println("============================================================")
	log.Printf("🚀 Server starting on port %s", cfg.Port)
	log.Println("📝 API Documentation:")
	log.Printf("   - Health Check: GET http://localhost:%s/health", cfg.Port)
	log.Printf("   - Register:     POST http://localhost:%s/api/v1/auth/register", cfg.Port)
	log.Printf("   - Login:        POST http://localhost:%s/api/v1/auth/login", cfg.Port)
	log.Printf("   - Get Tasks:    GET http://localhost:%s/api/v1/tasks", cfg.Port)
	log.Println("============================================================")

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
