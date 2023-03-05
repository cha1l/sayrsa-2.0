package models

type SignInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateConversionsInput struct {
	Title     string   `mapstructure:"title"`
	Usernames []string `mapstructure:"usernames"`
}

type SendMessageInput struct {
	Text string `mapstructure:"text"`
}

type StandardInput struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}
