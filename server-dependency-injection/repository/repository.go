package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type IRepository interface {
	GetTaskById(id int) (Task, error)
	CreateTask(task Task) (Task, error)
	UpdateTaskStatus(id int, status TaskStatus) (Task, error)
	GetAllTasks() ([]Task, error)
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) IRepository {
	repository := Repository{
		pool: pool,
	}

	return &repository
}

func (r *Repository) GetTaskById(id int) (Task, error) {

	query := `SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE id = $1`
	var task Task
	err := r.pool.QueryRow(context.Background(), query, id).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return Task{}, fmt.Errorf("failed to get task: %w", err)
	}
	return task, nil
}

func (r *Repository) CreateTask(task Task) (Task, error) {

	// Insert the task into the database
	query := `INSERT INTO tasks (title, description, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int
	err := r.pool.QueryRow(context.Background(), query, task.Title, task.Description, task.Status, time.Now().UTC(), time.Now().UTC()).Scan(&id)
	if err != nil {
		return Task{}, fmt.Errorf("failed to create task: %w", err)
	}

	// Get the created task
	return r.GetTaskById(id)
}

func (r *Repository) UpdateTaskStatus(id int, status TaskStatus) (Task, error) {

	query := `UPDATE tasks SET status = $1, updated_at = $2 WHERE id = $3 RETURNING id`
	var taskId int
	err := r.pool.QueryRow(context.Background(), query, status, time.Now().UTC(), id).Scan(&taskId)
	if err != nil {
		return Task{}, fmt.Errorf("failed to update task status: %w", err)
	}

	// Get the updated task
	return r.GetTaskById(taskId)
}

func (r *Repository) GetAllTasks() ([]Task, error) {
	query := `SELECT id, title, description, status, created_at, updated_at FROM tasks ORDER BY id`
	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over tasks: %w", err)
	}

	return tasks, nil
}
