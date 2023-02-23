package service

import (
	"github.com/cha1l/sayrsa-2.0/models"
	"github.com/cha1l/sayrsa-2.0/pkg/repository"
)

type Authorization interface {
	CreateUser(u models.User) (string, error)
}

type Service struct {
	Authorization
}

func New(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
	}
}
