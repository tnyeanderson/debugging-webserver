package main

import (
	"bytes"
	"fmt"
	"strings"
)

func ExampleFliesLogPretty() {
	buf := &bytes.Buffer{}
	w := NewRequestWriterPretty(buf)
	testRequestWriterInit(&w.DefaultRequestWriter)
	req := newTestRequest()
	w.WriteRequest(req)

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

// TODO: Fix this test
//func TestPrettyLoggerPrintsBanner(t *testing.T) {
//	buf := &bytes.Buffer{}
//	w := NewRequestWriterPretty(buf)
//
//	if !strings.Contains(out.String(), defaultBanner) {
//		t.Fail()
//	}
//}
