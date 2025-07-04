package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sumup/dependency-injection-go/server-dependency-injection/handlers"
	"github.com/sumup/dependency-injection-go/server-dependency-injection/repository"
)

func main() {
	port := ":8080"

	// Create a new Gin router with default middleware
	router := gin.Default()

	// Enable CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Declare dependencies
	pool, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to create connection pool: %v", err)
	}
	repository := repository.NewRepository(pool)

	createTaskHandler := handlers.NewCreateTaskHandler(repository)

	// Define routes
	router.GET("/health", handleHealth)

	router.GET("/tasks", handleGetTasks(repository))
	router.POST("/tasks", handleCreateTask(createTaskHandler))
	router.POST("/tasks/:id", provide(handleUpdateTaskStatus, repository))

	// Start the server
	fmt.Printf("Server starting on port %s\n", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// handleHealth processes the health check request
func handleHealth(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

// handleGetTasks handles GET requests to retrieve all tasks
func handleGetTasks(repository repository.IRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler := handlers.NewGetTasksHandler(repository)

		output, err := handler.Handle()
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching tasks: %v", err))
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

// handleCreateTask handles POST requests to create new tasks
func handleCreateTask(handler *handlers.CreateTaskHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input handlers.CreateTaskInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Error parsing request body: %v", err))
			return
		}

		output, err := handler.Handle(input)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error creating task: %v", err))
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}

// handleUpdateTaskStatus handles POST requests to update the status of a task
func handleUpdateTaskStatus(c *gin.Context, repository repository.IRepository) {
	var input handlers.UpdateTaskStatusInput
	if err := c.ShouldBindUri(&input); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid task ID: %v", err))
		return
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Error parsing request body: %v", err))
		return
	}

	handler := handlers.NewUpdateTaskStatusHandler(repository)

	output, err := handler.Handle(input)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error updating task: %v", err))
		return
	}

	c.JSON(http.StatusCreated, output)
}

type HandleFunc func(c *gin.Context, repository repository.IRepository)

func provide(handleFunc HandleFunc, repository repository.IRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		handleFunc(c, repository)
	}
}
