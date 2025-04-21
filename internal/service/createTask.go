package service

import (
	"context"
	"fmt"
	"test-workmate/internal/domain"
	"time"

	"github.com/google/uuid"
)

func (s *TaskService) CreateTask(ctx context.Context) (string, error) {
	task := &domain.Task{
		ID:        uuid.NewString(),
		Status:    domain.StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.repo.SaveTask(ctx, task)
	if err != nil {
		return "", err
	}

	bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	go func() {
		defer cancel()
		s.runTask(bgCtx, task)
	}()

	return task.ID, nil
}

func (s *TaskService) runTask(ctx context.Context, task *domain.Task) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Minute)
	defer cancel()

	task.Status = domain.StatusRunning
	task.UpdatedAt = time.Now()
	fmt.Println("runTask status changed to running")

	err := s.repo.UpdateTask(ctx, task)
	if err != nil {
		fmt.Printf("UpdateTask error: %v\n", err)
		return
	}

	time.Sleep(3 * time.Minute)

	task.Status = domain.StatusDone
	task.Result = "completed successfully"
	task.UpdatedAt = time.Now()
	fmt.Println("runTask status changed to done")

	err = s.repo.UpdateTask(ctx, task)
	if err != nil {
		fmt.Printf("UpdateTask error: %v\n", err)
	}
}
