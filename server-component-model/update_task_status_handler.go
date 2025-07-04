package main

import "fmt"

type UpdateTaskStatusHandler struct {
	repository *Repository
}

func NewUpdateTaskStatusHandler() (*UpdateTaskStatusHandler, error) {
	repository, err := NewRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}
	handler := UpdateTaskStatusHandler{
		repository: repository,
	}
	return &handler, nil
}

type UpdateTaskStatusInput struct {
	TaskID int    `json:"taskId"`
	Status string `json:"status"`
}

type UpdateTaskStatusOutput struct {
	Task Task `json:"task"`
}

func (h *UpdateTaskStatusHandler) Handle(input UpdateTaskStatusInput) (UpdateTaskStatusOutput, error) {

	taskStatus := TaskStatus(input.Status)

	// Validate status value
	if taskStatus != TaskStatusPending && taskStatus != TaskStatusCompleted {
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
