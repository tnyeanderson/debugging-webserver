package main

import (
	"bytes"
	"fmt"
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

func mockInit(l *DefaultLogger) error {
	l.Write([]byte("<banner>"))
	return nil
}

func mockTimestamp() time.Time {
	return testTimestamp
}

func testLoggerInit(l *DefaultLogger) {
	l.initFunc = mockInit
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
	l.handlerFunc(l, httptest.NewRecorder(), req)

	// Output:
	// {"wire":"UE9TVCAvbXkvdGVzdC9wYXRoP3BhcmFtMT12YWx1ZTEmcGFyYW0yJm11bHRpPWZpcnN0dmFsdWUmbXVsdGk9c2Vjb25kdmFsdWUgSFRUUC8xLjENCkhvc3Q6IGV4YW1wbGUuY29tDQpVc2VyLUFnZW50OiBHby1odHRwLWNsaWVudC8xLjENCkNvbnRlbnQtTGVuZ3RoOiAxOQ0KQWNjZXB0LUVuY29kaW5nOiBnemlwLCBkZWZsYXRlDQpBY2NlcHQtTGFuZ3VhZ2U6IGVuLXVzDQpGb286IEJhcg0KRm9vOiB0d28NCg0KdGhpcyBpcyBhIHRlc3QgYm9keQ==","body":"dGhpcyBpcyBhIHRlc3QgYm9keQ==","fragment":"","headers":{"Accept-Encoding":["gzip, deflate"],"Accept-Language":["en-us"],"Foo":["Bar","two"]},"method":"POST","path":"/my/test/path","query":"param1=value1\u0026param2\u0026multi=firstvalue\u0026multi=secondvalue","queryParams":{"multi":["firstvalue","secondvalue"],"param1":["value1"],"param2":[""]},"receivedAt":3735928559,"totalRequests":1,"errors":null}
}

func ExampleFliesLogText() {
	l := NewDefaultLogger()
	testLoggerInit(l)
	req := newTestRequest()
	buf := &bytes.Buffer{}
	l.Out = buf
	l.handlerFunc(l, httptest.NewRecorder(), req)

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
