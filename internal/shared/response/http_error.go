package response

import (
	"errors"
	"log"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/apperror"
)

// ErrorMapping defines HTTP status and code for a domain error
type ErrorMapping struct {
	Status int
	Code   string
}

// HTTPErrorResponse represents the structured HTTP error response
type HTTPErrorResponse struct {
	StatusCode int                    `json:"status_code"`
	RequestID  string                 `json:"request_id"`
	Timestamp  string                 `json:"timestamp"`
	Path       string                 `json:"path"`
	Error      string                 `json:"error"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

// DomainError interface for errors with metadata
type DomainError interface {
	error
	GetCode() string
	GetCategory() interface{} // ErrorCategory from domain layer
}

// Global error registry - modules register their errors here
var (
	errorMappings = make(map[error]ErrorMapping)
	registryMu    sync.RWMutex
)

// RegisterErrorMapping allows modules to register their domain errors
// This should be called during module initialization
func RegisterErrorMapping(err error, status int, code string) {
	registryMu.Lock()
	defer registryMu.Unlock()
	errorMappings[err] = ErrorMapping{
		Status: status,
		Code:   code,
	}
}

// RegisterErrorMappings allows bulk registration of error mappings
func RegisterErrorMappings(mappings map[error]ErrorMapping) {
	registryMu.Lock()
	defer registryMu.Unlock()
	for err, mapping := range mappings {
		errorMappings[err] = mapping
	}
}

// HandleError maps domain errors and apperror to HTTP status codes
func HandleError(c *gin.Context, err error) {
	// Get or generate request ID
	requestID := getRequestID(c)

	// Default values
	status := http.StatusInternalServerError
	code := "INTERNAL_ERROR"
	message := "an unexpected error occurred"
	var details map[string]interface{}

	// 1. Check for apperror (shared errors from infrastructure)
	if appErr, ok := err.(*apperror.Error); ok {
		status := mapAppErrorToStatus(appErr)
		respondWithError(c, HTTPErrorResponse{
			StatusCode: status,
			RequestID:  requestID,
			Timestamp:  time.Now().UTC().Format(time.RFC3339),
			Path:       c.Request.URL.Path,
			Error:      appErr.Code,
			Message:    appErr.Message,
			Details:    convertDetails(appErr.Details),
		}, status)
		logError(c, appErr, status)
		return
	}

	// 2. Check for validator validation errors
	if handleValidationError(c, err, requestID) {
		return
	}

	// 3. Check if it's a domain error with metadata
	if domainErr, ok := err.(DomainError); ok {
		code = domainErr.GetCode()
		message = domainErr.Error()
		status = categoryToHTTPStatus(domainErr.GetCategory())
		// Extract details if available using reflection
		details = extractErrorDetails(err)
	} else if mapping, found := findErrorMapping(err); found {
		// 4. Fallback to registered error mappings
		status = mapping.Status
		code = mapping.Code
		message = err.Error()
	} else {
		// Unknown error - use default 500
		message = err.Error()
	}

	respondWithError(c, HTTPErrorResponse{
		StatusCode: status,
		RequestID:  requestID,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Path:       c.Request.URL.Path,
		Error:      code,
		Message:    message,
		Details:    details,
	}, status)

	logError(c, err, status)
}

// findErrorMapping looks up the error in the registry
func findErrorMapping(err error) (ErrorMapping, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()

	for registeredErr, mapping := range errorMappings {
		if errors.Is(err, registeredErr) {
			return mapping, true
		}
	}
	return ErrorMapping{}, false
}

// handleValidationError handles go-playground/validator errors
func handleValidationError(c *gin.Context, err error, requestID string) bool {
	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return false
	}

	// Build field-specific errors
	fieldErrors := make(map[string]interface{})
	for _, fieldErr := range validationErrs {
		fieldErrors[fieldErr.Field()] = map[string]string{
			"tag":   fieldErr.Tag(),
			"value": fieldErr.Param(),
		}
	}

	respondWithError(c, HTTPErrorResponse{
		StatusCode: http.StatusBadRequest,
		RequestID:  requestID,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Path:       c.Request.URL.Path,
		Error:      "VALIDATION_ERROR",
		Message:    "validation failed",
		Details:    fieldErrors,
	}, http.StatusBadRequest)

	logError(c, err, http.StatusBadRequest)
	return true
}

// categoryToHTTPStatus maps error category to HTTP status code
func categoryToHTTPStatus(category interface{}) int {
	categoryStr := ""
	if cat, ok := category.(interface{ String() string }); ok {
		categoryStr = cat.String()
	} else if cat, ok := category.(string); ok {
		categoryStr = cat
	}

	mapping := map[string]int{
		"VALIDATION":    http.StatusBadRequest,
		"NOT_FOUND":     http.StatusNotFound,
		"CONFLICT":      http.StatusConflict,
		"UNAUTHORIZED":  http.StatusUnauthorized,
		"FORBIDDEN":     http.StatusForbidden,
		"BUSINESS_RULE": http.StatusUnprocessableEntity,
	}

	if status, found := mapping[categoryStr]; found {
		return status
	}
	return http.StatusInternalServerError
}

// mapAppErrorToStatus maps apperror.Error types to HTTP status codes
func mapAppErrorToStatus(err *apperror.Error) int {
	switch err.Type {
	case apperror.ErrorTypeNotFound:
		return http.StatusNotFound
	case apperror.ErrorTypeValidation:
		return http.StatusBadRequest
	case apperror.ErrorTypeConflict:
		return http.StatusConflict
	case apperror.ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case apperror.ErrorTypeForbidden:
		return http.StatusForbidden
	case apperror.ErrorTypeBadRequest:
		return http.StatusBadRequest
	case apperror.ErrorTypeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// respondWithError sends the error response
func respondWithError(c *gin.Context, errResp HTTPErrorResponse, status int) {
	c.JSON(status, errResp)
}

// logError logs the error with context
func logError(c *gin.Context, err error, status int) {
	log.Printf("[ERROR] [%d] [%s] %s - %v",
		status,
		c.Request.Method,
		c.Request.URL.Path,
		err,
	)
}

// getRequestID retrieves or generates a request ID
func getRequestID(c *gin.Context) string {
	if id, exists := c.Get("request_id"); exists {
		if requestID, ok := id.(string); ok {
			return requestID
		}
	}
	// Generate new request ID if not present
	requestID := "req_" + uuid.New().String()[:8]
	c.Set("request_id", requestID)
	return requestID
}

// convertDetails converts interface{} to map[string]interface{}
func convertDetails(details interface{}) map[string]interface{} {
	if details == nil {
		return nil
	}
	if m, ok := details.(map[string]interface{}); ok {
		return m
	}
	return map[string]interface{}{"data": details}
}

// extractErrorDetails tries to extract Details field from domain errors
func extractErrorDetails(err error) map[string]interface{} {
	// Import domain error types and check each one
	// Since we can't import domain packages in shared/, we use a simple approach
	// Domain errors that implement the DomainError interface can provide their details
	// through their public Details field which we access via reflection pattern

	// Simply try to get the Details field if it exists on the struct
	// This is safe because Go will return the zero value if the field doesn't exist
	v := reflect.ValueOf(err)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		detailsField := v.FieldByName("Details")
		if detailsField.IsValid() && detailsField.CanInterface() {
			if details, ok := detailsField.Interface().(map[string]interface{}); ok {
				return details
			}
		}
	}

	return nil
}

// Legacy support for AppError (kept for backward compatibility)
type AppError struct {
	HttpStatus int         `json:"-"`
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
}

func (e AppError) Error() string {
	return e.Message
}

func NewAppError(httpStatus int, code, message string, details interface{}) error {
	return AppError{
		HttpStatus: httpStatus,
		Code:       code,
		Message:    message,
		Details:    details,
	}
}

func GetStatusCode(err error) int {
	if appErr, ok := err.(AppError); ok {
		return appErr.HttpStatus
	}

	// Check for new apperror.Error
	if appErr, ok := err.(*apperror.Error); ok {
		return mapAppErrorToStatus(appErr)
	}

	return http.StatusInternalServerError
}
