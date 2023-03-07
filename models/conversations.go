package models

type Conversation struct {
	Id    int    `json:"conv_id"`
	Title string `json:"title"`
	//List of conversation members with username and public key
	Members []string `json:"members"`
}
