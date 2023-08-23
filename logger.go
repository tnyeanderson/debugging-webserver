package main

import (
	"io"
	"net/http"
)

// Logger is an interface for logging HTTP requests.
type Logger interface {
	io.Writer

	// Init prepares a log destination, usually by writing a banner.
	Init() error

	// WriteRequest writes a single request to the log, along with any
	// prefix/suffix/transformation as defined by the Logger.
	//
	// This is useful if your logger is keeping track of how many requests it has
	// processed, etc. However for some Loggers it may be a no-op, and it can
	// even be bypassed completely by the caller.  For example, TCPServer uses
	// Logger.Write() directly.
	WriteRequest(*http.Request) error
}
