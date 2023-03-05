package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
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

func (s *AuthService) CreateUser(u models.User) (string, error) {
	u.Password = GeneratePasswordHash(u.Password)
	token := GenerateSecureToken(tokenLength)
	tokenT := time.Now().Add(tokenT)

	return token, s.repo.CreateUser(u, token, tokenT)
}

func (s *AuthService) GetUsersToken(u models.SignInInput) (string, error) {
	u.Password = GeneratePasswordHash(u.Password)
	token, err := s.repo.GetUsersToken(u)
	if err != nil {
		return "", err
	}
	if time.Now().After(token.Expires_at) {
		log.Println("token is not valid: creating new token ...")
		token.Token = GenerateSecureToken(tokenLength)
		token.Expires_at = time.Now().Add(tokenT)
		err = s.repo.UpdateUsersToken(token)
	}
	return token.Token, err
}

func (s *AuthService) GetUsernameByToken(token string) (string, error) {
	return s.repo.GetUsernameByToken(token)
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
