package service

import (
	"context"
	"test-workmate/internal/domain"
)

func (s *TaskService) GetTask(ctx context.Context, id string) (*domain.Task, error) {
	return s.repo.GetTask(ctx, id)
}
