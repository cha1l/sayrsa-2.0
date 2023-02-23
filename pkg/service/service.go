package service

import "github.com/cha1l/sayrsa-2.0/pkg/repository"

type Authorization interface {
}

type Service struct {
	Authorization
}

func New(repo *repository.Repository) *Service {
	return &Service{}
}
