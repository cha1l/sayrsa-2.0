package repository

import (
	"fmt"
	"github.com/cha1l/sayrsa-2.0/models"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

type MessagesRepo struct {
	db *sqlx.DB
}

func NewMessagesRepo(db *sqlx.DB) *MessagesRepo {
	return &MessagesRepo{db: db}
}

func (m *MessagesRepo) GetMessages(convID int, offset int, amount int) ([]models.Message, error) {
	//todo
	return nil, nil
}

func (m *MessagesRepo) SendMessage(msg *models.Message) error {
	var messageGlobalId int
	var messageConvId int

	tx, err := m.db.Beginx()
	if err != nil {
		return err
	}

	createMessageQuery := fmt.Sprintf(`INSERT INTO %s (sender_username, conv_id, send_date)
		VALUES ($1, $2, $3) RETURNING id`, messagesTable)
	row := tx.QueryRowx(createMessageQuery, msg.Sender, msg.ConversationID, msg.SendDate)
	if err := row.Scan(&messageGlobalId); err != nil {
		if err0 := tx.Rollback(); err0 != nil {
			return err0
		}
		return err
	}

	conversationQuery := fmt.Sprintf(`INSERT INTO %s (conv_id, message_id) 
		VALUES ($1, $2) RETURNING id`, conversationMessagesTable)
	row = tx.QueryRowx(conversationQuery, msg.ConversationID, messageGlobalId)
	if err := row.Scan(&messageConvId); err != nil {
		if err0 := tx.Rollback(); err0 != nil {
			return err0
		}
		return err
	}

	args := make([]interface{}, 0)
	valueList := make([]string, 0)
	cnt := 1

	for key, element := range msg.Text {
		args = append(args, messageConvId)
		args = append(args, element)
		args = append(args, key)

		object := fmt.Sprintf(`($%s, $%s, $%s)`, strconv.Itoa(cnt), strconv.Itoa(cnt+1), strconv.Itoa(cnt+2))
		valueList = append(valueList, object)
		cnt += 3
	}

	values := strings.Join(valueList, ", ")

	textQuery := fmt.Sprintf(`INSERT INTO %s (conv_message_id, text, for_user) VALUES %s`, messageTextTable, values)
	_, err = tx.Exec(textQuery, args...)
	if err != nil {
		if err0 := tx.Rollback(); err0 != nil {
			return err0
		}
		return err
	}

	return tx.Commit()
}
