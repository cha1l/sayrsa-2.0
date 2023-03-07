package models

type Conversation struct {
	Id    int    `json:"conv_id"`
	Title string `json:"title"`
	//List of conversation members with username and public key
	Members []PublicKey `json:"members"`
}

func NewConversation(id int, title string, members []PublicKey) *Conversation {
	return &Conversation{Id: id, Title: title, Members: members}
}
