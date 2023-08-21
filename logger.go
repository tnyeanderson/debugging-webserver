package main

import (
	"io"
	"net/http"
)

// Logger is an interface for logging HTTP requests.
type Logger interface {
	io.Writer

	// Init prepares a log destination, usually by displaying a banner.
	Init() error

	// WriteRequest converts an http.Request into the Logger's specified format
	// and calls Write with the result.
	WriteRequest(*http.Request) error
}
