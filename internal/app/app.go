package app

import (
	"log"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mmaruf23/go-task-management/internal/config"
	"github.com/mmaruf23/go-task-management/internal/feature/auth"
	"github.com/mmaruf23/go-task-management/internal/feature/task"
	"github.com/mmaruf23/go-task-management/internal/repository"
)

type App struct {
	router *gin.Engine
	pool   *pgxpool.Pool
}

func New() *App {

	cfg := config.Load()
	pool := InitDB(cfg)

	repo := repository.New(pool)

	jwtService := auth.NewJWTService(cfg.JWTSecret)
	authService := auth.NewAuthService(repo, jwtService)
	taskService := task.NewTaskService(repo)

	authHandler := auth.NewAuthHandler(authService)
	taskHandler := task.NewTaskHandler(taskService)

	authMiddleware := auth.AuthMiddleware(jwtService)

	r := InitRouter()
	api := r.Group("/")

	authHandler.Routes(api, authMiddleware)
	taskHandler.Routes(api, authMiddleware)

	return &App{
		router: r,
		pool:   pool,
	}
}

func (a *App) Run() {
	go func() {
		if err := a.router.Run(":8080"); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("shutting down server ... and closing pool?")
	a.pool.Close()
}
