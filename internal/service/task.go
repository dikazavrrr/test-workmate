package service

import (
	"context"
	"test-workmate/internal/domain"
)

type TaskRepository interface {
	SaveTask(ctx context.Context, task *domain.Task) error
	UpdateTask(ctx context.Context, task *domain.Task) error
	GetTask(ctx context.Context, id string) (*domain.Task, error)
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(r TaskRepository) *TaskService {
	return &TaskService{repo: r}
}
