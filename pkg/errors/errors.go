package errors

import (
	"errors"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrTaskNotFound        = errors.New("task not found")
	ErrTaskAlreadyExists   = errors.New("task already exists")
	ErrInvalidInput        = errors.New("invalid input data")
	ErrValidationFailed    = errors.New("validation failed")
	ErrServiceUnavailable  = errors.New("service temporarily unavailable")
	ErrInternalError       = errors.New("internal server error")
)

func HTTPStatusFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if st, ok := status.FromError(err); ok {
		return grpcCodeToHTTPStatus(st.Code())
	}

	switch {
	case isNotFoundError(err):
		return http.StatusNotFound
	case isValidationError(err):
		return http.StatusBadRequest
	case isAlreadyExistsError(err):
		return http.StatusConflict
	case isServiceUnavailableError(err):
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

func MessageFromError(err error) string {
	if err == nil {
		return ""
	}

	if st, ok := status.FromError(err); ok {
		return cleanGRPCMessage(st.Message())
	}

	return err.Error()
}

func grpcCodeToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError
	}
}

func cleanGRPCMessage(msg string) string {
	prefixes := []string{
		"create_task: ",
		"get_task_by_id: ",
		"update_task: ",
		"delete_task: ",
		"list_tasks: ",
		"failed to create task: ",
		"failed to get task: ",
		"failed to update task: ",
		"failed to delete task: ",
		"failed to list tasks: ",
	}

	cleaned := msg
	for _, prefix := range prefixes {
		cleaned = strings.TrimPrefix(cleaned, prefix)
	}

	return cleaned
}

func isNotFoundError(err error) bool {
	return errors.Is(err, ErrTaskNotFound) ||
		strings.Contains(err.Error(), "not found") ||
		strings.Contains(err.Error(), "task not found")
}

func isValidationError(err error) bool {
	return errors.Is(err, ErrValidationFailed) ||
		errors.Is(err, ErrInvalidInput) ||
		strings.Contains(err.Error(), "validation failed") ||
		strings.Contains(err.Error(), "invalid") ||
		strings.Contains(err.Error(), "required") ||
		strings.Contains(err.Error(), "cannot be empty")
}

func isAlreadyExistsError(err error) bool {
	return errors.Is(err, ErrTaskAlreadyExists) ||
		strings.Contains(err.Error(), "already exists")
}

func isServiceUnavailableError(err error) bool {
	return errors.Is(err, ErrServiceUnavailable) ||
		strings.Contains(err.Error(), "unavailable") ||
		strings.Contains(err.Error(), "connection") ||
		strings.Contains(err.Error(), "timeout")
}