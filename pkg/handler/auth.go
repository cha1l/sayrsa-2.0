package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cha1l/sayrsa-2.0/models"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var input models.User
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		NewErrorResponce(w, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := h.service.Authorization.CreateUser(input)
	if err != nil {
		NewErrorResponce(w, http.StatusInternalServerError, err.Error())
		return
	}

	ans, err := json.Marshal(map[string]interface{}{
		"token": token,
	})
	if err != nil {
		NewErrorResponce(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(ans)

	log.Printf("Registation succsed: token created for %s\n", input.Username)

}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {

}
