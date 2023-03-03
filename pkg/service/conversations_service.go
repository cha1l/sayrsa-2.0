package service

import "github.com/cha1l/sayrsa-2.0/pkg/repository"

type ConversationService struct {
	repo repository.Conversations
}

func NewConversationService(repo repository.Conversations) *ConversationService {
	return &ConversationService{
		repo: repo,
	}
}

func (s *ConversationService) CreateConversations(title string, users []int) (int, error) {
	return 0, nil
}
