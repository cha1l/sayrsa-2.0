package repository

import (
	"fmt"
	"github.com/cha1l/sayrsa-2.0/models"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

type ConversationsRepo struct {
	db *sqlx.DB
}

func NewConversationsRepo(db *sqlx.DB) *ConversationsRepo {
	return &ConversationsRepo{db: db}
}

func (r *ConversationsRepo) GetUsersPublicKeys(usernames ...string) ([]models.PublicKey, error) {
	args := make([]interface{}, 0)
	searchIndexes := make([]string, 0)
	for i, _ := range usernames {
		args = append(args, usernames[i])
		searchIndexes = append(searchIndexes, "$"+strconv.Itoa(i+1))
	}
	indexes := strings.Join(searchIndexes, " OR username=")

	publicKeys := make([]models.PublicKey, 0)

	query := fmt.Sprintf(`SELECT username, public_key FROM %s WHERE username=%s`, usersTable, indexes)

	if err := r.db.Select(&publicKeys, query, args...); err != nil {
		return nil, err
	}

	return publicKeys, nil
}

func (r *ConversationsRepo) CreateConversation(title string, members []string) (int, error) {
	var convID int

	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}

	createConvQuery := fmt.Sprintf(`INSERT INTO %s (title) VALUES ($1) RETURNING id`, conversationsTable)
	row := tx.QueryRow(createConvQuery, title)
	if err := row.Scan(&convID); err != nil {
		if err0 := tx.Rollback(); err0 != nil {
			return 0, err0
		}
		return 0, err
	}

	valuesList := make([]string, 0)
	args := make([]interface{}, 0)
	cnt := 1

	for _, val := range members {
		args = append(args, convID)
		args = append(args, val)

		object := "(" + "$" + strconv.Itoa(cnt) + "," + "$" + strconv.Itoa(cnt+1) + ")"
		valuesList = append(valuesList, object)
		cnt += 2
	}

	values := strings.Join(valuesList, ", ")

	insertingUserQuery := fmt.Sprintf(`INSERT INTO %s (conv_id, user_username) VALUES %s`, conversationMembersTable, values)
	if _, err := tx.Exec(insertingUserQuery, args...); err != nil {
		if err0 := tx.Rollback(); err0 != nil {
			return 0, err0
		}
		return 0, err
	}

	return convID, tx.Commit()
}

func (r *ConversationsRepo) GetUserToken(username string) (models.Token, error) {
	var token models.Token

	query := fmt.Sprintf(`SELECT id, expires_at FROM %s WHERE user_id=(SELECT id FROM %s WHERE username=$1)`, tokensTable, usersTable)
	err := r.db.Get(&token, query, username)

	return token, err
}

func (r *ConversationsRepo) UpdateUserToken(token models.Token) error {
	query := fmt.Sprintf(`UPDATE %s SET expires_at=$1 WHERE id=$2`, tokensTable)
	_, err := r.db.Exec(query, token.ExpiresAt, token.Id)
	return err
}

type SqlResultConv struct {
	Id     int    `db:"id"`
	Title  string `db:"title"`
	Member string `db:"members"`
}
