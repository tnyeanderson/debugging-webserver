package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

var testTimestamp = time.Unix(3735928559, 0)

var headers = map[string][]string{
	"Accept-Encoding": {"gzip, deflate"},
	"Accept-Language": {"en-us"},
	"Foo":             {"Bar", "two"},
}

func mockTimestamp() time.Time {
	return testTimestamp
}

func testLoggerInit(l *DefaultLogger) {
	l.getTimestamp = mockTimestamp
}

func newTestRequest() *http.Request {
	body := "this is a test body"
	target := "/my/test/path?param1=value1&param2&multi=firstvalue&multi=secondvalue"
	req := httptest.NewRequest("POST", target, strings.NewReader(body))
	// Lifted from the net/http docs
	req.Header = headers
	return req
}
