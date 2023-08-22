package main

import (
	"io"
	"log"
	"net"
	"os"
)

const (
	defaultLogFormat = "pretty"
	defaultPort      = "8080"
)

type TCPServer struct {
	LogFormat string
	Logger    Logger
	Port      string
}

func (s *TCPServer) Init() {
	s.LogFormat = os.Getenv("FLIES_LOG_FORMAT")
	s.Port = os.Getenv("FLIES_PORT")
	s.Logger = s.GetLogger()
}

func (s *TCPServer) GetLogger() Logger {
	if s.Logger != nil {
		return s.Logger
	}
	return NewDefaultLogger()
}

func (s *TCPServer) Listen() error {
	l, err := net.Listen("tcp", s.getAddr())
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			// Echo all incoming data.
			io.Copy(s.GetLogger(), c)
			// Shut down the connection.
			c.Close()
		}(conn)
	}
	return err
}

func (s *TCPServer) getAddr() string {
	port := s.Port
	if port == "" {
		port = defaultPort
	}
	return net.JoinHostPort("", port)
}
