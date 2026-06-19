package persistence

import (
	"context"
	"dispatcher/internal/models"

	"github.com/jmoiron/sqlx"
)

type JobRepository struct {
	db *sqlx.DB
}

func NewJobRepository(db *sqlx.DB) *JobRepository {
	return &JobRepository{db: db}
}

// Insert adds a new job record to PostgreSQL
func (r *JobRepository) Insert(ctx context.Context, job *models.Job) error {
	query := `
		INSERT INTO jobs (id, user_id, prompt, status, created_at, updated_at)
		VALUES (:id, :user_id, :prompt, :status, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, job)
	return err
}

// GetByID fetches a single job by its ID.
func (r *JobRepository) GetByID(ctx context.Context, id string) (*models.Job, error) {
	var job models.Job
	query := `SELECT id, user_id, prompt, status, created_at, updated_at FROM jobs WHERE id = $1`
	if err := r.db.GetContext(ctx, &job, query, id); err != nil {
		return nil, err
	}
	return &job, nil
}

// ListByUserID fetches all jobs belonging to a given user, most recent first.
func (r *JobRepository) ListByUserID(ctx context.Context, userID string) ([]models.Job, error) {
	var jobs []models.Job
	query := `SELECT id, user_id, prompt, status, created_at, updated_at FROM jobs WHERE user_id = $1 ORDER BY created_at DESC`
	if err := r.db.SelectContext(ctx, &jobs, query, userID); err != nil {
		return nil, err
	}
	return jobs, nil
}

// UpdateStatus updates a job's status and refreshes its updated_at timestamp.
func (r *JobRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE jobs SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

// Delete removes a job by its ID.
func (r *JobRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM jobs WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
