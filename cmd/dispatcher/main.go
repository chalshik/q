package main

import (
	"context"
	"log"
	"net/http"

	transportHTTP "dispatcher/internal/http"
	"dispatcher/internal/persistence"
	"dispatcher/internal/queue"
	"dispatcher/internal/scheduler"
	service "dispatcher/internal/services"
)

func main() {
	log.Println("Dispatcher starting...")

	// 1. Connect to Postgres (reads its own config internally)
	db, err := persistence.NewPostgresConnection()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Close()
	log.Println("Database connected successfully.")

	// 2. Initialize Infrastructure Layers
	jobRepo := persistence.NewJobRepository(db)
	jobQueue := queue.NewQueue(1000) // Set a baseline capacity for your memory channel

	// 3. Initialize Core Business Logic (Service Layer)
	jobService := service.NewJobService(jobRepo, jobQueue)

	// 4. Initialize Transport Layer (HTTP)
	handler := transportHTTP.NewHandler(jobService)
	scheduler := scheduler.NewScheduler(jobService, jobQueue)

	// Start the scheduler in a separate goroutine
	go scheduler.Start(context.Background())
	// 5. Setup Router and Start Server
	mux := http.NewServeMux()
	mux.HandleFunc("/jobs", handler.CreateJob)

	log.Println("HTTP Server listening on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("HTTP Server crashed: %v", err)
	}
}
