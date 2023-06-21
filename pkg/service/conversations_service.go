package service

import (
	"errors"
	"time"

	"github.com/cha1l/sayrsa-2.0/models"
	"github.com/cha1l/sayrsa-2.0/pkg/repository"
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

func (s *ConversationService) CreateConversation(username string, title string, members []string) (int, []models.PublicKey, error) {
	convID, err := s.repo.CreateConversation(title, members)
	if err != nil {
		return 0, nil, err
	}
	publicKeys, err := s.repo.GetUsersPublicKeys(members...)
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

	token.ExpiresAt = token.ExpiresAt.Add(UpdateTokenTime)

	return s.repo.UpdateUserToken(token)
}

func (s *ConversationService) GetPublicKey(username string) (string, error) {
	publicKeys, err := s.repo.GetUsersPublicKeys(username)
	if err != nil {
		return "", err
	}
	if len(publicKeys) == 1 {
		publicKey := publicKeys[0]
		return publicKey.PublicKey, s.UpdateToken(username)
	}
	return "", errors.New("wrong length of slice")
}

func (s *ConversationService) GetAllConversations(username string) ([]*models.Conversation, error) {
	if err := s.UpdateToken(username); err != nil {
		return nil, err
	}
	return s.repo.GetAllConversations(username)
}
