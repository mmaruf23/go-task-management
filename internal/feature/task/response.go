package task

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mmaruf23/go-task-management/internal/repository"
)

type TaskResponse struct {
	ID          uuid.UUID             `json:"id"`
	Title       string                `json:"title"`
	Description *string               `json:"description"`
	Status      repository.TaskStatus `json:"status"`
	CreatedAt   pgtype.Timestamp      `json:"createdAt"`
	UpdatedAt   pgtype.Timestamp      `json:"updatedAt"`
}

func ToTaskResponse(task *repository.Task) *TaskResponse {
	return &TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func ToTaskResponses(tasks *[]repository.Task) *[]TaskResponse {
	result := make([]TaskResponse, 0, len(*tasks))

	for _, t := range *tasks {
		result = append(result, *ToTaskResponse(&t))
	}
	return &result
}
