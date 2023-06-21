package models

import "time"

type Token struct {
	Id           int       `db:"id"`
	Token        string    `db:"token"`
	ExpiresAt    time.Time `db:"expires_at"`
	UserUsername string    `db:"user_username"`
}

func NewToken(id int, token, userUsername string, expiresAt time.Time) Token {
	return Token{
		Id:           id,
		Token:        token,
		ExpiresAt:    expiresAt,
		UserUsername: userUsername,
	}
}
