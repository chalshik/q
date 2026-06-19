package scheduler

import (
	"context"
	"dispatcher/internal/queue"
	service "dispatcher/internal/services"
	"log"
)

type Scheduler struct {
	jobService *service.JobService
	queue      *queue.Queue
}

func NewScheduler(jobService *service.JobService, queue *queue.Queue) *Scheduler {
	return &Scheduler{
		jobService: jobService,
		queue:      queue,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	log.Println("scheduler: starting main loop")
	for {
		select {
		case <-ctx.Done():
			log.Println("scheduler: received shutdown signal, exiting main loop")
			return
		default:
			jobID := s.queue.Pop()
			log.Printf("scheduler: popped job ID %s from queue", jobID)
			log.Printf("scheduler: processing job ID %s", jobID)
		}
	}
}
