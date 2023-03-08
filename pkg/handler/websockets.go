package handler

import (
	"github.com/cha1l/sayrsa-2.0/models"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"log"
	"net/http"
)

type CreateConversionsInput struct {
	Title     string   `mapstructure:"title"`
	Usernames []string `mapstructure:"members"`
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		WsErrorResponse(conn, err.Error())
		return
	}

	username := GetParams(r.Context())
	h.clients[username] = NewClient(conn)

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			WsErrorResponse(conn, err.Error())
			return
		}
	}(conn)
	defer delete(h.clients, username)

	type StandardInput struct {
		Action string                 `json:"action"`
		Data   map[string]interface{} `json:"data"`
	}

	log.Printf("Client %s connected", username)

	for {
		var input StandardInput
		err := conn.ReadJSON(&input)
		if err != nil {
			log.Printf("client %s disconnect", username)
			break
		}

		if input.Action == createConversationAction {
			var conv CreateConversionsInput

			//algorithm???  data -> conv
			err := mapstructure.Decode(input.Data, &conv)
			if err != nil {
				WsErrorResponse(conn, err.Error())
				continue
			}

			log.Printf("User %s wants to create a conversation with users %s", username, conv.Usernames)

			convID, publicKeys, err := h.service.Conversations.CreateConversation(username, conv.Title, conv.Usernames)
			if err != nil {
				WsErrorResponse(conn, err.Error())
				continue
			}

			data := map[string]interface{}{
				"event": "new_conv",
				"data": map[string]interface{}{
					"conv_id":     convID,
					"public_keys": publicKeys,
				},
			}

			go h.SendInfo(data, conv.Usernames...)

		} else if input.Action == sendMessageAction {
			var message models.Message

			if err := mapstructure.Decode(input.Data, &message); err != nil {
				WsErrorResponse(conn, err.Error())
				continue
			}

			if err := h.service.Messages.SendMessage(username, &message); err != nil {
				WsErrorResponse(conn, err.Error())
				continue
			}

			for key, value := range message.Text {
				go func(message models.Message, text string, user string) {
					log.Printf("User %s send message to user %s", username, user)
					data := GenerateMessage(text, message)
					h.SendInfo(data, user)
				}(message, value, key)

			}

		} else {
			WsErrorResponse(conn, "invalid action")
			continue
		}
	}

}

func (h *Handler) SendInfo(data map[string]interface{}, users ...string) {
	if len(users) != 0 {
		for _, username := range users {
			go func(username string) {
				if val, ok := h.clients[username]; ok {
					if err := val.connection.WriteJSON(data); err != nil {
						WsErrorResponse(val.connection, err.Error())
						return
					}
					return
				}
				log.Printf("client %s is not connected", username)
			}(username)
		}
		return
	}
	log.Println("empty users list")
}

func GenerateMessage(text string, message models.Message) map[string]interface{} {
	data := map[string]interface{}{
		"event": "new_message",
		"data": map[string]interface{}{
			"conv_id":   message.ConversationID,
			"from":      message.Sender,
			"send_date": message.SendDate,
			"text":      text,
		},
	}

	return data
}
