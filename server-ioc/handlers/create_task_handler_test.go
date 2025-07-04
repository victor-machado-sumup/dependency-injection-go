package handlers_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/sumup/dependency-injection-go/server-ioc/handlers"
	"github.com/sumup/dependency-injection-go/server-ioc/repository"
)

func TestCreateTaskHandler_Handle(t *testing.T) {
	// Set up database connection
	pool, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable")
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM tasks")
		pool.Close()
	})
	// Create dependencies
	repo := repository.NewRepository(pool)
	handler := handlers.NewCreateTaskHandler(repo)

	// Create test input
	input := handlers.CreateTaskInput{
		Title:       "Test Task",
		Description: "This is a test task",
	}

	// Execute the handler
	output, err := handler.Handle(input)

	// Assert no error occurred
	require.NoError(t, err)
	require.NotEmpty(t, output.Task.ID)
	require.Equal(t, input.Title, output.Task.Title)
	require.Equal(t, input.Description, output.Task.Description)
	require.Equal(t, repository.TaskStatusPending, output.Task.Status)

	// Verify the task was actually created in the database
	var task repository.Task
	err = pool.QueryRow(context.Background(),
		"SELECT id, title, description, status FROM tasks WHERE id = $1",
		output.Task.ID,
	).Scan(&task.ID, &task.Title, &task.Description, &task.Status)

	require.NoError(t, err)
	require.Equal(t, input.Title, task.Title)
	require.Equal(t, input.Description, task.Description)
	require.Equal(t, repository.TaskStatusPending, task.Status)
}
