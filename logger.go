package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Logger interface {
	io.Writer
	Init() error
	WriteRequest(*http.Request) error
}

type DefaultLogger struct {
	Width            int
	Out              io.Writer
	TotalRequests    int64
	initFunc         func(*DefaultLogger) error
	writeRequestFunc func(*DefaultLogger, *http.Request) error
	getTimestamp     func() time.Time
}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		Out:              os.Stdout,
		Width:            80,
		initFunc:         writeBanner,
		writeRequestFunc: writeRequestText,
		getTimestamp:     time.Now,
	}
}

func (l *DefaultLogger) Init() error {
	return l.initFunc(l)
}

func (l *DefaultLogger) WriteRequest(r *http.Request) error {
	l.defaultLoggerWriteHook()
	return l.writeRequestFunc(l, r)
}

func (l *DefaultLogger) Write(b []byte) (int, error) {
	return l.Out.Write(b)
}

func writeRequestText(l *DefaultLogger, r *http.Request) error {
	writeSeparator(l)
	r.Write(l)
	l.writeNewline()
	return nil
}

func writeBanner(l *DefaultLogger) error {
	l.Write([]byte(defaultBanner))
	l.writeNewline()
	l.writeSeparatorLine('+')
	return nil
}

func writeSeparator(l *DefaultLogger) error {
	l.writeSeparatorLine('*')
	l.writeSeparatorMessage('-', l.getTimestamp().Format(time.UnixDate))
	l.writeSeparatorMessage('-', fmt.Sprintf("Total requests: %d", l.TotalRequests))
	return nil
}

func (l *DefaultLogger) writeNewline() {
	l.Write([]byte("\n"))
}

func (l *DefaultLogger) writeSeparatorLine(char rune) {
	l.Write([]byte(strings.Repeat(string(char), l.Width)))
	l.writeNewline()
}

func (l *DefaultLogger) writeSeparatorMessage(char rune, message string) {
	prefix := 3
	// prefix, space, message, space, suffix
	suffixLength := l.Width - 3 - 1 - len(message) - 1
	out := fmt.Sprintf("%s %s %s\n", strings.Repeat(string(char), prefix), message, strings.Repeat(string(char), suffixLength))
	l.Write([]byte(out))
}

func (l *DefaultLogger) defaultLoggerWriteHook() {
	l.TotalRequests++
}
