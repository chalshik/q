package service

import (
	"context"
	"fmt"
	"time"

	"dispatcher/internal/models"
	"dispatcher/internal/persistence"
	"dispatcher/internal/queue"

	"github.com/google/uuid"
)

type JobService struct {
	repo  *persistence.JobRepository
	queue *queue.Queue
}

func NewJobService(repo *persistence.JobRepository, q *queue.Queue) *JobService {
	return &JobService{
		repo:  repo,
		queue: q,
	}
}

func (s *JobService) Create(ctx context.Context, userID string, prompt string) (*models.Job, error) {
	job := &models.Job{
		ID:        uuid.New().String(),
		UserID:    userID,
		Prompt:    prompt,
		Status:    "queued",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 1. Delegate data management to the persistence layer
	if err := s.repo.Insert(ctx, job); err != nil {
		return nil, fmt.Errorf("service failed to save job: %w", err)
	}

	// 2. Push to the in-memory queue
	if err := s.queue.Push(job.ID); err != nil {
		return nil, fmt.Errorf("service failed to enqueue job: %w", err)
	}

	return job, nil
}
func (s *JobService) GetByID(ctx context.Context, id string) (*models.Job, error) {
	job, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service failed to retrieve job: %w", err)
	}
	return job, nil
}
func (s *JobService) UpdateStatus(ctx context.Context, id string, status string) error {
	return s.repo.UpdateStatus(ctx, id, status)
}
