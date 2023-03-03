package service

import (
	"github.com/cha1l/sayrsa-2.0/models"
	"github.com/cha1l/sayrsa-2.0/pkg/repository"
)

type Authorization interface {
	CreateUser(u models.User) (string, error)
	GetUsersToken(u models.SignInInput) (string, error)
	GetUserIdByToken(token string) (int, error)
}

type Conversations interface {
}

type Service struct {
	Authorization
	Conversations
}

func New(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		Conversations: NewConversationService(repo.Conversations),
	}
}
