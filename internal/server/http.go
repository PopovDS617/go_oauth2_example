package server

import (
	"fmt"
	"net/http"
	"os"
)

type Server interface {
	Run() error
}

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(router http.Handler) *HTTPServer {
	httpPort := os.Getenv("HTTP_PORT")

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", httpPort),
		Handler: router,
	}

	return &HTTPServer{
		server,
	}

}

func (s *HTTPServer) Run() error {
	return s.server.ListenAndServe()

}
