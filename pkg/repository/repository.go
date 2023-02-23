package repository

import "github.com/jmoiron/sqlx"

type Authorization interface {
}

type Repository struct {
	Authorization
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		//interface
	}
}
