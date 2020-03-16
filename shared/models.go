package shared

import "net/http"

type PublicError interface {
	publicError()
}

// ShowableError is an error to show to the user
type ShowableError struct {
	Message    string
	StatusCode int
}

func (e ShowableError) Error() string {
	return e.Message
}

func (e ShowableError) publicError() {}

// NewShowableError returns a new error to show to the user
func NewShowableError(message string, statusCode int) ShowableError {
	return ShowableError{message, statusCode}
}

// NewBadRequestError returns a bad request error
func NewBadRequestError(message string) ShowableError {
	return NewShowableError(message, http.StatusBadRequest)
}
