package handlers

import (
	"fmt"

	"github.com/sumup/dependency-injection-go/server-ioc/repository"
)

type CreateTaskHandler struct {
	repository repository.IRepository
}

func NewCreateTaskHandler(repository repository.IRepository) *CreateTaskHandler {
	handler := CreateTaskHandler{
		repository: repository,
	}
	return &handler
}

type CreateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateTaskOutput struct {
	Task repository.Task `json:"task"`
}

func (h *CreateTaskHandler) Handle(input CreateTaskInput) (CreateTaskOutput, error) {

	// Validate required fields
	if input.Title == "" {
		return CreateTaskOutput{}, fmt.Errorf("title is required")
	}

	task := repository.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      repository.TaskStatusPending,
	}

	// Create the task using repository
	createdTask, err := h.repository.CreateTask(task)
	if err != nil {
		return CreateTaskOutput{}, fmt.Errorf("failed to create task: %w", err)
	}

	return CreateTaskOutput{
		Task: createdTask,
	}, nil
}
