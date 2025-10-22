package api

import (
	"net/http"
)

// Represents API server
type Server struct {
	router http.Handler
	addr   string
}

// New API server
func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

// starts API server
func (s *Server) Start() error {
	return http.ListenAndServe(s.addr, s.router)
}
