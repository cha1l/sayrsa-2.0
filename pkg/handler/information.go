package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) GetPublicKeyHandler(w http.ResponseWriter, r *http.Request) {
	selfUser := GetParams(r.Context())
	w.Header().Set("Content-Type", "application/json")

	names := mux.Vars(r)
	username := names["username"]

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

func (h *Handler) GetConversationInfoHandler(w http.ResponseWriter, r *http.Request) {
	selfUser := GetParams(r.Context())
	w.Header().Set("Content-Type", "application/json")

	names := mux.Vars(r)
	id, err := strconv.Atoi(names["id"])
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	conv, err := h.service.Conversations.GetConversationInfo(id)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	ans, err := json.Marshal(conv)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = w.Write(ans)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("user %s got information about conversation with id %d", selfUser, id)
}

func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {

}
