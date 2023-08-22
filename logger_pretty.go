package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// PrettyLogger is like DefaultLogger, but it also prints a separator between
// requests containing the timestamp and TotalRequests
type PrettyLogger struct {
	DefaultLogger

	// Width is the amount of columns used for separator lines. Usually 80 for
	// VT100 reasons.
	Width int
}

func NewPrettyLogger() *PrettyLogger {
	return &PrettyLogger{
		DefaultLogger: *NewDefaultLogger(),
		Width:         80,
	}
}

func (l *PrettyLogger) Init() error {
	return l.writeBanner()
}

func (l *PrettyLogger) WriteRequest(r *http.Request) error {
	l.defaultLoggerWriteHook()
	l.writeSeparator()
	r.Write(l)
	return nil
}

func (l *PrettyLogger) writeBanner() error {
	l.Write([]byte(defaultBanner))
	l.writeNewline()
	l.writeSeparatorLine('+')
	return nil
}

func (l *PrettyLogger) writeSeparatorLine(char rune) {
	l.Write([]byte(strings.Repeat(string(char), l.Width)))
	l.writeNewline()
}

func (l *PrettyLogger) writeSeparatorMessage(char rune, message string) {
	prefix := 3
	// prefix, space, message, space, suffix
	suffixLength := l.Width - 3 - 1 - len(message) - 1
	out := fmt.Sprintf("%s %s %s\n", strings.Repeat(string(char), prefix), message, strings.Repeat(string(char), suffixLength))
	l.Write([]byte(out))
}

func (l *PrettyLogger) writeSeparator() error {
	l.writeNewline()
	l.writeSeparatorLine('*')
	l.writeSeparatorMessage('-', l.getTimestamp().Format(time.UnixDate))
	l.writeSeparatorMessage('-', fmt.Sprintf("Total requests: %d", l.TotalRequests))
	return nil
}
