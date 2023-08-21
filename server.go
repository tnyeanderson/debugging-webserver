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

type Server interface {
	Init()
	GetLogger() Logger
	Listen() error
}

type DefaultServer struct {
	LogFormat string
	Logger    Logger
	Port      string
}

func (s *DefaultServer) Init() {
	s.LogFormat = os.Getenv("FLIES_LOG_FORMAT")
	s.Port = os.Getenv("FLIES_PORT")
}

func (s *DefaultServer) GetLogger() Logger {
	switch s.LogFormat {
	case "json":
		return NewJSONLogger()
	default:
		return NewDefaultLogger()
	}
}

func (s *DefaultServer) Listen() error {
	http.HandleFunc("/", Handler(s.Logger))
	return http.ListenAndServe(s.getAddr(), nil)
}

func (s *DefaultServer) getAddr() string {
	port := s.Port
	if port == "" {
		port = defaultPort
	}
	return fmt.Sprintf(":%s", port)
}
