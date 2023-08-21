package main

import (
	"fmt"
	"net/http"
	"os"
)

const (
	defaultLogFormat = "text"
	defaultPort      = "8080"
)

// defaultServer is the server used by flies.
type defaultServer struct {
	LogFormat string
	Logger    Logger
	Port      string
}

func (s *defaultServer) Init() {
	s.LogFormat = os.Getenv("FLIES_LOG_FORMAT")
	s.Port = os.Getenv("FLIES_PORT")
}

func (s *defaultServer) GetLogger() Logger {
	switch s.LogFormat {
	case "json":
		return NewJSONLogger()
	default:
		return NewDefaultLogger()
	}
}

func (s *defaultServer) Listen() error {
	http.HandleFunc("/", Handler(s.Logger))
	return http.ListenAndServe(s.getAddr(), nil)
}

func (s *defaultServer) getAddr() string {
	port := s.Port
	if port == "" {
		port = defaultPort
	}
	return fmt.Sprintf(":%s", port)
}
