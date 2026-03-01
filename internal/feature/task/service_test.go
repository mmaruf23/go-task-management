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
		createTaskFunc: func(ctx context.Context, arg db.CreateTaskParams) (db.Task, error) {
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

	task, err := service.CreateTask(ctx, &userID, req)

	assert.NoError(t, err)
	assert.NotEqual(t, 0, task.ID)
	assert.Equal(t, task.Title, "ini sample title")

}

func TestGetTaskByUser_Success(t *testing.T) {
	mockRepo := &mockTaskRepo{
		listTasksByUserFunc: func(ctx context.Context, arg db.ListTaskByUserParams) ([]db.Task, error) {
			var tasks []db.Task
			tasks = append(tasks, db.Task{
				ID:          uuid.New(),
				UserID:      arg.UserID,
				Title:       "sample task",
				Description: nil,
				Status:      "pending",
			})

			return tasks, nil
		},
	}

	service := NewTaskService(mockRepo)

	ctx := context.Background()
	userID := uuid.New()

	tasks, err := service.GetUserTasks(ctx, &userID)

	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, "sample task", tasks[0].Title)

}

// MOCKING

type mockTaskRepo struct {
	createTaskFunc      func(ctx context.Context, arg db.CreateTaskParams) (db.Task, error)
	listTasksByUserFunc func(ctx context.Context, arg db.ListTaskByUserParams) ([]db.Task, error)
}

func (r *mockTaskRepo) CreateTask(ctx context.Context, arg db.CreateTaskParams) (db.Task, error) {
	return r.createTaskFunc(ctx, arg)
}

func (r *mockTaskRepo) ListTaskByUser(ctx context.Context, arg db.ListTaskByUserParams) ([]db.Task, error) {
	return r.listTasksByUserFunc(ctx, arg)
}
