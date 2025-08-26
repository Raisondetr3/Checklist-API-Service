package model

import (
	"fmt"
	"time"
)

type Task struct {
	ID          string
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTask(title, description string) *Task {
	return &Task{
		Title:       title,
		Description: description,
		Completed:   false,
	}
}

func (t *Task) Validate() error {
	if t.Title == "" {
		return fmt.Errorf("title is required")
	}
	
	if len(t.Title) > 255 {
		return fmt.Errorf("title too long (max 255 characters)")
	}
	
	if len(t.Description) > 1000 {
		return fmt.Errorf("description too long (max 1000 characters)")
	}
	
	return nil
}

func (t *Task) IsCompleted() bool {
	return t.Completed
}

func (t *Task) MarkCompleted() {
	t.Completed = true
}

func (t *Task) MarkIncomplete() {
	t.Completed = false
}

func (t *Task) Update(title, description *string, completed *bool) {
	if title != nil {
		t.Title = *title
	}
	if description != nil {
		t.Description = *description
	}
	if completed != nil {
		t.Completed = *completed
	}
}