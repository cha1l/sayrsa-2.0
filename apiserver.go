package sayrsa20

import (
	"log"
	"net/http"
	"time"
)

type APIServer struct {
	httpServer *http.Server
}

func (s *APIServer) StartServer(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		ReadTimeout:    1 * time.Minute,
		WriteTimeout:   1 * time.Minute,
		MaxHeaderBytes: 1 << 20, //1mb
	}

	log.Println("Server starting...")

	return s.httpServer.ListenAndServe()
}
