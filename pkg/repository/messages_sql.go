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
	//messages := make([]models.Message, amount)

	return nil, nil
}

func (m *MessagesRepo) SendMessage(msg *models.Message) error {
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
