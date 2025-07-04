package main

import "fmt"

type GetTasksHandler struct {
	repository *Repository
}

type GetTasksOutput struct {
	Tasks []Task `json:"tasks"`
}

func NewGetTasksHandler() (*GetTasksHandler, error) {
	repository, err := NewRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}
	handler := GetTasksHandler{
		repository: repository,
	}
	return &handler, nil
}

func (h *GetTasksHandler) Handle() (GetTasksOutput, error) {
	tasks, err := h.repository.GetAllTasks()
	if err != nil {
		return GetTasksOutput{}, fmt.Errorf("failed to retrieve tasks: %w", err)
	}

	return GetTasksOutput{
		Tasks: tasks,
	}, nil
}
