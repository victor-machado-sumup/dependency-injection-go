package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
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
			case r.URL.Path == "/tasks" && r.Method == "GET":
				handleGetTasks(w)
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

// Task represents a task in our system
type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
}

// TaskStatus represents the possible status values for a task
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusCompleted TaskStatus = "completed"
)

// handleCreateTask handles POST requests to create new tasks
func handleCreateTask(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var task Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing request body: %v", err)
		return
	}

	// Validate required fields
	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Title is required")
		return
	}

	// Database connection parameters
	connStr := "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error connecting to the database: %v", err)
	}

	// Defer closing the database connection
	defer func() {
		if err := conn.Close(context.Background()); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Insert the task into the database
	query := `INSERT INTO tasks (title, description) VALUES ($1, $2) RETURNING id`
	var id int
	err = conn.QueryRow(context.Background(), query, task.Title, task.Description).Scan(&id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating task: %v", err)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"task": Task{
			ID:          id,
			Title:       task.Title,
			Description: task.Description,
			Status:      TaskStatusPending,
		},
	})
}

// UpdateTaskStatus represents the request body for updating a task's status
type UpdateTaskStatus struct {
	Status TaskStatus `json:"status"`
}

// handleUpdateTaskStatus handles POST requests to update the status of a task
func handleUpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	// Extract task ID from URL
	taskID, err := strconv.Atoi(r.URL.Path[7:]) // Remove "/tasks/" prefix
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid taskID: %v", err)
		return
	}

	// Parse the request body
	var updateReq UpdateTaskStatus
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updateReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing request body: %v", err)
		return
	}

	// Validate status value
	if updateReq.Status != TaskStatusPending && updateReq.Status != TaskStatusCompleted {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid status value. Must be either 'pending' or 'completed'")
		return
	}

	// Database connection parameters
	connStr := "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error connecting to the database: %v", err)
		return
	}

	// Defer closing the database connection
	defer func() {
		if err := conn.Close(context.Background()); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Update task status in the database
	result, err := conn.Exec(context.Background(),
		"UPDATE tasks SET status = $1 WHERE id = $2",
		updateReq.Status,
		taskID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error updating task status: %v", err)
		return
	}

	if result.RowsAffected() == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Task with ID %v not found", taskID)
		return
	}

	// Fetch the updated task from the database
	var updatedTask Task
	err = conn.QueryRow(context.Background(),
		"SELECT id, title, description, status FROM tasks WHERE id = $1",
		taskID).Scan(&updatedTask.ID, &updatedTask.Title, &updatedTask.Description, &updatedTask.Status)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error fetching updated task: %v", err)
		return
	}

	// Return success response with updated task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"task": updatedTask,
	})
}

// handleGetTasks handles GET requests to retrieve all tasks
func handleGetTasks(w http.ResponseWriter) {
	// Database connection parameters
	connStr := "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error connecting to the database: %v", err)
		return
	}

	// Defer closing the database connection
	defer func() {
		if err := conn.Close(context.Background()); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Query all tasks from the database
	rows, err := conn.Query(context.Background(),
		"SELECT id, title, description, status FROM tasks ORDER BY id")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error fetching tasks: %v", err)
		return
	}
	defer rows.Close()

	// Create a slice to hold all tasks
	tasks := []Task{}

	// Iterate through the rows
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error scanning task row: %v", err)
			return
		}
		tasks = append(tasks, task)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error iterating task rows: %v", err)
		return
	}

	// Set content type and encode response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tasks": tasks,
	})
}
