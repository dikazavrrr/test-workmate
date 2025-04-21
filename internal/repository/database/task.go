package database

import (
	"context"
	"database/sql"
	"fmt"
	"test-workmate/internal/domain"
)

type TaskRepo struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) SaveTask(ctx context.Context, task *domain.Task) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO tasks (id, status, result, error, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		task.ID, task.Status, task.Result, task.Error, task.CreatedAt, task.UpdatedAt,
	)
	return err
}

func (r *TaskRepo) UpdateTask(ctx context.Context, task *domain.Task) error {
	if ctx.Err() != nil {
		fmt.Println("Context error before query:", ctx.Err())
	}

	query := `
        UPDATE tasks 
        SET status = $1, result = $2, updated_at = $3 
        WHERE id = $4
    `
	_, err := r.db.ExecContext(ctx, query, task.Status, task.Result, task.UpdatedAt, task.ID)
	if err != nil {
		fmt.Printf("UpdateTask error: %v\n", err)
	}
	return err
}

func (r *TaskRepo) GetTask(ctx context.Context, id string) (*domain.Task, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, status, result, error, created_at, updated_at FROM tasks WHERE id=$1`, id)

	var task domain.Task
	err := row.Scan(&task.ID, &task.Status, &task.Result, &task.Error, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &task, nil
}
