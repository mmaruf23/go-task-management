package task

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/mmaruf23/go-task-management/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTask_Success(t *testing.T) {

	userID := uuid.New()
	description := "ini sample description"
	arg := repo.CreateTaskParams{
		UserID:      userID,
		Title:       "ini sample title",
		Description: &description,
	}

	mockRepo := new(MockTaskRepo)
	mockRepo.On("CreateTask", mock.Anything, arg).Return(repo.Task{
		ID:          userID,
		Title:       arg.Title,
		UserID:      arg.UserID,
		Description: arg.Description,
		Status:      "pending",
	}, nil)

	service := NewTaskService(mockRepo)

	ctx := context.Background()
	req := CreateTaskRequest{
		Title:       arg.Title,
		Description: arg.Description,
	}

	task, err := service.CreateTask(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotEqual(t, 0, task.ID)
	assert.Equal(t, task.Title, arg.Title)
	if task.Description != nil {
		assert.Equal(t, *arg.Description, *task.Description)
	}
}

func TestGetTaskByUser_Success(t *testing.T) {
	mockRepo := new(MockTaskRepo)
	var tasks []repo.Task
	tasks = append(tasks, repo.Task{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		Title:       "sample task",
		Description: nil,
		Status:      "pending",
	})
	mockRepo.On("ListTaskByUser", mock.Anything, mock.Anything).Return(tasks, nil)
	mockRepo.On("CountTaskByUser", mock.Anything, mock.Anything).Return(int64(1), nil)

	service := NewTaskService(mockRepo)

	ctx := context.Background()
	userID := uuid.New()
	pagination := PaginationRequest{}
	pagination.Normalize()

	paginatedtasks, err := service.GetUserTasks(ctx, userID, pagination)

	assert.NoError(t, err)
	assert.Len(t, *paginatedtasks.Data, 1)
	assert.Equal(t, "sample task", (*paginatedtasks.Data)[0].Title)

}

func TestUpdateStatus_Success(t *testing.T) {
	mockRepo := new(MockTaskRepo)
	mockRepo.On("UpdateStatus", mock.Anything, mock.Anything).Return(int64(1), nil)

	service := NewTaskService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()
	taskID := uuid.New()

	err := service.UpdateStatus(ctx, userID, taskID, repo.TaskStatusCompleted)

	assert.NoError(t, err)
}

func TestUpdateStatus_Fail_NotFound(t *testing.T) {
	mockRepo := new(MockTaskRepo)
	mockRepo.On("UpdateStatus", mock.Anything, mock.Anything).Return(int64(0), nil) // err nya nil aja karena asumsi result query mah nggak ada masalah.

	service := NewTaskService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()
	taskID := uuid.New()

	err := service.UpdateStatus(ctx, userID, taskID, repo.TaskStatusCompleted)

	assert.NotNil(t, err)
	assert.Equal(t, "task not found", err.Error())
}

func TestUpdateTask_Success(t *testing.T) {
	mockRepo := new(MockTaskRepo)
	userID := uuid.New()
	taskID := uuid.New()
	newTitle := "updated title"
	newDescp := "updated description"

	task := repo.Task{
		ID:          taskID,
		UserID:      userID,
		Title:       newTitle,
		Description: &newDescp,
		Status:      repo.TaskStatusPending,
		CreatedAt:   pgtype.Timestamp{Time: time.Now()},
		UpdatedAt:   pgtype.Timestamp{Time: time.Now()},
	}
	mockRepo.On("UpdateTask", mock.Anything, mock.Anything).Return(task, nil)
	mockRepo.On("GetTaskByID", mock.Anything, mock.Anything).Return(repo.Task{}, nil)

	service := NewTaskService(mockRepo)
	ctx := context.Background()

	req := UpdateTaskRequest{
		Title:       &newTitle,
		Description: &newDescp,
	}

	taskResponse, err := service.Update(ctx, userID, taskID, req)

	assert.NoError(t, err)
	assert.Equal(t, taskResponse.Title, newTitle)
	assert.Equal(t, *taskResponse.Description, newDescp)

}

// MOCKING

type MockTaskRepo struct {
	mock.Mock
}

func (m *MockTaskRepo) CreateTask(ctx context.Context, arg repo.CreateTaskParams) (repo.Task, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repo.Task), args.Error(1)
}

func (m *MockTaskRepo) CountTaskByUser(ctx context.Context, arg uuid.UUID) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTaskRepo) GetTaskByID(ctx context.Context, arg repo.GetTaskByIDParams) (repo.Task, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repo.Task), args.Error(1)
}

func (m *MockTaskRepo) ListTaskByUser(ctx context.Context, arg repo.ListTaskByUserParams) ([]repo.Task, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]repo.Task), args.Error(1)
}

func (m *MockTaskRepo) UpdateStatus(ctx context.Context, arg repo.UpdateStatusParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTaskRepo) UpdateTask(ctx context.Context, arg repo.UpdateTaskParams) (repo.Task, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repo.Task), args.Error(1)
}
