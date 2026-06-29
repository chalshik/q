package gprc

import (
	"context"
	"log"
	"net"

	"dispatcher/proto"

	"google.golang.org/grpc"
)

type WorkerServer struct {
	proto.UnimplementedWorkerServiceServer
}

func (s *WorkerServer) PushJob(ctx context.Context, req *proto.PushJobRequest) (*proto.PushJobResponse, error) {
	log.Printf("Received job: %s", req.JobId)
	return &proto.PushJobResponse{Accepted: true}, nil
}

func (s *WorkerServer) HealthCheck(ctx context.Context, req *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	return &proto.HealthCheckResponse{
		Alive:          true,
		ModelId:        "model-a",
		AvailableSlots: 4,
	}, nil
}

func NewWorkerGRPCServer() *grpc.Server {
	grpcServer := grpc.NewServer()
	proto.RegisterWorkerServiceServer(grpcServer, &WorkerServer{})
	return grpcServer
}

func ServeWorkerGRPC(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Printf("Worker gRPC server listening on %s", addr)
	return NewWorkerGRPCServer().Serve(lis)
}
