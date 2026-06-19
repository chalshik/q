package main

import (
	"context"
	"dispatcher/proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}
	client := proto.NewDispatcherClient(conn)
	for {
		jobResponse, err := client.GetJob(context.Background(), &proto.GetJobRequest{})
		if err != nil {
			panic(err)
		}
		log.Printf("Received Job: %+v\n", jobResponse)

		submitResponse, err := client.SubmitResult(context.Background(), &proto.SubmitResultRequest{
			JobId:  jobResponse.Id,
			Status: "completed",
		})
		if err != nil {
			panic(err)
		}
		log.Printf("Submit Result Response: %+v\n", submitResponse)
	}
}
