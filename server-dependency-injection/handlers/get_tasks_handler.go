package handlers

import (
	"fmt"

	"github.com/sumup/dependency-injection-go/server-dependency-injection/repository"
)

type GetTasksHandler struct {
	repository repository.IRepository
}

type GetTasksOutput struct {
	Tasks []repository.Task `json:"tasks"`
}

func NewGetTasksHandler(repository repository.IRepository) *GetTasksHandler {
	handler := GetTasksHandler{
		repository: repository,
	}
	return &handler
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
