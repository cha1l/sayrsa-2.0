package handler

import (
	"net/http"
)

func (h *Handler) CreateConversation(w http.ResponseWriter, r *http.Request) {
	userId := GetParams(r.Context())
	if userId == -1 {
		NewErrorResponce(w, http.StatusUnauthorized, "error getting params")
		return
	}
}
