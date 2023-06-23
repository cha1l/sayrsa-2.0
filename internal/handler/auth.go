package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cha1l/sayrsa-2.0/models"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var input models.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := h.service.Authorization.CreateUser(input)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	ans, err := json.Marshal(map[string]interface{}{
		"token": token,
	})
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(ans)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Registation succsed: token created for %s\n", input.Username)

}

type SignInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var input SignInInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	token, privateKey, err := h.service.Authorization.GetUserTokenPrivateKey(input.Username, input.Password)
	if err != nil {
		NewErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	ans, err := json.Marshal(map[string]interface{}{
		"token":      token,
		"privateKey": privateKey,
	})
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(ans)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("User %s signed in", input.Username)
}
