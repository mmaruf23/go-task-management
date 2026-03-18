package task

import (
	"errors"

	"github.com/mmaruf23/go-task-management/internal/repository"
)

type CreateTaskRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description" binding:"required"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type TaskStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (r *TaskStatusRequest) Parse() (repository.TaskStatus, error) {
	status := repository.TaskStatus(r.Status)
	switch status {
	case repository.TaskStatusPending,
		repository.TaskStatusCompleted:
	default:
		return "", errors.New("invalid status")
	}
	return status, nil
}

type PaginationRequest struct {
	Page  int32 `form:"page"`
	Limit int32 `form:"limit"`
}

func (r *PaginationRequest) Normalize() {
	if r.Page <= 0 {
		r.Page = 1
	}

	if r.Limit <= 0 {
		r.Limit = 10
	}

	if r.Limit > 100 {
		r.Limit = 100
	}
}

func (r *PaginationRequest) Offset() int32 {
	return (r.Limit) * (r.Page - 1)
}
