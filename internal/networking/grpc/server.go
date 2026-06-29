package grpc

import (
	"context"
	"dispatcher/internal/worker"
	proto "dispatcher/proto"
)

type DispatcherServer struct {
	registry *worker.WorkerRegistry

	proto.UnimplementedDispatcherServiceServer
}

func NewDispatcherServer(registry *worker.WorkerRegistry) *DispatcherServer {
	return &DispatcherServer{
		registry: registry,
	}
}
func (s *DispatcherServer) RegisterAvailable(ctx context.Context, req *proto.RegisterAvailableRequest) (*proto.RegisterAvailableResponse, error) {
	if err := s.registry.Register(req.WorkerId); err != nil {
		return &proto.RegisterAvailableResponse{
			Accepted: false,
		}, err
	}

	return &proto.RegisterAvailableResponse{
		Accepted: true,
	}, nil
}
