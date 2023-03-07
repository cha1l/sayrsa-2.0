package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/cha1l/sayrsa-2.0/models"
	"github.com/cha1l/sayrsa-2.0/pkg/repository"
)

const (
	salt        string        = "ojdflkghjfdlkj"
	tokenLength int           = 64
	tokenT      time.Duration = 7 * 24 * time.Hour
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) CreateUser(u *models.User) (*string, error) {
	u.Password = GeneratePasswordHash(u.Password)
	token := GenerateSecureToken(tokenLength)
	tokenT := time.Now().Add(tokenT)

	return &token, s.repo.CreateUser(u, &token, tokenT)
}

func (s *AuthService) GetUsersToken(username string, password string) (*string, error) {
	password = GeneratePasswordHash(password)
	token, err := s.repo.GetUsersToken(username, password)
	if err != nil {
		return nil, err
	}
	if time.Now().After(token.ExpiresAt) {
		log.Println("token is not valid: creating new token ...")
		token.Token = GenerateSecureToken(tokenLength)
		token.ExpiresAt = time.Now().Add(tokenT)
		err = s.repo.UpdateUsersToken(token)
	}
	return &token.Token, err
}

func (s *AuthService) GetUsernameByToken(token string) (string, error) {
	cToken, err := s.repo.GetToken(token)
	if err != nil {
		return "", err
	}
	if cToken.ExpiresAt.Before(time.Now()) {
		return "", errors.New("invalid token, register one more time")
	}

	return s.repo.GetUsernameByToken(cToken)

}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func GeneratePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
