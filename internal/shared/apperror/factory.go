package apperror

func NotFound(msg string) *Error {
	return &Error{
		Type:    ErrorTypeNotFound,
		Code:    string(ErrorTypeNotFound),
		Message: msg,
	}
}

func Validation(msg string) *Error {
	return &Error{
		Type:    ErrorTypeValidation,
		Code:    string(ErrorTypeValidation),
		Message: msg,
	}
}

func Internal(msg string) *Error {
	return &Error{
		Type:    ErrorTypeInternal,
		Code:    string(ErrorTypeInternal),
		Message: msg,
	}
}

func Conflict(msg string) *Error {
	return &Error{
		Type:    ErrorTypeConflict,
		Code:    string(ErrorTypeConflict),
		Message: msg,
	}
}

func Unauthorized(msg string) *Error {
	return &Error{
		Type:    ErrorTypeUnauthorized,
		Code:    string(ErrorTypeUnauthorized),
		Message: msg,
	}
}

func Forbidden(msg string) *Error {
	return &Error{
		Type:    ErrorTypeForbidden,
		Code:    string(ErrorTypeForbidden),
		Message: msg,
	}
}

func BadRequest(msg string) *Error {
	return &Error{
		Type:    ErrorTypeBadRequest,
		Code:    string(ErrorTypeBadRequest),
		Message: msg,
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, msg string) *Error {
	if appErr, ok := err.(*Error); ok {
		return &Error{
			Type:    appErr.Type,
			Code:    appErr.Code,
			Message: msg + ": " + appErr.Message,
			Details: appErr.Details,
		}
	}
	return Internal(msg + ": " + err.Error())
}
