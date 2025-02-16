package models

import (
	"time"
)

type Admin struct {
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	IsRoot    bool      `db:"is_root"`
	CreatedAt time.Time `db:"created_at"`
}
