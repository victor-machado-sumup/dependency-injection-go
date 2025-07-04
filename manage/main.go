package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

//go:embed tasks.sql
var sqlContent string

func main() {

	command := "migrate"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	fmt.Printf("Running %s command\n", command)

	// Database connection parameters (same as in main.go)
	connStr := "host=localhost port=5433 user=postgres password=postgres dbname=postgres sslmode=disable"

	// Connect to the database
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	if command == "migrate" {
		// Execute the SQL
		_, err = conn.Exec(context.Background(), string(sqlContent))
		if err != nil {
			log.Fatalf("Error executing SQL: %v\n", err)
		}

		fmt.Println("Migration completed successfully!")
	} else if command == "clean" {
		// Execute the clean command to delete all tasks
		_, err = conn.Exec(context.Background(), "DELETE FROM tasks")
		if err != nil {
			log.Fatalf("Error clearing tasks table: %v\n", err)
		}

		fmt.Println("Tasks table cleared successfully!")
	}

}
