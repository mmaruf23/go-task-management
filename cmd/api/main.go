package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mmaruf23/go-task-management/internal/config"
	"github.com/mmaruf23/go-task-management/internal/db"
	"github.com/mmaruf23/go-task-management/internal/feature/auth"
)

func main() {
	r := gin.Default()

	cfg := config.Load()
	pool, err := config.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer pool.Close()

	queries := db.New(pool)

	jwtService := auth.NewJWTService(cfg.JWTSecret)
	authService := auth.NewAuthService(queries, jwtService)
	authHandler := auth.NewAuthHandler(authService)

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
