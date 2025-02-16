package models

import (
	"time"
)

type Admin struct {
	Username  string
	Password  string // hashed
	CreatedAt time.Time
}
