package handler

import (
	"encoding/json"
	models "github.com/cha1l/sayrsa-2.0/models"
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

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			WsErrorResponse(conn, err.Error())
			return
		}
	}(conn)

	userId := GetParams(r.Context())
	h.clients[conn] = NewClient(userId, true)
	defer delete(h.clients, conn)

	for {
		var input models.StandardInput
		err := conn.ReadJSON(&input)
		if err != nil {
			WsErrorResponse(conn, err.Error())
			break
		}

		if input.Action == createConversationAction {
			var conv models.CreateConversionsInput

			//algorithm???  data -> conv
			js, err := json.Marshal(input.Data)
			if err != nil {
				WsErrorResponse(conn, err.Error())
				break
			}
			if err != json.Unmarshal(js, &conv) {
				WsErrorResponse(conn, err.Error())
				break
			}

			log.Printf("User %d wants to create a conversation with users %d", userId, conv.UsersID)

			//todo: call service for creating conversions

		} else if input.Action == sendMessageAction {
			var msg models.SendMessageInput

			js, err := json.Marshal(input.Data)
			if err != nil {
				WsErrorResponse(conn, err.Error())
				break
			}
			if err != json.Unmarshal(js, &msg) {
				WsErrorResponse(conn, err.Error())
				break
			}

			log.Printf("User %d wants to send message", userId)

		} else {
			WsErrorResponse(conn, "invalid action")
			continue
		}
	}
}
