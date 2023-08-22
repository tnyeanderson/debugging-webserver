package main

import (
	"io"
	"net/http"
	"os"
	"time"
)

// DefaultLogger logs the result in wire format to Out (stdout by default) and
// maintains a count of how many requests it has processed.
type DefaultLogger struct {
	// Out is the destination for the log, usually stdout.
	Out io.Writer

	// TotalRequests is an Incrementing counter of times WriteRequest has been
	// called.
	TotalRequests int64

	// getTimestamp returns the time.Time of the current request. Allowed to be
	// set here for tests. Usually time.Now
	getTimestamp func() time.Time
}

// NewDefaultLogger initializes and returns a DefaultLogger.
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		Out:          os.Stdout,
		getTimestamp: time.Now,
	}
}

func (l *DefaultLogger) Init() error { return nil }

func (l *DefaultLogger) WriteRequest(r *http.Request) error {
	l.defaultLoggerWriteHook()
	l.writeNewline()
	r.Write(l)
	return nil
}

func (l *DefaultLogger) Write(b []byte) (int, error) {
	return l.Out.Write(b)
}

func (l *DefaultLogger) writeNewline() {
	l.Write([]byte("\n"))
}

func (l *DefaultLogger) defaultLoggerWriteHook() {
	l.TotalRequests++
}
