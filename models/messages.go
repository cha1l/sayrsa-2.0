package models

import "time"

type MessageText struct {
	Username string `db:"for user"`
	Text     string `mapstructure:"text"`
}

type Message struct {
	Sender         string
	ConversationID int `mapstructure:"conversationID"`
	SendDate       time.Time
	Text           map[string]string `mapstructure:"text"`
}
