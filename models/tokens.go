package models

import "time"

type Token struct {
	Id         int       `db:"id"`
	Token      string    `db:"token"`
	Expires_at time.Time `db:"expires_at"`
	User_id    int
}
