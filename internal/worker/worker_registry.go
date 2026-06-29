package worker

import (
	"dispatcher/internal/queue"
)

type WorkerRegistry struct {
	queue queue.Queue
}

func NewWorkerRegistry(addr string) (*WorkerRegistry, error) {
	client, err := queue.NewRedisClient(addr)
	if err != nil {
		return nil, err
	}
	return &WorkerRegistry{queue: client.NewQueue("workers")}, nil
}

// Register adds a worker ID to the available pool.
func (r *WorkerRegistry) Register(workerID string) error {
	return r.queue.Enqueue(workerID)
}

// Next blocks until a worker is available and returns its ID.
func (r *WorkerRegistry) Next() string {
	return r.queue.Dequeue()
}
