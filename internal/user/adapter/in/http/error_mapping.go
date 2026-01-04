package http

import (
	"net/http"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/response"
	usererror "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/error"
)

// init registers user domain errors to HTTP status codes
// This is in the adapter layer, so it's OK to know about HTTP
func init() {
	response.RegisterErrorMappings(map[error]response.ErrorMapping{
		// Authentication errors
		usererror.ErrUserNotFound:      {Status: http.StatusNotFound, Code: "USER_NOT_FOUND"},
		usererror.ErrEmailAlreadyUsed:  {Status: http.StatusConflict, Code: "USER_ALREADY_EXISTS"},
		usererror.ErrUserAlreadyExists: {Status: http.StatusConflict, Code: "USER_ALREADY_EXISTS"},
		usererror.ErrInvalidCredential: {Status: http.StatusUnauthorized, Code: "INVALID_CREDENTIALS"},

		// Validation errors
		usererror.ErrInvalidEmail:  {Status: http.StatusBadRequest, Code: "VALIDATION_ERROR"},
		usererror.ErrWeakPassword:  {Status: http.StatusBadRequest, Code: "VALIDATION_ERROR"},
		usererror.ErrEmptyEmail:    {Status: http.StatusBadRequest, Code: "VALIDATION_ERROR"},
		usererror.ErrEmptyPassword: {Status: http.StatusBadRequest, Code: "VALIDATION_ERROR"},

		// Business rule errors
		usererror.ErrUserInactive: {Status: http.StatusForbidden, Code: "USER_INACTIVE"},
	})
}
