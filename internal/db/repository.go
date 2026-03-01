package db

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
}

type TaskRepository interface {
	// CountTaskByUser(ctx context.Context, userID uuid.UUID) (int64, error)
	CreateTask(ctx context.Context, arg *CreateTaskParams) (Task, error)
	// GetTaskByID(ctx context.Context, arg GetTaskByIDParams) (Task, error)
	// ListTaskByUser(ctx context.Context, arg ListTaskByUserParams) ([]Task, error)
	// UpdateStatus(ctx context.Context, arg UpdateStatusParams) (int64, error)
	// UpdateTask(ctx context.Context, arg UpdateTaskParams) (Task, error)
}
