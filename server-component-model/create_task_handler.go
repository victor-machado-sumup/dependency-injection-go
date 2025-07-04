package main

import "fmt"

type CreateTaskHandler struct {
	repository *Repository
}

func NewCreateTaskHandler() (*CreateTaskHandler, error) {
	repository, err := NewRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}
	handler := CreateTaskHandler{
		repository: repository,
	}
	return &handler, nil
}

type CreateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateTaskOutput struct {
	Task Task `json:"task"`
}

func (h *CreateTaskHandler) Handle(input CreateTaskInput) (CreateTaskOutput, error) {

	// Validate required fields
	if input.Title == "" {
		return CreateTaskOutput{}, fmt.Errorf("title is required")
	}

	task := Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      TaskStatusPending,
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
