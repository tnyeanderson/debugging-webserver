package flies

import "net/http"

// RequestWriterDiscard is a no-op RequestWriter.
type RequestWriterDiscard struct{}

// WriteRequest is a no-op.
func (r *RequestWriterDiscard) WriteRequest(req *http.Request) error {
	return nil
}
