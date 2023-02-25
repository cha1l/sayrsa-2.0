package handler

import (
	"net/http"

	"github.com/cha1l/sayrsa-2.0/pkg/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()

	//Authorization
	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/sign-up", h.SignUp).Methods(http.MethodPost)
	auth.HandleFunc("/sign-in", h.SignIn).Methods(http.MethodPost)

	//Main api
	api := r.PathPrefix("/api").Subrouter()
	api.Use(h.AuthorizationMiddleware)
	api.HandleFunc("/create-conv", h.CreateConversation).Methods(http.MethodPost)

	return r
}
