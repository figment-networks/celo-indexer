package model

import "time"

// JobID uniquely identifies a job
type JobID uint64

// Job represents a worker job
type Job struct {
	ID         JobID      `json:"-"`
	Height     *int64     `json:"height"`
	RunCount   uint64     `json:"run_count"`
	LastError  *string    `json:"last_error"`
	StartedAt  *time.Time `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"-"`
}
