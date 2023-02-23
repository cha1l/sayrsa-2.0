package handler

import (
	"github.com/cha1l/sayrsa-2.0/pkg/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	service service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{}
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()

	//Authorization
	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/sign-in", h.SignIn).Methods("POST")
	auth.HandleFunc("/sign-up", h.SignUp).Methods("POST")

	return r
}
