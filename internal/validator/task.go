package validator

import (
	"errors"
	"strings"

	"github.com/Raisondetr3/checklist-api-service/pkg/dto"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return e.Message
}

type ValidationErrors []ValidationError

func (errs ValidationErrors) Error() string {
	var messages []string
	for _, err := range errs {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

func (errs ValidationErrors) HasErrors() bool {
	return len(errs) > 0
}

func ValidateCreateTaskRequest(req dto.CreateTaskRequest) ValidationErrors {
	var errs ValidationErrors

	if req.Title == "" {
		errs = append(errs, ValidationError{
			Field:   "title",
			Message: "Title is required",
		})
	} else if len(req.Title) > 255 {
		errs = append(errs, ValidationError{
			Field:   "title", 
			Message: "Title is too long (max 255 characters)",
		})
	}

	if len(req.Description) > 1000 {
		errs = append(errs, ValidationError{
			Field:   "description",
			Message: "Description is too long (max 1000 characters)",
		})
	}

	return errs
}

func ValidateUpdateTaskRequest(req dto.UpdateTaskRequest) ValidationErrors {
	var errs ValidationErrors

	if req.Title == nil && req.Description == nil && req.Completed == nil {
		errs = append(errs, ValidationError{
			Field:   "request",
			Message: "At least one field must be provided for update",
		})
		return errs
	}

	if req.Title != nil {
		if *req.Title == "" {
			errs = append(errs, ValidationError{
				Field:   "title",
				Message: "Title cannot be empty",
			})
		} else if len(*req.Title) > 255 {
			errs = append(errs, ValidationError{
				Field:   "title",
				Message: "Title is too long (max 255 characters)",
			})
		}
	}

	if req.Description != nil && len(*req.Description) > 1000 {
		errs = append(errs, ValidationError{
			Field:   "description",
			Message: "Description is too long (max 1000 characters)",
		})
	}

	return errs
}

func ValidateTaskID(taskID string) error {
	if taskID == "" {
		return errors.New("Task ID is required")
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
		return nil, errors.New("Invalid 'completed' parameter. Use 'true' or 'false'")
	}
}