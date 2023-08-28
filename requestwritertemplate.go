package flies

import (
	"html/template"
	"io"
	"net/http"
)

// RequestWriterTemplate is a RequestWriter that writes each http.Request
// according to a provided html/template.
type RequestWriterTemplate struct {
	DefaultRequestWriter
	template *template.Template
}

// NewRequestWriterTemplate returns an initialized RequestWriterTemplate.
func NewRequestWriterTemplate(out io.Writer, tmpl *template.Template) *RequestWriterTemplate {
	return &RequestWriterTemplate{
		DefaultRequestWriter: *NewRequestWriter("", out),
		template:             tmpl,
	}
}

// WriteRequest increments TotalRequests, creates a request object based on the
// http.Request, and writes the request according to the template.
func (w *RequestWriterTemplate) WriteRequest(r *http.Request) error {
	w.TotalRequests++
	req := newRequest(r, w.getTimestamp())
	req.TotalRequests = w.TotalRequests
	w.template.Execute(w.Out, req)
	return nil
}
