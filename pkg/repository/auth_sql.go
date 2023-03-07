package repository

import (
	"fmt"
	"time"

	"github.com/cha1l/sayrsa-2.0/models"
	"github.com/jmoiron/sqlx"
)

type AuthRepo struct {
	db *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (s *AuthRepo) CreateUser(u *models.User, token *string, tokenT time.Time) error {
	var id int

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	createUserQuery := fmt.Sprintf(`INSERT INTO %s (username, bio, last_online, password_hash, public_key) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id`, usersTable)
	row := tx.QueryRow(createUserQuery, u.Username, u.Bio, u.LastOnline, u.Password, u.PublicKey)
	if err := row.Scan(&id); err != nil {
		err0 := tx.Rollback()
		if err0 != nil {
			return err0
		}
		return err
	}

	createTokenQuery := fmt.Sprintf(`INSERT INTO %s (user_id, token, expires_at) VALUES ($1, $2, $3)`, tokensTable)
	_, err = tx.Exec(createTokenQuery, id, token, tokenT)
	if err != nil {
		err0 := tx.Rollback()
		if err0 != nil {
			return err0
		}
		return err
	}

	return tx.Commit()
}

func (s *AuthRepo) GetUsersToken(username string, password string) (*models.Token, error) {
	var token models.Token

	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id=(SELECT id FROM %s WHERE username=$1 AND password_hash=$2)`,
		tokensTable, usersTable)

	err := s.db.Get(&token, query, username, password)
	return &token, err
}

func (s *AuthRepo) UpdateUsersToken(token *models.Token) error {
	query := fmt.Sprintf(`UPDATE %s SET token=$1, expires_at=$2 WHERE id=$3`, tokensTable)
	_, err := s.db.Exec(query, token.Token, token.ExpiresAt, token.Id)
	return err
}

func (s *AuthRepo) GetToken(token string) (*models.Token, error) {
	var cToken models.Token

	query := fmt.Sprintf(`SELECT user_id, expires_at FROM %s WHERE token=$1`, tokensTable)
	err := s.db.Get(&cToken, query, token)

	return &cToken, err
}

func (s *AuthRepo) GetUsernameByToken(token *models.Token) (string, error) {
	var username string

	query := fmt.Sprintf(`SELECT username FROM %s WHERE id=$1`, usersTable)
	row := s.db.QueryRowx(query, token.UserId)
	err := row.Scan(&username)

	return username, err
}
