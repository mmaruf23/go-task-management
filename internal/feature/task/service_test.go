package task

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/mmaruf23/go-task-management/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mockTaskRepo{
		createTaskFunc: func(ctx context.Context, arg *db.CreateTaskParams) (db.Task, error) {
			return db.Task{
				ID:          uuid.New(),
				Title:       arg.Title,
				UserID:      arg.UserID,
				Description: arg.Description,
				Status:      "pending",
			}, nil
		},
	}
	service := NewTaskService(mockRepo)

	userID := uuid.New()
	req := &CreateTaskRequest{
		Title:       "ini sample title",
		Description: nil,
	}

	task, err := service.CreateTask(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotEqual(t, 0, task.ID)
	assert.Equal(t, task.Title, "ini sample title")

}

// MOCKING

type mockTaskRepo struct {
	createTaskFunc func(ctx context.Context, arg *db.CreateTaskParams) (db.Task, error)
}

func (r *mockTaskRepo) CreateTask(ctx context.Context, arg *db.CreateTaskParams) (db.Task, error) {
	return r.createTaskFunc(ctx, arg)
}
