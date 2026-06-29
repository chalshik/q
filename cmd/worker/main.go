package main

import (
	"context"
	"log"

	workergrpc "dispatcher/cmd/worker/internal/networking/gprc"
	"dispatcher/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Start inbound server first
	go func() {
		if err := workergrpc.ServeWorkerGRPC(":50052"); err != nil {
			log.Fatalf("Worker server failed: %v", err)
		}
	}()

	// Outbound — connect to dispatcher
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := proto.NewDispatcherServiceClient(conn)

	// Register this worker as available
	_, err = client.RegisterAvailable(context.Background(), &proto.RegisterAvailableRequest{
		WorkerId:       "worker-1",
		Address:        "localhost:50052",
		ModelId:        "model-a",
		AvailableSlots: 4,
	})
	if err != nil {
		log.Fatalf("Failed to register worker: %v", err)
	}

	log.Println("Worker registered successfully")

	// Block forever
	select {}
}
