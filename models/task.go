package models

import "time"

type Task struct {
	ID        int
	ProjectID int
	Status    string
	Message   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
