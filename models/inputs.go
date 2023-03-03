package models

type SignInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateConversionsInput struct {
	Title   string `json:"title"`
	UsersID []int  `json:"users_id"`
}

type SendMessageInput struct {
	Text string `json:"text"`
}

type StandardInput struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}
