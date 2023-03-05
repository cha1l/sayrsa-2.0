package models

import "time"

type Token struct {
	Id        int       `db:"id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	UserId    int       `db:"user_id"`
}
