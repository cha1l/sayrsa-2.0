package sayrsa20

import (
	"log"
	"net/http"
	"os"
	"time"
)

type APIServer struct {
	httpServer *http.Server
}

func (s *APIServer) StartServer(handler http.Handler) error {
	AppIp := os.Getenv("APP_IP")
	AppPort := os.Getenv("APP_PORT")

	s.httpServer = &http.Server{
		Addr:           AppIp + ":" + AppPort,
		Handler:        handler,
		ReadTimeout:    1 * time.Minute,
		WriteTimeout:   1 * time.Minute,
		MaxHeaderBytes: 1 << 20, //1mb
	}

	log.Println("Server starting...")

	return s.httpServer.ListenAndServe()
}
