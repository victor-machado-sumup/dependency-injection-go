package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	port := ":8080"

	// Register single generic handler for all routes
	server := &http.Server{
		Addr: port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Log incoming request
			log.Printf("Received %s request to %s", r.Method, r.URL.Path)

			// Handle different paths and methods programmatically
			switch {
			case r.URL.Path == "/health" && r.Method == "GET":
				handleHealth(w)
			case r.URL.Path == "/tasks" && r.Method == "GET":
				handleGetTasks(w)
			case r.URL.Path == "/tasks" && r.Method == "POST":
				handleCreateTask(w, r)
			case r.Method == "POST" && len(r.URL.Path) > 7 && r.URL.Path[:7] == "/tasks/":
				handleUpdateTaskStatus(w, r)
			default:
				// Handle 404 Not Found
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "404 - Path %s not found\n", r.URL.Path)
			}
		}),
	}

	// Start the server
	fmt.Printf("Server starting on port %s\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// handleHealth processes the health check request
func handleHealth(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// handleCreateTask handles POST requests to create new tasks
func handleCreateTask(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var input CreateTaskInput
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing request body: %v", err)
		return
	}

	handler, err := NewCreateTaskHandler()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error instantiating handler: %v", err)
	}

	output, err := handler.Handle(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating task: %v", err)
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

// handleUpdateTaskStatus handles POST requests to update the status of a task
func handleUpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	// Extract task ID from URL
	taskIDStr := r.URL.Path[7:] // Remove "/tasks/" prefix
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid task ID: %v", err)
		return
	}

	// Parse the request body
	var input UpdateTaskStatusInput
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing request body: %v", err)
		return
	}
	input.TaskID = taskID

	handler, err := NewUpdateTaskStatusHandler()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error instantiating handler: %v", err)
	}

	output, err := handler.Handle(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error updating task: %v", err)
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

// handleGetTasks handles GET requests to retrieve all tasks
func handleGetTasks(w http.ResponseWriter) {

	handler, err := NewGetTasksHandler()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error instantiating handler: %v", err)
	}

	output, err := handler.Handle()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error fetching tasks: %v", err)
		return
	}

	// Set content type and encode response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
