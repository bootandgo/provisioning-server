package main

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"provisioning-server/handlers"
	"provisioning-server/middleware"
	"provisioning-server/store"
)

func main() {
	// Connect to Postgres
	connStr := os.Getenv("DB_CONNECTION_STRING")
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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET env var is required")
	}

	// Admin routes
	r.POST("/admin/register", handlers.AdminRegisterHandler(pgStore))
	r.POST("/admin/login", handlers.AdminLoginHandler(pgStore, jwtSecret))

	// Server routes
	serverRoutes := r.Group("/servers")
	{
		serverRoutes.POST("/register", handlers.RegisterServerHandler(pgStore))
	}

	// Admin-protected routes
	adminAuth := middleware.AdminAuthMiddleware(pgStore, jwtSecret)
	adminRoutes := r.Group("/servers")
	adminRoutes.Use(adminAuth)
	{
		adminRoutes.GET("", handlers.ListServersHandler(pgStore))
		adminRoutes.POST("/approve", handlers.ApproveServerHandler(pgStore))
	}

	err = r.Run(":8080")
	if err != nil {
		panic("failed to start server: " + err.Error())
	}
}
