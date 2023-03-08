package models

import "time"

type MessageText struct {
	Username string `db:"for user"`
	Text     string `mapstructure:"text"`
}

type Message struct {
	Sender         string            //`mapstructure:"from"`
	ConversationID int               `mapstructure:"conversationID"`
	SendDate       time.Time         //`mapstructure:"sendDate"`
	Text           map[string]string `mapstructure:"text"`
}
