package sayrsa20

import (
	"log"
	"net/http"
	"time"
)

type APIServer struct {
	httpServer *http.Server
}

func (s *APIServer) StartServer(handler http.Handler, ip string, port string) error {
	s.httpServer = &http.Server{
		Addr:           ip + ":" + port,
		Handler:        handler,
		ReadTimeout:    1 * time.Minute,
		WriteTimeout:   1 * time.Minute,
		MaxHeaderBytes: 1 << 20, //1mb
	}

	log.Println("Server starting on", port)

	return s.httpServer.ListenAndServe()
}
