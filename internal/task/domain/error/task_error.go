package error

var (
	// Resource not found errors
	ErrTaskNotFound = NewNotFoundError(
		"TASK_NOT_FOUND",
		"task not found",
	)

	// Conflict errors
	ErrTaskAlreadyExists = NewConflictError(
		"TASK_ALREADY_EXISTS",
		"task already exists",
	)
	ErrTaskAlreadyCompleted = NewConflictError(
		"TASK_ALREADY_COMPLETED",
		"task is already completed",
	)

	// Validation errors
	ErrEmptyTitle = NewValidationError(
		"TASK_EMPTY_TITLE",
		"task title cannot be empty",
	)
	ErrInvalidTaskID = NewValidationError(
		"TASK_INVALID_ID",
		"invalid task ID format",
	)
	ErrEmptyUserID = NewValidationError(
		"TASK_EMPTY_USER_ID",
		"user ID cannot be empty",
	)

	// Authorization errors
	ErrTaskAccessDenied = NewForbiddenError(
		"TASK_ACCESS_DENIED",
		"you don't have permission to access this task",
	)
)
