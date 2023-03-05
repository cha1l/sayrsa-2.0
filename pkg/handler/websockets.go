package handler

import (
	models "github.com/cha1l/sayrsa-2.0/models"
	"github.com/mitchellh/mapstructure"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

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
	defer delete(h.clients, username)
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			WsErrorResponse(conn, err.Error())
			return
		}
	}(conn)

	log.Printf("Client %s connected", username)

	for {
		var input models.StandardInput
		err := conn.ReadJSON(&input)
		if err != nil {
			log.Printf("client %s disconnect", username)
			break
		}

		if input.Action == createConversationAction {
			var conv models.CreateConversionsInput

			//algorithm???  data -> conv
			err := mapstructure.Decode(input.Data, &conv)
			if err != nil {
				WsErrorResponse(conn, err.Error())
				continue
			}

			log.Printf("User %s wants to create a conversation with users %s", username, conv.Usernames)

			convID, publicKeys, err := h.service.Conversations.CreateConversation(username, conv)
			if err != nil {
				WsErrorResponse(conn, err.Error())
				continue
			}

			data := map[string]interface{}{
				"event":       "new_conv",
				"conv_id":     convID,
				"public_keys": publicKeys,
			}

			go h.SendMessage(data, conv.Usernames...)

		} else if input.Action == sendMessageAction {
			var msg models.SendMessageInput

			err := mapstructure.Decode(input.Data, &msg)
			if err != nil {
				WsErrorResponse(conn, err.Error())
				continue
			}

			log.Printf("User %s wants to send message", username)

		} else {
			WsErrorResponse(conn, "invalid action")
			continue
		}

	}
}

func (h *Handler) SendMessage(data map[string]interface{}, users ...string) {
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
