package main

import (
	"context"
	"dispatcher/proto"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}
	client := proto.NewDispatcherClient(conn)

	jobResponse, err := client.GetJob(context.Background(), &proto.GetJobRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Received Job: %+v\n", jobResponse)

	// Example: Submit a result
	submitResponse, err := client.SubmitResult(context.Background(), &proto.SubmitResultRequest{
		JobId:  jobResponse.Id,
		Status: "completed",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Submit Result Response: %+v\n", submitResponse)
}
