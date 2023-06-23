package service

import (
	"errors"
	"log"
	"time"

	"github.com/cha1l/sayrsa-2.0/models"
	"github.com/cha1l/sayrsa-2.0/pkg/repository"
)

const (
	salt        string = "ojdflkghjfdlkj"
	tokenLength int    = 64
	tokenT             = 7 * 24 * time.Hour
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) CreateUser(u models.User) (string, error) {
	if u.PrivateKey == "" {
		return "", errors.New("empty private key")
	}

	u.Password = GeneratePasswordHash(u.Password)
	u.PrivateKey = EncryptPrivateKey(u.PrivateKey, u.Password)
	token := GenerateSecureToken(tokenLength)
	tokenTime := time.Now().Add(tokenT)

	if err := s.repo.CreateUser(u, token, tokenTime); err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) GetUserTokenPrivateKey(username, password string) (string, string, error) {
	password = GeneratePasswordHash(password)
	token, encoded, err := s.repo.GetUserTokenPrivateKey(username, password)
	if err != nil {
		return "", "", err
	}

	privateKey := DecryptPrivateKey(encoded, password)

	if time.Now().After(token.ExpiresAt) {
		log.Println("token is not valid: creating new token ...")
		token.Token = GenerateSecureToken(tokenLength)
		token.ExpiresAt = time.Now().Add(tokenT)
		err = s.repo.UpdateUsersToken(token)
	}
	return token.Token, privateKey, err
}

func (s *AuthService) GetUsernameByToken(token string) (string, error) {
	cToken, err := s.repo.GetToken(token)
	if err != nil {
		return "", err
	}
	if cToken.ExpiresAt.Before(time.Now()) {
		return "", errors.New("invalid token, register one more time")
	}

	return cToken.UserUsername, nil

}
