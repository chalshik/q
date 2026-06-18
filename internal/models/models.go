package models

import "time"

type Job struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Prompt    string    `db:"prompt"`
	Status    string    `db:"status"`
	WorkerID  *string   `db:"worker_id"` // Nullable if not assigned yet
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
