package flies

import (
	"io"
	"net/http"
	"time"
)

// RequestWriter is something that takes an *http.Request, marshals or
// transforms it in some way, and writes the result somewhere.
//
// The common use case is to daisy chain RequestWriters into a
// MultiRequestWriter, which automatically wraps the Body of the http.Request
// to a BodyReader. If calling multiple RequestWriters on the same request
// without using MultiRequestWriter and need to be able to read the Body more
// than once, be sure to use a BodyReader.
type RequestWriter interface {
	// WriteRequest writes a single request to the log, along with any
	// prefix/suffix/transformation as defined by the RequestLogger.
	WriteRequest(*http.Request) error
}

// DefaultRequestWriter will write each request, followed by the provided
// separator, to the provided io.Writer. It also keeps track of the number of
// requests it has written.
type DefaultRequestWriter struct {
	// Out is the Writer that the RequestWriter writes the request to.
	Out io.Writer

	// Separator is the string that will be written after each request.
	Separator string

	// TotalRequests is the amount of requests that have been received/written by
	// the RequestWriter.
	TotalRequests int64

	// getTimestamp returns the time.Time of the current request. Allowed to be
	// set here for tests. Usually time.Now.
	getTimestamp func() time.Time
}

// NewRequestWriter returns an initialized DefaultRequestWriter.
func NewRequestWriter(separator string, out io.Writer) *DefaultRequestWriter {
	w := &DefaultRequestWriter{}
	w.Separator = separator
	w.Out = out
	w.getTimestamp = time.Now
	return w
}

// WriteRequest increments the TotalRequests counter, writes the Request in
// wire format, then writes the defined separator.
func (w *DefaultRequestWriter) WriteRequest(r *http.Request) error {
	w.TotalRequests++
	r.Write(w.Out)
	w.Out.Write([]byte(w.Separator))
	return nil
}
