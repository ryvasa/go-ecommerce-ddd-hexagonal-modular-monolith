package http

import (
	"net/http"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/response"
	taskerror "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task/domain/error"
)

// init registers task domain errors to HTTP status codes
// This is in the adapter layer, so it's OK to know about HTTP
func init() {
	response.RegisterErrorMappings(map[error]response.ErrorMapping{
		// Resource errors
		taskerror.ErrTaskNotFound:         {Status: http.StatusNotFound, Code: "TASK_NOT_FOUND"},
		taskerror.ErrTaskAlreadyExists:    {Status: http.StatusConflict, Code: "TASK_CONFLICT"},
		taskerror.ErrTaskAlreadyCompleted: {Status: http.StatusConflict, Code: "TASK_CONFLICT"},

		// Validation errors
		taskerror.ErrEmptyTitle:    {Status: http.StatusBadRequest, Code: "VALIDATION_ERROR"},
		taskerror.ErrInvalidTaskID: {Status: http.StatusBadRequest, Code: "VALIDATION_ERROR"},
		taskerror.ErrEmptyUserID:   {Status: http.StatusBadRequest, Code: "VALIDATION_ERROR"},

		// Authorization errors
		taskerror.ErrTaskAccessDenied: {Status: http.StatusForbidden, Code: "TASK_ACCESS_DENIED"},
	})
}
