package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (h *Handler) GetPublicKeyHandler(w http.ResponseWriter, r *http.Request) {
	selfUser := GetParams(r.Context())

	names := mux.Vars(r)
	username := names["username"]

	w.Header().Set("Content-Type", "application/json")

	publicKey, err := h.service.Conversations.GetPublicKey(username)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	ans, err := json.Marshal(map[string]interface{}{
		"public_key": publicKey,
	})
	_, err = w.Write(ans)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("%s got %s token ", selfUser, username)
}
