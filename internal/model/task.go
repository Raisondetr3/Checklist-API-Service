package model

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	id          uuid.UUID
	title       string
	description string
	completed   bool
	created_at  time.Time
}

func NewTask(title, description string) *Task {
	return &Task{
		title: title,
		description: description,
		completed: false,
		// остальные параметры в DB-Service проставляются
	}
}