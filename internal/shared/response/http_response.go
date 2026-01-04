package response

import (
	"time"

	"github.com/gin-gonic/gin"
)

type HTTPSuccessResponse struct {
	StatusCode int         `json:"status_code"`
	RequestID  string      `json:"request_id"`
	Timestamp  string      `json:"timestamp"`
	Path       string      `json:"path"`
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

// Legacy Response struct (kept for backward compatibility)
type Response struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func NewResponse(status int, success bool, message string, data interface{}, errors interface{}) Response {
	return Response{
		Status:  status,
		Success: success,
		Message: message,
		Data:    data,
		Errors:  errors,
	}
}

func ErrorResponse(c *gin.Context, err error) {
	statusCode := GetStatusCode(err)
	response := NewResponse(statusCode, false, "failed", nil, err)
	c.JSON(statusCode, response)
}

// SuccessResponse creates a standardized success response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	requestID := getRequestID(c)

	response := HTTPSuccessResponse{
		StatusCode: statusCode,
		RequestID:  requestID,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Path:       c.Request.URL.Path,
		Success:    true,
		Message:    message,
		Data:       data,
	}

	c.JSON(statusCode, response)
}
