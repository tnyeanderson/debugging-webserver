package flies

import (
	"io"
	"net/http"
)

// MultiRequestWriter is a RequestWriter that calls multiple RequestWriters
// sequentially with the same request. Any errors will be written to errWriter,
// followed by a newline.
type MultiRequestWriter struct {
	// writers is the list of RequestWriters which will be called sequentially.
	writers []RequestWriter

	// errorWriter is where errors will be written. Leave nil to return instead,
	// bypassing all subsequent RequestWriters.
	errWriter *io.Writer
}

// WriteRequest calls each RequestWriter in the writers slice sequentially with
// the same request. Any errors will be written to errWriter, followed by a
// newline.
//
// If the Body of the http.Request is not already a BodyReader, it will be
// replaced with one first.
func (r *MultiRequestWriter) WriteRequest(req *http.Request) error {
	if _, ok := req.Body.(*BodyReader); !ok {
		req.Body = NewBodyReader(req.Body)
	}
	for _, w := range r.writers {
		err := w.WriteRequest(req)
		if err != nil {
			if r.errWriter != nil {
				writeError(*r.errWriter, err)
				continue
			}
			return err
		}
	}
	return nil
}

// NewMultiRequestWriter returns an initialized MultiRequestWriter.
func NewMultiRequestWriter(errWriter io.Writer, writers ...RequestWriter) RequestWriter {
	return &MultiRequestWriter{
		writers:   writers,
		errWriter: &errWriter,
	}
}

func writeError(w io.Writer, err error) {
	w.Write([]byte(err.Error()))
	w.Write([]byte("\n"))
}
