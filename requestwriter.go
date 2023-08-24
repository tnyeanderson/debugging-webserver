package main

import (
	"bufio"
	"io"
	"net/http"
	"time"
)

type RequestWriter interface {
	// WriteRequest writes a single request to the log, along with any
	// prefix/suffix/transformation as defined by the Logger.
	//
	// This is useful if your logger is keeping track of how many requests it has
	// processed, etc. However for some Loggers it may be a no-op, and it can
	// even be bypassed completely by the caller.  For example, TCPServer uses
	// Logger.Write() directly.
	WriteRequest(*http.Request) error
}

type DefaultRequestWriter struct {
	Out           io.Writer
	Separator     string
	TotalRequests int64

	// getTimestamp returns the time.Time of the current request. Allowed to be
	// set here for tests. Usually time.Now
	getTimestamp func() time.Time
}

func NewRequestWriter(separator string, out io.Writer) *DefaultRequestWriter {
	w := &DefaultRequestWriter{}
	w.Separator = separator
	w.Out = out
	w.getTimestamp = time.Now
	return w
}

func (w *DefaultRequestWriter) WriteRequest(r *http.Request) error {
	w.TotalRequests++
	r.Write(w.Out)
	w.Out.Write([]byte(w.Separator))
	return nil
}

func MultiRequestWriter(errWriter io.Writer, writers ...RequestWriter) io.Writer {
	r, w := io.Pipe()

	go func() {
		for {
			req, err := http.ReadRequest(bufio.NewReader(r))
			if err != nil {
				errWriter.Write([]byte(err.Error()))
				errWriter.Write([]byte("\n"))
				continue
			}
			for _, writer := range writers {
				writer.WriteRequest(req)
			}
		}
	}()

	return w
}
