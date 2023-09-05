package flies

import (
	"io"
	"net/http"
)

// MultiRequestWriter is a [RequestWriter] that calls each of the Writers
// sequentially with the same request. Any errors will be written to ErrWriter,
// followed by a newline.
type MultiRequestWriter struct {
	// Writers is the list of RequestWriters which will be called sequentially.
	Writers []RequestWriter

	// ErrorWriter is where errors will be written. Leave nil to return instead,
	// bypassing all subsequent Writers.
	ErrWriter *io.Writer
}

// WriteRequest calls each [RequestWriter] in the Writers slice sequentially
// with the same request. Any errors will be written to ErrWriter, followed by
// a newline.
//
// If the [http.Request.Body] is not already a [BodyReader], it will be
// replaced with one first.
func (r *MultiRequestWriter) WriteRequest(req *http.Request) error {
	if _, ok := req.Body.(*BodyReader); !ok {
		req.Body = NewBodyReader(req.Body)
	}
	for _, w := range r.Writers {
		err := w.WriteRequest(req)
		if err != nil {
			if r.ErrWriter != nil {
				writeError(*r.ErrWriter, err)
				continue
			}
			return err
		}
	}
	return nil
}

// NewMultiRequestWriter returns an initialized [MultiRequestWriter].
func NewMultiRequestWriter(errWriter io.Writer, writers ...RequestWriter) RequestWriter {
	return &MultiRequestWriter{
		Writers:   writers,
		ErrWriter: &errWriter,
	}
}

func writeError(w io.Writer, err error) {
	w.Write([]byte(err.Error()))
	w.Write([]byte("\n"))
}
