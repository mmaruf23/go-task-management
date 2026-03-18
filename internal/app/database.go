package app

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mmaruf23/go-task-management/internal/config"
	"github.com/mmaruf23/go-task-management/internal/db"
)

func InitDB(cfg *config.Config) *pgxpool.Pool {
	pool, err := db.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	return pool
}
