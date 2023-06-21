package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/cha1l/sayrsa-2.0/pkg/service"
	"github.com/gorilla/mux"
)

const (
	createConversationAction  = "create_conv"
	sendMessageAction         = "send_message"
	getAllConversationsAction = "get_all_user_conv"
)

type Handler struct {
	service *service.Service
	clients map[string]Client
}

func New(service *service.Service) *Handler {
	return &Handler{
		service: service,
		clients: make(map[string]Client),
	}
}

type Client struct {
	connection *websocket.Conn
}

func NewClient(conn *websocket.Conn) Client {
	return Client{
		connection: conn,
	}
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()

	//Authorization
	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/sign-up", h.SignUp).Methods(http.MethodPost)
	auth.HandleFunc("/sign-in", h.SignIn).Methods(http.MethodPost)

	//Main api handler
	api := r.PathPrefix("/api").Subrouter()
	api.Use(h.AuthorizationMiddleware)
	api.HandleFunc("/public-key/{username}", h.GetPublicKeyHandler).Methods(http.MethodGet)
	api.HandleFunc("/msg/{convID:[0-9]+}/", h.GetMessages).
		Queries("offset", "{offset}", "amount", "{amount}").
		Methods(http.MethodGet) //todo : into ws
	api.HandleFunc("/conv", h.GetAllConversations).Methods(http.MethodGet)

	//WebSockets handler
	api.HandleFunc("/", h.wsHandler)

	//Test endpoint (home page endpoint in the future)
	r.HandleFunc("/", h.TestEndpoint).Methods(http.MethodGet)

	return r
}

func (h *Handler) TestEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("here")

	body, err := json.Marshal(map[string]interface{}{
		"status": "ok",
	})
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if _, err = w.Write(body); err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
