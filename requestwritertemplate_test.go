package flies

import (
	"html/template"
	"os"
)

func ExampleFliesLogTemplate() {
	tmpl, _ := template.New("").Parse("{{.ReceivedAt}} | {{.Method}} {{.Path}}")
	w := NewRequestWriterTemplate(os.Stdout, tmpl)
	testRequestWriterInit(&w.DefaultRequestWriter)
	req := newTestRequest()
	w.WriteRequest(req)

	// Output:
	// 3735928559 | POST /my/test/path
}
