package main

import (
	"io"
	"log"
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
	return net.JoinHostPort("", port)
}

// ------------- TCP -----------------

type tcpServer struct {
	defaultServer
}

func (s *tcpServer) Init() {
	s.LogFormat = os.Getenv("FLIES_LOG_FORMAT")
	s.Port = os.Getenv("FLIES_PORT")
	s.Logger = s.GetLogger()
}

func (s *tcpServer) GetLogger() Logger {
	if s.Logger != nil {
		return s.Logger
	}
	return newTCPLogger()
}

func (s *tcpServer) Listen() error {
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

type tcpLogger struct {
	DefaultLogger
}

// NewTCPLogger initializes and returns a TCPLogger.
func newTCPLogger() (l *tcpLogger) {
	d := NewDefaultLogger()
	l = &tcpLogger{
		DefaultLogger: *d,
	}
	return
}

func (l *tcpLogger) Init() error { return nil }

func (l *tcpLogger) WriteRequest(*http.Request) error { return nil }
