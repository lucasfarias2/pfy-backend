package models

import (
	"packlify-cloud-backend/models/constants"
	"time"
)

type Task struct {
	ID         int
	ProjectID  int
	TaskTypeID int
	Status     constants.TaskStatus
	Message    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
