package error

var (
	// Authentication errors
	ErrEmailAlreadyUsed = NewConflictError(
		"USER_EMAIL_EXISTS",
		"email already used",
	)
	ErrInvalidCredential = NewUnauthorizedError(
		"USER_INVALID_CREDENTIALS",
		"invalid credentials",
	)
	ErrUserNotFound = NewNotFoundError(
		"USER_NOT_FOUND",
		"user not found",
	)

	// Validation errors
	ErrInvalidEmail = NewValidationError(
		"USER_INVALID_EMAIL",
		"invalid email format",
	)
	ErrWeakPassword = NewValidationError(
		"USER_WEAK_PASSWORD",
		"password does not meet security requirements",
	)
	ErrEmptyEmail = NewValidationError(
		"USER_EMPTY_EMAIL",
		"email cannot be empty",
	)
	ErrEmptyPassword = NewValidationError(
		"USER_EMPTY_PASSWORD",
		"password cannot be empty",
	)

	// Business rule errors
	ErrUserInactive = NewForbiddenError(
		"USER_INACTIVE",
		"user account is inactive",
	)
	ErrUserAlreadyExists = NewConflictError(
		"USER_ALREADY_EXISTS",
		"user already exists",
	)
)
