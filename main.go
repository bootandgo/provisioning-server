package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"provisioning-server/handlers"
	"provisioning-server/middleware"
	"provisioning-server/store"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to Postgres
	connStr := os.Getenv("DB_CONNECTION_STRING")
	if connStr == "" {
		panic("DB_CONNECTION_STRING env var is required")
	}

	pgStore, err := store.NewPostgresStore(connStr)
	if err != nil {
		panic("failed to connect to postgres: " + err.Error())
	}

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pgStore.Ping(ctx); err != nil {
		panic("postgres ping failed: " + err.Error())
	}

	r := gin.Default()

	// Create root user
	username := os.Getenv("ROOT_USERNAME")
	if username == "" {
		panic("ROOT_USERNAME env var is required")
	}

	password := os.Getenv("ROOT_PASSWORD")
	if password == "" {
		panic("ROOT_PASSWORD env var is required")
	}

	// Create root admin if none exists
	handlers.CreateRootAdminIfNeeded(pgStore, username, password)

	// JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET env var is required")
	}

	// Middleware
	adminAuth := middleware.AdminAuthMiddleware(pgStore, jwtSecret)

	// Admin routes
	adminRoutes := r.Group("/admin")
	{
		// Public admin endpoints
		adminRoutes.POST("/register", handlers.AdminRegisterHandler(pgStore))
		adminRoutes.POST("/login", handlers.AdminLoginHandler(pgStore, jwtSecret))

		// Authenticated admin endpoints
		authAdmin := adminRoutes.Group("").Use(adminAuth)
		{
			authAdmin.POST("/invite", handlers.CreateInviteHandler(pgStore))
		}
	}

	// Server routes
	serverRoutes := r.Group("/servers")
	{
		// Public server endpoints
		serverRoutes.POST("/register", handlers.RegisterServerHandler(pgStore))
		serverRoutes.GET("/:id/status", handlers.ServerStatusHandler(pgStore))

		// Authenticated server endpoints
		authServers := serverRoutes.Group("").Use(adminAuth)
		{
			authServers.GET("", handlers.ListServersHandler(pgStore))
			authServers.POST("/approve", handlers.ApproveServerHandler(pgStore))
		}
	}

	err = r.Run(":8080")
	if err != nil {
		panic("failed to start server: " + err.Error())
	}
}
