package handlers_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/sumup/dependency-injection-go/server-ioc/handlers"
	"github.com/sumup/dependency-injection-go/server-ioc/repository"
)

func TestUpdateTaskStatusHandler_Handle(t *testing.T) {
	// Set up database connection
	pool, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable")
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM tasks")
		pool.Close()
	})

	// Create dependencies
	repo := repository.NewRepository(pool)
	handler := handlers.NewUpdateTaskStatusHandler(repo)

	// First, create a task to update
	newTask := repository.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      repository.TaskStatusPending,
	}
	task, err := repo.CreateTask(newTask)
	require.NoError(t, err)
	require.Equal(t, repository.TaskStatusPending, task.Status)

	// Create test input for status update
	input := handlers.UpdateTaskStatusInput{
		TaskID: task.ID,
		Status: string(repository.TaskStatusCompleted),
	}

	// Execute the handler
	output, err := handler.Handle(input)

	// Assert no error occurred and status was updated
	require.NoError(t, err)
	require.Equal(t, task.ID, output.Task.ID)
	require.Equal(t, task.Title, output.Task.Title)
	require.Equal(t, task.Description, output.Task.Description)
	require.Equal(t, repository.TaskStatusCompleted, output.Task.Status)

	// Verify the task was actually updated in the database
	var updatedTask repository.Task
	err = pool.QueryRow(context.Background(),
		"SELECT id, title, description, status FROM tasks WHERE id = $1",
		task.ID,
	).Scan(&updatedTask.ID, &updatedTask.Title, &updatedTask.Description, &updatedTask.Status)

	require.NoError(t, err)
	require.Equal(t, task.ID, updatedTask.ID)
	require.Equal(t, task.Title, updatedTask.Title)
	require.Equal(t, task.Description, updatedTask.Description)
	require.Equal(t, repository.TaskStatusCompleted, updatedTask.Status)

	// Test invalid status
	invalidInput := handlers.UpdateTaskStatusInput{
		TaskID: task.ID,
		Status: "invalid_status",
	}

	// Execute the handler with invalid status
	_, err = handler.Handle(invalidInput)

	// Assert error occurred for invalid status
	require.Error(t, err)
	require.Contains(t, err.Error(), "Invalid status value")
}
