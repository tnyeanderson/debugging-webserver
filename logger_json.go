package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// JSONLogger logs the result in JSON format
type JSONLogger struct {
	DefaultLogger
	initFunc         func(*JSONLogger) error
	writeRequestFunc func(*JSONLogger, *http.Request) error
}

// NewJSONLogger initializes and returns a JSONLogger.
func NewJSONLogger() (l *JSONLogger) {
	d := NewDefaultLogger()
	l = &JSONLogger{
		DefaultLogger:    *d,
		initFunc:         func(j *JSONLogger) error { return nil },
		writeRequestFunc: writeRequestJSON,
	}
	return
}

func (l *JSONLogger) Init() error {
	return l.initFunc(l)
}

func (l *JSONLogger) WriteRequest(r *http.Request) error {
	l.defaultLoggerWriteHook()
	return l.writeRequestFunc(l, r)
}

func writeRequestJSON(l *JSONLogger, r *http.Request) error {
	b, _ := l.toJSON(r)
	l.Write(b)
	l.writeNewline()
	return nil
}

func (l *JSONLogger) toJSON(r *http.Request) ([]byte, error) {
	req := newRequest(r, l.getTimestamp())
	req.TotalRequests = l.TotalRequests
	return json.Marshal(req)
}

type request struct {
	Wire          []byte      `json:"wire"`
	Body          []byte      `json:"body"`
	Fragment      string      `json:"fragment"`
	Headers       http.Header `json:"headers"`
	Method        string      `json:"method"`
	Path          string      `json:"path"`
	Query         string      `json:"query"`
	QueryParams   url.Values  `json:"queryParams"`
	ReceivedAt    int64       `json:"receivedAt"`
	TotalRequests int64       `json:"totalRequests"`
	Errors        []string    `json:"errors"`
}

func newRequest(r *http.Request, timestamp time.Time) *request {
	req := request{}
	req.Method = r.Method
	req.Path = r.URL.Path
	req.Query = r.URL.RawQuery
	req.Fragment = r.URL.RawFragment
	req.ReceivedAt = timestamp.Unix()
	req.Headers = r.Header

	// Read the body into both Wire and Body
	wireBuffer := &bytes.Buffer{}
	bodyBuffer := &bytes.Buffer{}
	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, bodyBuffer))
	if err := r.Write(wireBuffer); err != nil {
		req.Errors = append(req.Errors, err.Error())
	}
	req.Body = bodyBuffer.Bytes()
	req.Wire = wireBuffer.Bytes()

	queryParams, err := url.ParseQuery(req.Query)
	if err != nil {
		req.Errors = append(req.Errors, err.Error())
	}
	req.QueryParams = queryParams

	return &req
}
