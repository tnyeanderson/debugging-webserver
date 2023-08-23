package main

import (
	"net/http"
)

// HTTPServer is the server used by flies.
type HTTPServer struct {
	TCPServer
}

func (s *HTTPServer) Init(l Logger) {
	s.TCPServer.Init(l)
}

func (s *HTTPServer) Listen() error {
	http.HandleFunc("/", s.Handler())
	return http.ListenAndServe(s.getAddr(), nil)
}

// Handler returns an http.Handler which calls WriteRequest on the Logger.
func (s *HTTPServer) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.GetLogger().WriteRequest(r)
	}
}
