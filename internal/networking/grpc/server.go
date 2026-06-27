package grpc

import (
	"dispatcher/internal/queue"
	service "dispatcher/internal/services"
	"fmt"

	"context"

	proto "dispatcher/proto"

	"github.com/jmoiron/sqlx"
)

type DispatcherServer struct {
	jobService *service.JobService
	db         *sqlx.DB
	queue      *queue.Queue
	proto.UnimplementedDispatcherServer
}

func NewDispatcherServer(jobService *service.JobService, db *sqlx.DB, queue *queue.Queue) *DispatcherServer {
	return &DispatcherServer{
		jobService: jobService,
		db:         db,
		queue:      queue,
	}
}
func (s *DispatcherServer) GetJob(ctx context.Context, req *proto.GetJobRequest) (*proto.GetJobResponse, error) {
	jobId := s.queue.Pop()
	job, err := s.jobService.GetByID(ctx, jobId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve job: %w", err)
	}
	return &proto.GetJobResponse{Id: job.ID, Prompt: job.Prompt, Status: job.Status}, nil
}
func (s *DispatcherServer) SubmitResult(ctx context.Context, req *proto.SubmitResultRequest) (*proto.SubmitResultResponse, error) {
	err := s.jobService.UpdateStatus(ctx, req.JobId, req.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to update job status: %w", err)
	}
	return &proto.SubmitResultResponse{Success: true}, nil
}
