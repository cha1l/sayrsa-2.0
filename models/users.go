package models

import "time"

type User struct {
	Id         int    `json:"-"`
	Username   string `json:"username"`
	Bio        string `json:"bio"`
	Password   string `json:"password"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
	LastOnline time.Time
}

type PublicKey struct {
	Username  string `json:"username" db:"username"`
	PublicKey string `json:"publicKey" db:"public_key"`
}

func NewPublicKey(username string, publicKey string) PublicKey {
	return PublicKey{Username: username, PublicKey: publicKey}
}
