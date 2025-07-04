package handlers

import (
	"fmt"

	"github.com/sumup/dependency-injection-go/server-dependency-injection/repository"
)

type UpdateTaskStatusHandler struct {
	repository repository.IRepository
}

func NewUpdateTaskStatusHandler(repository repository.IRepository) *UpdateTaskStatusHandler {
	handler := UpdateTaskStatusHandler{
		repository: repository,
	}
	return &handler
}

type UpdateTaskStatusInput struct {
	TaskID int    `uri:"id" json:"taskId"`
	Status string `json:"status"`
}

type UpdateTaskStatusOutput struct {
	Task repository.Task `json:"task"`
}

func (h *UpdateTaskStatusHandler) Handle(input UpdateTaskStatusInput) (UpdateTaskStatusOutput, error) {

	taskStatus := repository.TaskStatus(input.Status)

	// Validate status value
	if taskStatus != repository.TaskStatusPending && taskStatus != repository.TaskStatusCompleted {
		return UpdateTaskStatusOutput{}, fmt.Errorf("Invalid status value. Must be either 'pending' or 'completed'")
	}

	// Update the task status using repository
	updatedTask, err := h.repository.UpdateTaskStatus(input.TaskID, taskStatus)
	if err != nil {
		return UpdateTaskStatusOutput{}, fmt.Errorf("failed to update task: %w", err)
	}

	return UpdateTaskStatusOutput{
		Task: updatedTask,
	}, nil
}
