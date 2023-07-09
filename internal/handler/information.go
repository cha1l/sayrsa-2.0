package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = w.Write(ans)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("%s got %s token ", selfUser, username)
}

func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	selfUser := GetParams(r.Context())
	w.Header().Set("Content-Type", "application/json")

	convID, err := strconv.Atoi(mux.Vars(r)["convID"])
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	offset, err := strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	amount, err := strconv.Atoi(r.FormValue("amount"))
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	messages, err := h.service.Messages.GetMessages(selfUser, convID, offset, amount)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	ans, err := json.Marshal(map[string]interface{}{
		"messages": messages,
	})
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = w.Write(ans)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) GetAllConversations(w http.ResponseWriter, r *http.Request) {

	username := GetParams(r.Context())

	conversations, err := h.service.GetAllConversations(username)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := map[string]interface{}{
		"conversations": conversations,
	}

	resp, err := json.Marshal(data)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if _, err := w.Write(resp); err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
