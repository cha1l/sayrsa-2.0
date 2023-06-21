package service

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
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
	token := GenerateSecureToken(tokenLength)
	tokenT := time.Now().Add(tokenT)

	id, tx, err := s.repo.CreateUser(u, token, tokenT)
	fmt.Println("id is ", id)
	if err != nil {
		return "", err
	}

	privateKey := EncodePrivateKey(u.PrivateKey, id)

	if err := s.repo.SetUserPrivateKey(id, privateKey, tx); err != nil {
		return "", err
	}

	return token, err
}

func (s *AuthService) GetUserTokenPrivateKey(username, password string) (string, string, error) {
	password = GeneratePasswordHash(password)
	token, decoded, err := s.repo.GetUserTokenPrivateKey(username, password)
	if err != nil {
		return "", "", err
	}

	privateKey, err := DecodePrivateKey(decoded)
	if err != nil {
		return "", "", err
	}

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

func EncodePrivateKey(privateKey string, userID int) string {
	userIDstr := strconv.Itoa(userID)

	toEncode := []byte(strings.Join([]string{privateKey, userIDstr, salt}, " "))

	return base64.StdEncoding.EncodeToString(toEncode)
}

func DecodePrivateKey(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		return "", err
	}

	privateKey := bytes.Split(data, []byte(" "))[0]

	return string(privateKey), nil
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
