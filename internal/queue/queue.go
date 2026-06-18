package queue

import (
	"errors"
)

var ErrQueueFull = errors.New("queue is full")

type Queue struct {
	ch chan string
}

// NewQueue initializes a queue with a maximum capacity.
func NewQueue(capacity int) *Queue {
	return &Queue{
		ch: make(chan string, capacity),
	}
}

// Push adds a job ID to the queue. Non-blocking with an error if full.
func (q *Queue) Push(jobID string) error {
	select {
	case q.ch <- jobID:
		return nil
	default:
		return ErrQueueFull
	}
}

// Pop blocks until a job ID is available and returns it.
func (q *Queue) Pop() string {
	return <-q.ch
}
