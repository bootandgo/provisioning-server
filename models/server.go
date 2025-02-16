package models

import (
	"time"
)

type Server struct {
	ID           string    `db:"id"`
	SerialNumber string    `db:"serial_number"`
	IPAddress    string    `db:"ip_address"`
	Status       string    `db:"status"` // "pending", "approved"
	CreatedAt    time.Time `db:"created_at"`
}
