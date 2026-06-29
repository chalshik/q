package grpc

import (
	"context"

	proto "dispatcher/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const workerAddr = "localhost:50052"

type WorkerClient struct {
	client proto.WorkerServiceClient
}

func NewWorkerClient() (*WorkerClient, error) {
	conn, err := grpc.NewClient(workerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &WorkerClient{client: proto.NewWorkerServiceClient(conn)}, nil
}

func (w *WorkerClient) PushJob(ctx context.Context, jobID, prompt, priority string) (*proto.PushJobResponse, error) {
	return w.client.PushJob(ctx, &proto.PushJobRequest{
		JobId:    jobID,
		Prompt:   prompt,
		Priority: priority,
	})
}

func (w *WorkerClient) HealthCheck(ctx context.Context, workerID string) (*proto.HealthCheckResponse, error) {
	return w.client.HealthCheck(ctx, &proto.HealthCheckRequest{
		WorkerId: workerID,
	})
}
