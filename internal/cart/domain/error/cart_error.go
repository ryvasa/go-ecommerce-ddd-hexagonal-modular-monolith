package error

var (
	// Resource not found errors
	ErrCartNotFound = NewNotFoundError(
		"CART_NOT_FOUND",
		"cart not found",
	)

	// Conflict errors
	ErrCartAlreadyExists = NewConflictError(
		"CART_ALREADY_EXISTS",
		"cart already exists",
	)

	// Validation errors
	ErrEmptyTitle = NewValidationError(
		"CART_EMPTY_TITLE",
		"cart title cannot be empty",
	)
	ErrInvalidCartID = NewValidationError(
		"CART_INVALID_ID",
		"invalid cart ID format",
	)
	ErrEmptyUserID = NewValidationError(
		"CART_EMPTY_USER_ID",
		"user ID cannot be empty",
	)

	// Authorization errors
	ErrCartAccessDenied = NewForbiddenError(
		"CART_ACCESS_DENIED",
		"you don't have permission to access this cart",
	)
)
