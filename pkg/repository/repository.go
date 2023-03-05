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
	GetToken(token string) (models.Token, error)
	GetUsernameByToken(token models.Token) (string, error)
}

type Conversations interface {
	CreateConversation(input models.CreateConversionsInput) (int, error)
	GetUsersPublicKeys(usernames []string) ([]models.PublicKey, error)
	GetUserToken(username string) (models.Token, error)
	UpdateUserToken(token models.Token) error
}

type Repository struct {
	Authorization
	Conversations
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepo(db),
		Conversations: NewConversationsRepo(db),
	}
}
