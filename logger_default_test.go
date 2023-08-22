package main

import (
	"io"
	"testing"
)

func TestDefaultLoggerRequestIncrementsTotal(t *testing.T) {
	l := NewDefaultLogger()
	testLoggerInit(l)
	l.Init()
	req := newTestRequest()
	l.Out = io.Discard
	l.WriteRequest(req)
	l.WriteRequest(req)
	l.WriteRequest(req)

	if l.TotalRequests != 3 {
		t.Fail()
	}
}
