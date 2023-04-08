package service

import (
	"github.com/cha1l/sayrsa-2.0/models"
	"github.com/cha1l/sayrsa-2.0/pkg/repository"
)

type Authorization interface {
	CreateUser(u models.User) (string, error)
	GetUsersToken(username string, password string) (string, error)
	GetUsernameByToken(token string) (string, error)
}

type Conversations interface {
	CreateConversation(username string, title string, members []string) (int, []models.PublicKey, error)
	UpdateToken(username string) error
	GetPublicKey(username string) (string, error)
	GetConversationInfo(username string, convID int) (*models.Conversation, error)
	GetAllConversations(username string) ([]*models.Conversation, error)
}

type Messages interface {
	SendMessage(username string, message *models.SendMessage) error
	GetMessages(username string, convID int, offset int, amount int) ([]models.GetMessage, error)
}

type Service struct {
	Authorization
	Conversations
	Messages
}

func New(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		Conversations: NewConversationService(repo.Conversations),
		Messages:      NewMessagesService(repo.Messages),
	}
}
