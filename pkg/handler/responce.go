package handler

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func NewErrorResponse(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	log.Println(msg)
	ans, _ := json.Marshal(map[string]string{
		"status": msg,
	})
	_, err := w.Write(ans)
	if err != nil {
		return
	}
}

func WsErrorResponse(conn *websocket.Conn, msg string) {
	log.Println(msg)
	if err := conn.WriteJSON(map[string]interface{}{
		"status": msg,
	}); err != nil {
		return
	}
}
