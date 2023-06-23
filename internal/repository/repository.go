package repository

import (
	"time"

	"github.com/cha1l/sayrsa-2.0/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(u models.User, token string, tokenT time.Time) error
	GetUserTokenPrivateKey(username string) (models.Token, string, string, error)
	UpdateUsersToken(token models.Token) error
	GetToken(token string) (models.Token, error)
}

type Conversations interface {
	CreateConversation(title string, members []string) (int, error)
	GetUsersPublicKeys(usernames ...string) ([]models.PublicKey, error)
	GetUserToken(username string) (models.Token, error)
	UpdateUserToken(token models.Token) error
	GetAllConversations(username string) ([]*models.Conversation, error)
}

type Messages interface {
	GetMessages(username string, convID int, offset int, amount int) (*[]models.Message, error)
	SendMessage(message *models.SendMessage) error
}

type Repository struct {
	Authorization
	Conversations
	Messages
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepo(db),
		Conversations: NewConversationsRepo(db),
		Messages:      NewMessagesRepo(db),
	}
}
