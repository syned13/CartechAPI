package order

import "time"

// ServiceOrderStatus is the status of a service order
type ServiceOrderStatus string

const (
	// ServiceOrderStatusPending status pending
	ServiceOrderStatusPending ServiceOrderStatus = "pending"
	// ServiceOrderStatusInProgress status in progress
	ServiceOrderStatusInProgress ServiceOrderStatus = "in_progress"
	// ServiceOrderStatusCancelled status cancelled
	ServiceOrderStatusCancelled ServiceOrderStatus = "cancelled"
	// ServiceOrderStatusFinished status finished
	ServiceOrderStatusFinished ServiceOrderStatus = "finished"
	// ServiceOrderStatusFailure status failure
	ServiceOrderStatusFailure ServiceOrderStatus = "failure"
)

// ServiceOrder represents a service order
type ServiceOrder struct {
	ServiceOrderID int                `json:"service_order_id"`
	ServiceID      int                `json:"service_id"`
	UserID         int                `json:"user_id"`
	MechanicID     int                `json:"mechanic_id"`
	CreatedAt      *time.Time         `json:"created_at"`
	StartedAt      *time.Time         `json:"started_at"`
	Status         ServiceOrderStatus `json:"status"`
	FinishedAt     *time.Time         `json:"finished_at"`
	CancelledAt    *time.Time         `json:"cancelled_at"`
}
