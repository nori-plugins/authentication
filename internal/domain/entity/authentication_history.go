package entity

import "time"

type AuthenticationHistory struct {
	ID        uint64
	UserID    uint64
	SigninAt  time.Time
	Meta      string
	SignoutAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}