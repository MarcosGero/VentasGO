package sale

import "time"

type Status string

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
)

var validStatus = map[Status]struct{}{
	StatusPending:  {},
	StatusApproved: {},
	StatusRejected: {},
}

func IsValidStatus(s Status) bool {
	_, ok := validStatus[s]
	return ok
}

// Sale representa una transacción de venta.
type Sale struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Amount    float64   `json:"amount"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}
