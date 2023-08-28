package flies

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

// RequestWriterJSON is a RequestWriter that writes each http.Request in JSON
// format. Each request is written as a single line.
type RequestWriterJSON struct {
	DefaultRequestWriter
}

// NewRequestWriterJSON returns an initialized RequestWriterJSON.
func NewRequestWriterJSON(out io.Writer) *RequestWriterJSON {
	return &RequestWriterJSON{
		DefaultRequestWriter: *NewRequestWriter("", out),
	}
}

// WriteRequest increments TotalRequests and writes the request in one-line
// JSON format, followed by a newline character.
func (w *RequestWriterJSON) WriteRequest(r *http.Request) error {
	w.TotalRequests++
	req := newRequest(r, w.getTimestamp())
	req.TotalRequests = w.TotalRequests
	b, _ := json.Marshal(req)
	w.Out.Write(b)
	w.Out.Write([]byte("\n"))
	return nil
}

// request is the JSON representation of a request
type request struct {
	Wire          []byte      `json:"wire"`
	Body          []byte      `json:"body"`
	Fragment      string      `json:"fragment"`
	Headers       http.Header `json:"headers"`
	Method        string      `json:"method"`
	Path          string      `json:"path"`
	Query         string      `json:"query"`
	QueryParams   url.Values  `json:"queryParams"`
	ReceivedAt    int64       `json:"receivedAt"`
	TotalRequests int64       `json:"totalRequests"`
	Errors        []string    `json:"errors"`
}

func newRequest(r *http.Request, timestamp time.Time) *request {
	req := request{}
	req.Method = r.Method
	req.Path = r.URL.Path
	req.Query = r.URL.RawQuery
	req.Fragment = r.URL.RawFragment
	req.ReceivedAt = timestamp.Unix()
	req.Headers = r.Header

	if _, ok := r.Body.(*BodyReader); !ok {
		r.Body = NewBodyReader(r.Body)
	}

	wireBuffer := &bytes.Buffer{}
	if err := r.Write(wireBuffer); err != nil {
		req.Errors = append(req.Errors, err.Error())
	}
	req.Wire = wireBuffer.Bytes()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		req.Errors = append(req.Errors, err.Error())
	}
	req.Body = b

	// Parse query parameters
	queryParams, err := url.ParseQuery(req.Query)
	if err != nil {
		req.Errors = append(req.Errors, err.Error())
	}
	req.QueryParams = queryParams

	return &req
}
