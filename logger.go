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
	WritePrefix() error
	WriteSuffix() error
	Handler() http.HandlerFunc
}

type DefaultLogger struct {
	Width           int
	Out             io.Writer
	TotalRequests   int64
	initFunc        func(*DefaultLogger) error
	handlerFunc     func(*DefaultLogger, http.ResponseWriter, *http.Request)
	writePrefixFunc func(*DefaultLogger) error
	writeSuffixFunc func(*DefaultLogger) error
	getTimestamp    func() time.Time
}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		Out:             os.Stdout,
		Width:           80,
		handlerFunc:     handlerText,
		initFunc:        logBanner,
		writePrefixFunc: logSeparator,
		writeSuffixFunc: func(d *DefaultLogger) error { return nil },
		getTimestamp:    time.Now,
	}
}

func (l *DefaultLogger) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.handlerFunc(l, w, r)
	}
}

func (l *DefaultLogger) Init() error {
	return l.initFunc(l)
}

func (l *DefaultLogger) WritePrefix() error {
	return l.writePrefixFunc(l)
}

func (l *DefaultLogger) WriteSuffix() error {
	return l.writeSuffixFunc(l)
}

func (l *DefaultLogger) Write(b []byte) (int, error) {
	return l.Out.Write(b)
}

func handlerText(l *DefaultLogger, w http.ResponseWriter, r *http.Request) {
	l.TotalRequests++
	logSeparator(l)
	r.Write(l)
	l.logNewline()
}

func logBanner(l *DefaultLogger) error {
	l.Write([]byte(defaultBanner))
	l.logNewline()
	l.logSeparatorLine('+')
	return nil
}

func logSeparator(l *DefaultLogger) error {
	l.logSeparatorLine('*')
	l.logSeparatorMessage('-', l.getTimestamp().Format(time.UnixDate))
	l.logSeparatorMessage('-', fmt.Sprintf("Total requests: %d", l.TotalRequests))
	return nil
}

func (l *DefaultLogger) logNewline() {
	l.Write([]byte("\n"))
}

func (l *DefaultLogger) logSeparatorLine(char rune) {
	l.Write([]byte(strings.Repeat(string(char), l.Width)))
	l.logNewline()
}

func (l *DefaultLogger) logSeparatorMessage(char rune, message string) {
	prefix := 3
	// prefix, space, message, space, suffix
	suffixLength := l.Width - 3 - 1 - len(message) - 1
	out := fmt.Sprintf("%s %s %s\n", strings.Repeat(string(char), prefix), message, strings.Repeat(string(char), suffixLength))
	l.Write([]byte(out))
}
