package apperror

func NotFound(msg string) *Error {
	return &Error{
		Code:    "NOT_FOUND",
		Message: msg,
	}
}

func Validation(msg string) *Error {
	return &Error{
		Code:    "VALIDATION_ERROR",
		Message: msg,
	}
}

func Internal(msg string) *Error {
	return &Error{
		Code:    "INTERNAL_ERROR",
		Message: msg,
	}
}
