package shared

import (
	"net/http"
	"time"
)

// PatchOp represents a patch operation
type PatchOp string

const (
	// PatchOpReplace represents replace operation on a resource property
	PatchOpReplace PatchOp = "replace"
)

type TokenClaims struct {
	ClientType ClientType `json:"type"`
	ID         int        `json:"id"`
	IAT        time.Time  `json:"iat"`
}

// Client is the type of requester of a resource
type Client interface{ client() }

// ClientType identifies type of client
type ClientType string

// ClientTypeUser identifies regular user
var ClientTypeUser ClientType = "user"

// ClientTypeMechanic identifies mechanic
var ClientTypeMechanic ClientType = "mechanic"

// ClientTypeAdmin identifies admin
var ClientTypeAdmin ClientType = "admin"

// PatchRequestBody is the representation of the body of a PATCH request
type PatchRequestBody []struct {
	Op    PatchOp `json:"op"`
	Path  string  `json:"path"`
	Value string  `json:"value,omitempty"`
}

// PublicError is an error to show to the user
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
