package models

import (
	"time"
)

type Invitation struct {
	Token     string    `db:"token"`
	CreatedBy string    `db:"created_by"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
	Used      bool      `db:"used"`
}
