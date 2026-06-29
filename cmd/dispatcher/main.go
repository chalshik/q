package main

import (
	"context"
	"log"
	"net"
	"net/http"

	transportHTTP "dispatcher/internal/networking/http"
	"dispatcher/internal/persistence"
	"dispatcher/internal/queue"
	"dispatcher/internal/scheduler"
	service "dispatcher/internal/services"

	internalgrpc "dispatcher/internal/networking/grpc"
	proto "dispatcher/proto"

	"google.golang.org/grpc"
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
	jobQueue, err := queue.NewJobQueue("localhost:6379")
	if err != nil {
		log.Fatalf("Redis initialization failed: %v", err)
	}

	// 3. Initialize Core Business Logic (Service Layer)
	jobService := service.NewJobService(jobRepo, jobQueue)

	// 4. Initialize Transport Layer (HTTP)
	handler := transportHTTP.NewHandler(jobService)
	scheduler := scheduler.NewScheduler(jobService, jobQueue)

	// Start the scheduler in a separate goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go scheduler.Start(ctx)
	// 5. Setup Router and Start Server
	mux := http.NewServeMux()
	mux.HandleFunc("/jobs", handler.CreateJob)

	go func() {
		log.Println("HTTP Server listening on :8080...")
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Fatalf("HTTP Server crashed: %v", err)
		}
	}()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on :50051: %v", err)
	}
	grpcServer := grpc.NewServer()
	dispatcherServer := internalgrpc.NewDispatcherServer(jobService, db, jobQueue)

	proto.RegisterDispatcherServer(grpcServer, dispatcherServer)
	log.Println("gRPC Server listening on :50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC Server crashed: %v", err)
	}
}
