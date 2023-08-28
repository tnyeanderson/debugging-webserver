package flies

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// RequestWriterPretty is like DefaultLogger, but it also prints a separator between
// requests containing the timestamp and TotalRequests
type RequestWriterPretty struct {
	DefaultRequestWriter

	// Width is the amount of columns used for separator lines. Usually 80 for
	// VT100 reasons.
	Width int
}

// NewRequestWriterPretty returns an initialized RequestWriterPretty.
func NewRequestWriterPretty(out io.Writer) *RequestWriterPretty {
	return &RequestWriterPretty{
		DefaultRequestWriter: *NewRequestWriter("", out),
		Width:                80,
	}
}

// Write requests increments TotalRequests, writes the separator banner
// including the timestamp of the request and the amount of requests received,
// then writes the request in wire format.
func (w *RequestWriterPretty) WriteRequest(r *http.Request) error {
	w.TotalRequests++
	w.writeSeparator()
	r.Write(w.Out)
	return nil
}

func (w *RequestWriterPretty) writeSeparatorLine(char rune) {
	w.Out.Write([]byte(strings.Repeat(string(char), w.Width)))
	w.Out.Write([]byte("\n"))
}

func (w *RequestWriterPretty) writeSeparatorMessage(char rune, message string) {
	prefix := 3
	// prefix, space, message, space, suffix
	suffixLength := w.Width - 3 - 1 - len(message) - 1
	out := fmt.Sprintf("%s %s %s\n", strings.Repeat(string(char), prefix), message, strings.Repeat(string(char), suffixLength))
	w.Out.Write([]byte(out))
}

func (w *RequestWriterPretty) writeSeparator() error {
	w.Out.Write([]byte("\n"))
	w.writeSeparatorLine('*')
	w.writeSeparatorMessage('-', w.getTimestamp().Format(time.UnixDate))
	w.writeSeparatorMessage('-', fmt.Sprintf("Total requests: %d", w.TotalRequests))
	return nil
}
