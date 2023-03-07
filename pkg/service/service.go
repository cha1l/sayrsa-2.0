package service

import (
	"github.com/cha1l/sayrsa-2.0/models"
	"github.com/cha1l/sayrsa-2.0/pkg/repository"
)

type Authorization interface {
	CreateUser(u *models.User) (*string, error)
	GetUsersToken(username string, password string) (*string, error)
	GetUsernameByToken(token string) (string, error)
}

type Conversations interface {
	CreateConversation(username string, title string, members []string) (int, []models.PublicKey, error)
	UpdateToken(username string) error
	GetPublicKey(username string) (string, error)
	GetConversationInfo(convID int) (models.Conversation, error)
}

type Service struct {
	Authorization
	Conversations
}

func New(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		Conversations: NewConversationService(repo),
	}
}
