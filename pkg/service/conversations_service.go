package service

import (
	"github.com/cha1l/sayrsa-2.0/models"
	"github.com/cha1l/sayrsa-2.0/pkg/repository"
	"time"
)

const (
	UpdateTokenTime = 24 * time.Hour
)

type ConversationService struct {
	repo repository.Conversations
}

func NewConversationService(repo repository.Conversations) *ConversationService {
	return &ConversationService{
		repo: repo,
	}
}

func (s *ConversationService) CreateConversation(username string, input models.CreateConversionsInput) (int, []models.PublicKey, error) {
	input.Usernames = append(input.Usernames, username)

	convID, err := s.repo.CreateConversation(input)
	if err != nil {
		return 0, nil, err
	}
	publicKeys, err := s.repo.GetUsersPublicKeys(input.Usernames)
	if err != nil {
		return 0, nil, err
	}

	return convID, publicKeys, s.UpdateToken(username)
}

func (s *ConversationService) UpdateToken(username string) error {
	token, err := s.repo.GetUserToken(username)
	if err != nil {
		return err
	}

	token.Expires_at = token.Expires_at.Add(UpdateTokenTime)

	return s.repo.UpdateUserToken(token)
}
