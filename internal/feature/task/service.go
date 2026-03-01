package task

import (
	"context"
	"errors"

	"github.com/google/uuid"
	repo "github.com/mmaruf23/go-task-management/internal/repository"
	"github.com/mmaruf23/go-task-management/pkg/util"
)

type TaskRepository interface {
	CountTaskByUser(ctx context.Context, userID uuid.UUID) (int64, error)
	CreateTask(ctx context.Context, arg repo.CreateTaskParams) (repo.Task, error)
	// GetTaskByID(ctx context.Context, arg GetTaskByIDParams) (Task, error)
	ListTaskByUser(ctx context.Context, arg repo.ListTaskByUserParams) ([]repo.Task, error)
	// UpdateStatus(ctx context.Context, arg UpdateStatusParams) (int64, error)
	// UpdateTask(ctx context.Context, arg UpdateTaskParams) (Task, error)
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, userID uuid.UUID, req CreateTaskRequest) (*repo.Task, error) {
	params := repo.CreateTaskParams{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
	}

	task, err := s.repo.CreateTask(ctx, params)
	if err != nil {
		return nil, errors.New("failed create new task")
	}

	return &task, nil

}

func (s *TaskService) GetUserTasks(ctx context.Context, userID uuid.UUID, pagination PaginationRequest) (*util.Paginated[*[]repo.Task], error) {

	params := repo.ListTaskByUserParams{
		UserID: userID,
		Limit:  pagination.Limit,
		Offset: pagination.Offset(),
	}

	tasks, err := s.repo.ListTaskByUser(ctx, params)
	if err != nil {
		return nil, err
	}

	count, err := s.repo.CountTaskByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &util.Paginated[*[]repo.Task]{
		Data: &tasks,
		Meta: *util.BuildPaginationMeta(pagination.Page, pagination.Limit, count),
	}, nil
}
