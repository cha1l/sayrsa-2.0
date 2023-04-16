package models

type Conversation struct {
	Id          int         `json:"conv_id" db:"conv_id"`
	Title       string      `json:"title" db:"title"`
	Members     []PublicKey `json:"members"` //List of conversation members with username and public key
	LastMessage Message     `json:"lastMessage"`
}

func NewConversation(id int, title string, lastMessage Message, members ...PublicKey) *Conversation {
	return &Conversation{
		Id:          id,
		Title:       title,
		Members:     members,
		LastMessage: lastMessage,
	}
}
