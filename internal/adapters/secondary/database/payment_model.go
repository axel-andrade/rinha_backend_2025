package database

import "time"

type Payment struct {
	ID          string    `json:"id"`
	Amount      float64   `json:"amount"`
	Processor   string    `json:"processor"` // "default" or "fallback"
	RequestedAt time.Time `json:"requested_at"`
}
