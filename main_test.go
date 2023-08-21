package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

func ExampleFliesLogJSON() {
	l := NewJSONLogger()
	testLoggerInit(&l.DefaultLogger)
	req := newTestRequest()
	l.WriteRequest(req)

	// Output:
	// {"wire":"UE9TVCAvbXkvdGVzdC9wYXRoP3BhcmFtMT12YWx1ZTEmcGFyYW0yJm11bHRpPWZpcnN0dmFsdWUmbXVsdGk9c2Vjb25kdmFsdWUgSFRUUC8xLjENCkhvc3Q6IGV4YW1wbGUuY29tDQpVc2VyLUFnZW50OiBHby1odHRwLWNsaWVudC8xLjENCkNvbnRlbnQtTGVuZ3RoOiAxOQ0KQWNjZXB0LUVuY29kaW5nOiBnemlwLCBkZWZsYXRlDQpBY2NlcHQtTGFuZ3VhZ2U6IGVuLXVzDQpGb286IEJhcg0KRm9vOiB0d28NCg0KdGhpcyBpcyBhIHRlc3QgYm9keQ==","body":"dGhpcyBpcyBhIHRlc3QgYm9keQ==","fragment":"","headers":{"Accept-Encoding":["gzip, deflate"],"Accept-Language":["en-us"],"Foo":["Bar","two"]},"method":"POST","path":"/my/test/path","query":"param1=value1\u0026param2\u0026multi=firstvalue\u0026multi=secondvalue","queryParams":{"multi":["firstvalue","secondvalue"],"param1":["value1"],"param2":[""]},"receivedAt":3735928559,"totalRequests":1,"errors":null}
}

func ExampleFliesLogText() {
	l := NewDefaultLogger()
	testLoggerInit(l)
	req := newTestRequest()
	buf := &bytes.Buffer{}
	l.Out = buf
	l.WriteRequest(req)

	// NOTE: Wire format is \r\n, but example output here is \n
	out := strings.ReplaceAll(string(buf.Bytes()), "\r\n", "\n")
	fmt.Println(out)

	// Output:
	// ********************************************************************************
	// --- Thu May 20 17:55:59 EDT 2088 -----------------------------------------------
	// --- Total requests: 1 ----------------------------------------------------------
	// POST /my/test/path?param1=value1&param2&multi=firstvalue&multi=secondvalue HTTP/1.1
	// Host: example.com
	// User-Agent: Go-http-client/1.1
	// Content-Length: 19
	// Accept-Encoding: gzip, deflate
	// Accept-Language: en-us
	// Foo: Bar
	// Foo: two
	//
	// this is a test body
}

func TestJSONWriteRequestIncrementsTotal(t *testing.T) {
	l := NewJSONLogger()
	testLoggerInit(&l.DefaultLogger)
	req := newTestRequest()
	l.Out = io.Discard
	l.WriteRequest(req)
	l.WriteRequest(req)
	l.WriteRequest(req)

	if l.TotalRequests != 3 {
		t.Fail()
	}
}

func TestTextWriteRequestIncrementsTotal(t *testing.T) {
	l := NewDefaultLogger()
	testLoggerInit(l)
	req := newTestRequest()
	l.Out = io.Discard
	l.WriteRequest(req)
	l.WriteRequest(req)
	l.WriteRequest(req)

	if l.TotalRequests != 3 {
		t.Fail()
	}
}

func TestTextLoggerPrintsBanner(t *testing.T) {
	l := NewDefaultLogger()
	out := &strings.Builder{}
	l.Out = out
	l.Init()

	if !strings.Contains(out.String(), defaultBanner) {
		t.Fail()
	}
}

func TestTextLoggerNoPrintsBanner(t *testing.T) {
	l := NewJSONLogger()
	out := &strings.Builder{}
	l.Out = out
	l.Init()

	if out.Len() > 0 {
		t.Fail()
	}
}
