package auth

import (
	"time"

	"github.com/CartechAPI/shared"
)

// Session represents an users session
type Session struct {
	CreatedAt time.Time         `json:"created_at"`
	SessionID int               `json:"session_id"`
	UserID    int               `json:"user_id"`
	UserType  shared.ClientType `json:"user_type"`
	Token     string            `json:"token"`
}
