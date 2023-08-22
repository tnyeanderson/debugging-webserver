package main

import (
	"net"
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
	s.Logger = s.GetLogger()
}

func (s *defaultServer) GetLogger() Logger {
	if s.Logger != nil {
		return s.Logger
	}
	switch s.LogFormat {
	case "json":
		return NewJSONLogger()
	default:
		return NewPrettyLogger()
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
	return net.JoinHostPort("", port)
}
