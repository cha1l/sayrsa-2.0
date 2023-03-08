package service

import (
	"github.com/cha1l/sayrsa-2.0/models"
	"github.com/cha1l/sayrsa-2.0/pkg/repository"
	"time"
)

type MessagesService struct {
	repo repository.Messages
}

func NewMessagesService(repo repository.Messages) *MessagesService {
	return &MessagesService{repo: repo}
}

func (m *MessagesService) SendMessage(username string, message *models.Message) error {
	message.Sender = username
	message.SendDate = time.Now()
	return m.repo.SendMessage(message)
}
