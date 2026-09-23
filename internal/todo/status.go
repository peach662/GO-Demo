package todo

import "errors"

type Status string

const (
	StatusPending    Status = "PENDING"
	StatusProcessing Status = "PROCESSING"
	StatusCompleted  Status = "COMPLETED"
	StatusFailed     Status = "FAILED"
)

var ErrInvalidTransition = errors.New("invalid status transition")

func CanTransition(from, to Status) bool {
	if from == to {
		return true
	}

	switch from {
	case StatusPending:
		return to == StatusProcessing
	case StatusProcessing:
		return to == StatusCompleted || to == StatusFailed
	case StatusFailed:
		return to == StatusProcessing
	default:
		return false
	}
}
