package main

import (
	"io"
	"testing"
)

func TestDefaultLoggerRequestIncrementsTotal(t *testing.T) {
	w := NewRequestWriter("", io.Discard)
	req := newTestRequest()
	w.Out = io.Discard
	w.WriteRequest(req)
	w.WriteRequest(req)
	w.WriteRequest(req)

	if w.TotalRequests != 3 {
		t.Fail()
	}
}
