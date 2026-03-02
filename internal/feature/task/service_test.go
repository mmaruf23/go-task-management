package task

import (
	"context"
	"testing"

	"github.com/google/uuid"
	repo "github.com/mmaruf23/go-task-management/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mockTaskRepo{
		createTaskFunc: func(ctx context.Context, arg repo.CreateTaskParams) (repo.Task, error) {
			return repo.Task{
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
	req := CreateTaskRequest{
		Title:       "ini sample title",
		Description: nil,
	}

	task, err := service.CreateTask(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotEqual(t, 0, task.ID)
	assert.Equal(t, task.Title, "ini sample title")

}

func TestGetTaskByUser_Success(t *testing.T) {
	mockRepo := &mockTaskRepo{
		listTasksByUserFunc: func(ctx context.Context, arg repo.ListTaskByUserParams) ([]repo.Task, error) {
			var tasks []repo.Task
			tasks = append(tasks, repo.Task{
				ID:          uuid.New(),
				UserID:      arg.UserID,
				Title:       "sample task",
				Description: nil,
				Status:      "pending",
			})

			return tasks, nil
		},
		countTaskFunc: func(ctx context.Context, userID uuid.UUID) (int64, error) {
			return 1, nil
		},
	}

	service := NewTaskService(mockRepo)

	ctx := context.Background()
	userID := uuid.New()
	pagination := PaginationRequest{}
	pagination.Normalize()

	tasks, err := service.GetUserTasks(ctx, userID, pagination)

	assert.NoError(t, err)
	assert.Len(t, *tasks.Data, 1)
	assert.Equal(t, "sample task", (*tasks.Data)[0].Title)

}

func TestUpdateStatus_Success(t *testing.T) {
	mockRepo := &mockTaskRepo{
		updateStatusFun: func(ctx context.Context, arg repo.UpdateStatusParams) (int64, error) {
			return 1, nil
		},
	}

	service := NewTaskService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()
	taskID := uuid.New()

	err := service.UpdateStatus(ctx, userID, taskID, repo.TaskStatusCompleted)

	assert.NoError(t, err)
}
func TestUpdateStatus_Fail_NotFound(t *testing.T) {
	mockRepo := &mockTaskRepo{
		updateStatusFun: func(ctx context.Context, arg repo.UpdateStatusParams) (int64, error) {
			return 0, nil
		},
	}

	service := NewTaskService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()
	taskID := uuid.New()

	err := service.UpdateStatus(ctx, userID, taskID, repo.TaskStatusCompleted)

	assert.NotNil(t, err)
	assert.Equal(t, "task not found", err.Error())
}

// MOCKING

type mockTaskRepo struct {
	createTaskFunc      func(ctx context.Context, arg repo.CreateTaskParams) (repo.Task, error)
	listTasksByUserFunc func(ctx context.Context, arg repo.ListTaskByUserParams) ([]repo.Task, error)
	countTaskFunc       func(ctx context.Context, userID uuid.UUID) (int64, error)
	updateStatusFun     func(ctx context.Context, arg repo.UpdateStatusParams) (int64, error)
}

func (r *mockTaskRepo) CreateTask(ctx context.Context, arg repo.CreateTaskParams) (repo.Task, error) {
	return r.createTaskFunc(ctx, arg)
}

func (r *mockTaskRepo) ListTaskByUser(ctx context.Context, arg repo.ListTaskByUserParams) ([]repo.Task, error) {
	return r.listTasksByUserFunc(ctx, arg)
}

func (r *mockTaskRepo) CountTaskByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.countTaskFunc(ctx, userID)
}

func (r *mockTaskRepo) UpdateStatus(ctx context.Context, arg repo.UpdateStatusParams) (int64, error) {
	return r.updateStatusFun(ctx, arg)
}
