package flies

import (
	"bufio"
	"io"
	"net/http"
)

// MultiRequestWriter returns a Writer which, when written to, will parse the
// request it receives and call each RequestWriter sequentially with the
// result.  Any errors will be written to errWriter, followed by a newline.
//
// The Body of the http.Request is replaced with a BodyReader, which allows
// reading the body multiple times (so each RequestWriter in the chain can read
// the body as needed).
func MultiRequestWriter(errWriter io.Writer, writers ...RequestWriter) io.Writer {
	r, w := io.Pipe()

	go func() {
		for {
			req, err := http.ReadRequest(bufio.NewReader(r))
			if err != nil {
				writeError(errWriter, err)
				continue
			}
			req.Body = NewBodyReader(req.Body)
			for _, writer := range writers {
				writer.WriteRequest(req)
			}
		}
	}()

	return w
}

func writeError(w io.Writer, err error) {
	w.Write([]byte(err.Error()))
	w.Write([]byte("\n"))
}
