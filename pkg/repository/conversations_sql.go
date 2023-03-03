package repository

import "github.com/jmoiron/sqlx"

type ConversationsRepo struct {
	db *sqlx.DB
}

func NewConversationsRepo(db *sqlx.DB) *ConversationsRepo {
	return &ConversationsRepo{db: db}
}
