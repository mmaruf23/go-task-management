package task

import (
	"context"
	"errors"

	"github.com/google/uuid"
	repo "github.com/mmaruf23/go-task-management/internal/repository"
	"github.com/mmaruf23/go-task-management/pkg/util"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, arg repo.CreateTaskParams) (repo.Task, error)
	CountTaskByUser(ctx context.Context, userID uuid.UUID) (int64, error)
	GetTaskByID(ctx context.Context, arg repo.GetTaskByIDParams) (repo.Task, error)
	ListTaskByUser(ctx context.Context, arg repo.ListTaskByUserParams) ([]repo.Task, error)
	UpdateStatus(ctx context.Context, arg repo.UpdateStatusParams) (int64, error)
	UpdateTask(ctx context.Context, arg repo.UpdateTaskParams) (repo.Task, error)
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, userID uuid.UUID, req CreateTaskRequest) (*TaskResponse, error) {
	params := repo.CreateTaskParams{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
	}

	task, err := s.repo.CreateTask(ctx, params)
	if err != nil {
		return nil, errors.New("failed create new task")
	}

	return ToTaskResponse(&task), nil

}

func (s *TaskService) GetUserTasks(ctx context.Context, userID uuid.UUID, pagination PaginationRequest) (*util.Paginated[*[]TaskResponse], error) {

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

	return &util.Paginated[*[]TaskResponse]{
		Data: ToTaskResponses(&tasks),
		Meta: *util.BuildPaginationMeta(pagination.Page, pagination.Limit, count),
	}, nil
}

func (s *TaskService) UpdateStatus(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, status repo.TaskStatus) error {
	params := repo.UpdateStatusParams{
		UserID: userID,
		ID:     taskID,
		Status: status,
	}

	rowsAffected, err := s.repo.UpdateStatus(ctx, params)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

func (s *TaskService) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req UpdateTaskRequest) (*TaskResponse, error) {
	task, err := s.repo.GetTaskByID(ctx, repo.GetTaskByIDParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return nil, errors.New("no task found")
	}

	params := repo.UpdateTaskParams{
		UserID:      userID,
		ID:          id,
		Title:       *req.Title,
		Description: req.Description,
	}

	if params.Title == "" {
		params.Title = task.Title
	}

	if params.Description == nil {
		params.Description = task.Description
	}

	updatedTask, err := s.repo.UpdateTask(ctx, params)
	if err != nil {
		return nil, err
	}
	return ToTaskResponse(&updatedTask), nil
}
