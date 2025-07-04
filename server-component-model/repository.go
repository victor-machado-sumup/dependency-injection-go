package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	connString string
}

func NewRepository() (*Repository, error) {
	repository := Repository{
		connString: "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable",
	}

	return &repository, nil
}

func (r *Repository) getConnection() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), r.connString)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (r *Repository) GetTaskById(id int) (Task, error) {
	conn, err := r.getConnection()
	if err != nil {
		return Task{}, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	query := `SELECT id, title, description, status FROM tasks WHERE id = $1`
	var task Task
	err = conn.QueryRow(context.Background(), query, id).Scan(&task.ID, &task.Title, &task.Description, &task.Status)
	if err != nil {
		return Task{}, fmt.Errorf("failed to get task: %w", err)
	}
	return task, nil
}

func (r *Repository) CreateTask(task Task) (Task, error) {
	conn, err := r.getConnection()
	if err != nil {
		return Task{}, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Insert the task into the database
	query := `INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err = conn.QueryRow(context.Background(), query, task.Title, task.Description, task.Status).Scan(&id)
	if err != nil {
		return Task{}, fmt.Errorf("failed to create task: %w", err)
	}

	// Get the created task
	return r.GetTaskById(id)
}

func (r *Repository) UpdateTaskStatus(id int, status TaskStatus) (Task, error) {
	conn, err := r.getConnection()
	if err != nil {
		return Task{}, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	query := `UPDATE tasks SET status = $1 WHERE id = $2 RETURNING id`
	var taskId int
	err = conn.QueryRow(context.Background(), query, status, id).Scan(&taskId)
	if err != nil {
		return Task{}, fmt.Errorf("failed to update task status: %w", err)
	}

	// Get the updated task
	return r.GetTaskById(taskId)
}

func (r *Repository) GetAllTasks() ([]Task, error) {
	conn, err := r.getConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	query := `SELECT id, title, description, status FROM tasks`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status)
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
