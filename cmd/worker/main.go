package main

import (
	"context"
	"log"

	workergrpc "dispatcher/cmd/worker/internal/networking/gprc"
	"dispatcher/proto"
	"os"

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

	// 1. Fetch and validate the target address
	dispatcherAddr := os.Getenv("DISPATCHER_URL")
	if dispatcherAddr == "" {
		// Fallback to the bridge address out of the Kind cluster to your laptop
		dispatcherAddr = "192.168.0.60:50051"
	}
	log.Printf("Connecting to dispatcher at: %s", dispatcherAddr)

	// Outbound — connect to dispatcher
	conn, err := grpc.NewClient(dispatcherAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := proto.NewDispatcherServiceClient(conn)

	// Register this worker as available
	_, err = client.RegisterAvailable(context.Background(), &proto.RegisterAvailableRequest{
		WorkerId:       "worker-1",
		Address:        "localhost:50052", // Note: you may need to update this to the pod IP later so dispatcher can call back
		ModelId:        "model-a",
		AvailableSlots: 4,
	})
	if err != nil {
		log.Fatalf("Failed to register worker: %v", err)
	}

	log.Println("Worker registered successfully")

	select {}
}
