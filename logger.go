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
	LogBanner()
	LogSeparator()
	Handler() http.HandlerFunc
}

type DefaultLogger struct {
	Width         int
	Out           io.Writer
	TotalRequests int64
	HandlerFunc   func(*DefaultLogger, http.ResponseWriter, *http.Request)
	LogBannerFunc func(*DefaultLogger)
	GetTimestamp  func() time.Time
}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		Out:           os.Stdout,
		Width:         80,
		HandlerFunc:   handlerText,
		LogBannerFunc: logBanner,
		GetTimestamp:  time.Now,
	}
}

func (l *DefaultLogger) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.HandlerFunc(l, w, r)
	}
}

func (l *DefaultLogger) LogBanner() {
	l.LogBannerFunc(l)
}

func (l *DefaultLogger) LogSeparator() {
	l.logSeparatorLine('*')
	l.logSeparatorMessage('-', l.GetTimestamp().Format(time.UnixDate))
	l.logSeparatorMessage('-', fmt.Sprintf("Total requests: %d", l.TotalRequests))
}

func (l *DefaultLogger) Write(b []byte) (int, error) {
	return l.Out.Write(b)
}

func handlerText(l *DefaultLogger, w http.ResponseWriter, r *http.Request) {
	l.TotalRequests++
	l.LogSeparator()
	r.Write(l)
	l.logNewline()
}

func logBanner(l *DefaultLogger) {
	l.Write([]byte(defaultBanner))
	l.logNewline()
	l.logSeparatorLine('+')
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
