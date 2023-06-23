package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

const (
	usersTable               = "users"
	tokensTable              = "tokens"
	conversationsTable       = "conversations"
	conversationMembersTable = "conversation_members"
	messagesTable            = "messages"
	messageTextTable         = "message_text"
)

type Config struct {
	Host     string
	User     string
	Password string
	DBname   string
}

func NewDB(c Config) (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s",
		c.Host, c.User, c.Password, c.DBname)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		log.Println("Db pinging failed")
		return nil, err
	}

	log.Println("Database created")

	return db, nil
}
