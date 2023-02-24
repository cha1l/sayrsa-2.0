package repository

import (
	"time"

	"github.com/cha1l/sayrsa-2.0/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(u models.User, token string, tokenT time.Time) error
	GetUsersToken(u models.SignInInput) (models.Token, error)
	UpdateUsersToken(token models.Token) error
}

type Repository struct {
	Authorization
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepo(db),
	}
}
