package models

import "time"

type MessageText struct {
	Username string `db:"for user"`
	Text     string `mapstructure:"text"`
}

type SendMessage struct {
	Sender         string
	ConversationID int `mapstructure:"conversationID"`
	SendDate       time.Time
	Text           map[string]string `mapstructure:"text"`
}

type GetMessage struct {
	Sender         string    `json:"sender"`
	ConversationID int       `json:"conversationID"`
	SendDate       time.Time `json:"sendDate"`
	Text           string    `json:"text"`
}
