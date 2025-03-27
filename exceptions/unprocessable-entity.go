package exceptions

import (
	"fmt"
	"net/http"
)

// UnprocessableEntityException represents a 422 Unprocessable Entity error
// indicating that the server cannot process the request due to semantic issues.
type UnprocessableEntityException struct {
	StatusCode int      `json:"status_code"`
	Message    string   `json:"message"`
	Errors     []string `json:"errors"`
}

// NewUnprocessableEntityException creates a new instance of UnprocessableEntityException.
func NewUnprocessableEntityException(message string, errors []string) *UnprocessableEntityException {
	if message == "" {
		message = "Unprocessable Entity"
	}

	return &UnprocessableEntityException{
		StatusCode: http.StatusUnprocessableEntity,
		Message:    message,
		Errors:     errors,
	}
}

// Error implements the error interface for UnprocessableEntityException.
func (e *UnprocessableEntityException) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Errors)
}
