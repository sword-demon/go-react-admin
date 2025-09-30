package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sword-demon/go-react-admin/internal/admin/biz"
	"github.com/sword-demon/go-react-admin/internal/admin/store"
	"github.com/sword-demon/go-react-admin/internal/pkg/cache"
	"github.com/sword-demon/go-react-admin/internal/pkg/config"
	"github.com/sword-demon/go-react-admin/internal/pkg/db"
)

func main() {
	// 1. Load configuration
	configPath := "configs/config.yml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("⚠️  Config file not found, using default configuration")
		configPath = "" // Will use default config
	}

	var cfg *config.Config
	var err error
	if configPath != "" {
		cfg, err = config.Load(configPath)
		if err != nil {
			log.Fatalf("❌ Failed to load config: %v", err)
		}
		log.Println("✅ Configuration loaded from", configPath)
	} else {
		cfg = config.Default()
		log.Println("✅ Using default configuration")
	}

	// 2. Initialize database (GORM)
	database, err := db.InitDB(cfg.Database.ToDBConfig())
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer func() {
		if err := db.Close(database); err != nil {
			log.Printf("⚠️  Failed to close database: %v", err)
		}
	}()

	// 3. Initialize Redis (optional, warn if fails)
	redisClient, err := cache.InitRedis(cfg.Redis.ToRedisConfig())
	if err != nil {
		log.Printf("⚠️  Failed to initialize Redis (will continue without cache): %v", err)
		redisClient = nil
	} else {
		defer func() {
			if err := redisClient.Close(); err != nil {
				log.Printf("⚠️  Failed to close Redis: %v", err)
			}
		}()
	}

	// 4. Run database migrations
	if err := db.AutoMigrate(database); err != nil {
		log.Fatalf("❌ Failed to run database migrations: %v", err)
	}

	// 5. Initialize store layer
	dataStore := store.NewStore(database)
	log.Println("✅ Store layer initialized")

	// 6. Initialize biz layer
	bizLayer := biz.NewBiz(dataStore, redisClient)
	log.Println("✅ Biz layer initialized")

	// Prevent unused variable error (will be used in controllers later)
	_ = bizLayer

	// TODO: Initialize controllers (pass bizLayer)
	// TODO: Setup middleware (JWT, Permission, CORS)
	// TODO: Register routes (see internal/admin/router.go)

	// Set Gin mode based on config
	gin.SetMode(cfg.Server.Mode)
	r := gin.Default()

	// Health check endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"version": "v1.0.0-alpha",
			"status":  "healthy",
		})
	})

	// API v1 routes (will be moved to internal/admin/router.go)
	v1 := r.Group("/api/v1")
	{
		// Auth routes
		v1.POST("/auth/login", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "login endpoint"})
		})

		// User routes (protected by JWT + Permission middleware)
		users := v1.Group("/users")
		{
			users.GET("", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "list users"})
			})
			users.POST("", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "create user"})
			})
		}
	}

	port := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("\n🚀 Server starting on http://localhost:%d\n", cfg.Server.Port)
	fmt.Printf("📊 Mode: %s\n", cfg.Server.Mode)
	fmt.Println("📚 API Documentation: http://localhost:8080/swagger/index.html (coming soon)")
	fmt.Printf("💚 Health Check: http://localhost:%d/ping\n", cfg.Server.Port)

	if err := r.Run(port); err != nil {
		log.Fatal("❌ Server startup failed:", err)
	}
}
