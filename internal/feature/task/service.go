package task

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mmaruf23/go-task-management/internal/db"
)

type TaskService struct {
	repo db.TaskRepository
}

func NewTaskService(repo db.TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, userID *uuid.UUID, req *CreateTaskRequest) (db.Task, error) {
	params := db.CreateTaskParams{
		UserID:      *userID,
		Title:       req.Title,
		Description: req.Description,
	}

	task, err := s.repo.CreateTask(ctx, params)
	if err != nil {
		return task, errors.New("failed create new task")
	}

	return task, nil

}

func (s *TaskService) GetUserTasks(ctx context.Context, userID *uuid.UUID) ([]db.Task, error) {
	params := db.ListTaskByUserParams{
		UserID: *userID,
		Limit:  10,
		Offset: 1,
	}

	tasks, err := s.repo.ListTaskByUser(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, errors.New("cannot found any tasks")
	}

	return tasks, nil
}
