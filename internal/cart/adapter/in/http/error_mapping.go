package http

import (
	"net/http"

	carterror "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/domain/error"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/response"
)

// init registers cart domain errors to HTTP status codes
// This is in the adapter layer, so it's OK to know about HTTP
func init() {
	response.RegisterErrorMappings(map[error]response.ErrorMapping{
		// Resource errors
		carterror.ErrCartNotFound:      {Status: http.StatusNotFound, Code: "CART_NOT_FOUND"},
		carterror.ErrCartAlreadyExists: {Status: http.StatusConflict, Code: "CART_ALREADY_EXISTS"},

		// Validation errors
		carterror.ErrEmptyTitle:    {Status: http.StatusBadRequest, Code: "VALIDATION_ERROR"},
		carterror.ErrInvalidCartID: {Status: http.StatusBadRequest, Code: "VALIDATION_ERROR"},
		carterror.ErrEmptyUserID:   {Status: http.StatusBadRequest, Code: "VALIDATION_ERROR"},

		// Authorization errors
		carterror.ErrCartAccessDenied: {Status: http.StatusForbidden, Code: "CART_ACCESS_DENIED"},
	})
}
