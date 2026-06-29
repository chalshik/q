package grpc

import (
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
