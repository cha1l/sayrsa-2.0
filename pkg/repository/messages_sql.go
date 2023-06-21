package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cha1l/sayrsa-2.0/models"
	"github.com/jmoiron/sqlx"
)

type MessagesRepo struct {
	db *sqlx.DB
}

func NewMessagesRepo(db *sqlx.DB) *MessagesRepo {
	return &MessagesRepo{db: db}
}

type GetMessagesResult struct {
	IdInConv       int       `db:"id_in_conv"`
	ConversationID int       `db:"conv_id"`
	SenderUsername string    `db:"sender_username"`
	SendDate       time.Time `db:"send_date"`
	Text           string    `db:"text"`
}

func (m *MessagesRepo) GetMessages(username string, convID int, offset int, amount int) (*[]models.Message, error) {
	messages := make([]models.Message, 0)

	query := fmt.Sprintf(`SELECT m.id_in_conv, m.sender_username, m.send_date, t.text FROM %s AS m INNER JOIN %s AS t
	     ON t.id=m.id WHERE m.id_in_conv <= $1 AND m.id_in_conv > $2 AND m.conv_id=$3 AND t.for_user=$4`, messagesTable, messageTextTable)
	rows, err := m.db.Queryx(query, offset, offset-amount, convID, username)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var r GetMessagesResult
		if err := rows.StructScan(&r); err != nil {
			return nil, err
		}
		message := models.NewMessage(r.IdInConv, r.SenderUsername, convID, r.SendDate, r.Text)
		messages = append(messages, message)
	}

	return &messages, nil
}

func (m *MessagesRepo) SendMessage(msg *models.SendMessage) error {
	var messageGlobalId int
	var messageConvId int

	tx, err := m.db.Beginx()
	if err != nil {
		return err
	}

	idQuery := fmt.Sprintf(`SELECT COALESCE(MAX(id_in_conv), 0) FROM %s WHERE conv_id=$1 `, messagesTable)
	row := tx.QueryRowx(idQuery, msg.ConversationID)
	if err = row.Scan(&messageConvId); err != nil {
		if err0 := tx.Rollback(); err0 != nil {
			return err0
		}
		return err
	}

	createMessageQuery := fmt.Sprintf(`INSERT INTO %s (id_in_conv, sender_username, conv_id, send_date)
		VALUES ($1, $2, $3, $4) RETURNING id`, messagesTable)
	row = tx.QueryRowx(createMessageQuery, messageConvId+1, msg.Sender, msg.ConversationID, msg.SendDate)
	if err := row.Scan(&messageGlobalId); err != nil {
		if err0 := tx.Rollback(); err0 != nil {
			return err0
		}
		return err
	}

	args := make([]interface{}, 0)
	valueList := make([]string, 0)
	cnt := 1

	for key, element := range msg.Text {
		args = append(args, messageGlobalId)
		args = append(args, element)
		args = append(args, key)

		object := fmt.Sprintf(`($%s, $%s, $%s)`, strconv.Itoa(cnt), strconv.Itoa(cnt+1), strconv.Itoa(cnt+2))
		valueList = append(valueList, object)
		cnt += 3
	}

	values := strings.Join(valueList, ", ")

	textQuery := fmt.Sprintf(`INSERT INTO %s (id, text, for_user) VALUES %s`, messageTextTable, values)
	_, err = tx.Exec(textQuery, args...)
	if err != nil {
		if err0 := tx.Rollback(); err0 != nil {
			return err0
		}
		return err
	}

	return tx.Commit()
}
