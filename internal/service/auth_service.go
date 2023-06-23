package service

import (
	"errors"
	"log"
	"time"

	"github.com/cha1l/sayrsa-2.0/internal/repository"
	"github.com/cha1l/sayrsa-2.0/models"
	"golang.org/x/crypto/bcrypt"
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

	hashedPassword, err := GeneratePasswordHash(u.Password)
	if err != nil {
		return "", err
	}

	u.Password = hashedPassword
	u.PrivateKey = EncryptPrivateKey(u.PrivateKey, u.Password)
	token := GenerateSecureToken(tokenLength)
	tokenTime := time.Now().Add(tokenT)

	if err := s.repo.CreateUser(u, token, tokenTime); err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) GetUserTokenPrivateKey(username, password string) (string, string, error) {
	token, encoded, password_hash, err := s.repo.GetUserTokenPrivateKey(username)
	if err != nil {
		return "", "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password)); err != nil {
		return "", "", err
	}

	privateKey := DecryptPrivateKey(encoded, password_hash)

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
