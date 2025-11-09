package validator

import (
	"errors"

	"github.com/Raisondetr3/checklist-api-service/pkg/dto"
)

var (
	ErrTitleRequired             = errors.New("Title is required")
	ErrTitleTooLong              = errors.New("Title is too long (max 255 characters)")
	ErrTitleEmpty                = errors.New("Title cannot be empty")
	ErrDescriptionTooLong        = errors.New("Description is too long (max 1000 characters)")
	ErrTaskIDRequired            = errors.New("Task ID is required")
	ErrNoFieldsProvided          = errors.New("At least one field must be provided for update")
	ErrInvalidCompletedParameter = errors.New("Invalid 'completed' parameter. Use 'true' or 'false'")

	MaxTitleLength       = 255
	MaxDescriptionLength = 1000
)

func ValidateCreateTaskRequest(req dto.CreateTaskRequest) error {
	if req.Title == "" {
		return ErrTitleRequired
	}

	if len(req.Title) > MaxTitleLength {
		return ErrTitleTooLong
	}

	if len(req.Description) > MaxDescriptionLength {
		return ErrDescriptionTooLong
	}

	return nil
}

func ValidateUpdateTaskRequest(req dto.UpdateTaskRequest) error {
	if req.Title == nil && req.Description == nil && req.Completed == nil {
		return ErrNoFieldsProvided
	}

	if req.Title != nil {
		if *req.Title == "" {
			return ErrTitleEmpty
		}
		if len(*req.Title) > MaxTitleLength {
			return ErrTitleTooLong
		}
	}

	if req.Description != nil && len(*req.Description) > MaxDescriptionLength {
		return ErrDescriptionTooLong
	}

	return nil
}

func ValidateTaskID(taskID string) error {
	if taskID == "" {
		return ErrTaskIDRequired
	}

	return nil
}

func ValidateCompletedParam(completedStr string) (*bool, error) {
	if completedStr == "" {
		return nil, nil
	}

	switch completedStr {
	case "true":
		completed := true
		return &completed, nil
	case "false":
		completed := false
		return &completed, nil
	default:
		return nil, ErrInvalidCompletedParameter
	}
}
