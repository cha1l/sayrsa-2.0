package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func NewErrorResponce(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	log.Println(msg)
	ans, _ := json.Marshal(map[string]string{
		"status": msg,
	})
	w.Write(ans)
}
