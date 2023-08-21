package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

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
