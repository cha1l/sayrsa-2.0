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

func (s *AuthRepo) CreateUser(u models.User, token string, tokenT time.Time) (int, *sqlx.Tx, error) {
	var id int

	tx, err := s.db.Beginx()
	if err != nil {
		return 0, nil, err
	}

	createUserQuery := fmt.Sprintf(`INSERT INTO %s (username, bio, last_online, password_hash, public_key) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id`, usersTable)
	row := tx.QueryRowx(createUserQuery, u.Username, u.Bio, u.LastOnline, u.Password, u.PublicKey)
	if err := row.Scan(&id); err != nil {
		if err0 := tx.Rollback(); err0 != nil {
			return 0, nil, err0
		}
		return 0, nil, err
	}

	createTokenQuery := fmt.Sprintf(`INSERT INTO %s (user_username, token, expires_at) VALUES ($1, $2, $3)`, tokensTable)
	_, err = tx.Exec(createTokenQuery, u.Username, token, tokenT)
	if err != nil {
		err0 := tx.Rollback()
		if err0 != nil {
			return 0, nil, err0
		}
		return 0, nil, err
	}

	return id, tx, nil
}

type GetUserTokenPrivateKeyResut struct {
	Id           int       `db:"id"`
	Token        string    `db:"token"`
	ExpiresAt    time.Time `db:"expires_at"`
	UserUsername string    `db:"user_username"`
	PrivateKey   string    `db:"private_key"`
}

func (s *AuthRepo) GetUserTokenPrivateKey(username string, password string) (models.Token, string, error) {
	var res GetUserTokenPrivateKeyResut

	query := fmt.Sprintf(`SELECT t.*, u.private_key FROM %s AS t LEFT JOIN %s AS u 
	ON u.username=t.user_username AND u.username=$1 AND u.password_hash=$2`,
		tokensTable, usersTable)

	err := s.db.Get(&res, query, username, password)

	token := models.NewToken(res.Id, res.Token, res.UserUsername, res.ExpiresAt)

	return token, res.PrivateKey, err
}

func (s *AuthRepo) UpdateUsersToken(token models.Token) error {
	query := fmt.Sprintf(`UPDATE %s SET token=$1, expires_at=$2 WHERE id=$3`, tokensTable)
	_, err := s.db.Exec(query, token.Token, token.ExpiresAt, token.Id)
	return err
}

func (s *AuthRepo) GetToken(token string) (models.Token, error) {
	var cToken models.Token

	query := fmt.Sprintf(`SELECT user_username, expires_at FROM %s WHERE token=$1`, tokensTable)
	err := s.db.Get(&cToken, query, token)

	return cToken, err
}

func (s *AuthRepo) SetUserPrivateKey(userID int, key string, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`UPDATE %s SET private_key = $1`, usersTable)

	_, err := tx.Exec(query, key)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}
