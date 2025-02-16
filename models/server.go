package models

import (
	"time"
)

type Server struct {
	ID           string
	SerialNumber string
	IPAddress    string
	Status       string // "pending", "approved"
	CreatedAt    time.Time
}
