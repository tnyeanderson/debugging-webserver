package flies

import (
	"io"
	"os"
	"strings"
	"testing"
)

func ExampleFliesLogJSON() {
	w := NewRequestWriterJSON(os.Stdout)
	testRequestWriterInit(&w.DefaultRequestWriter)
	req := newTestRequest()
	w.WriteRequest(req)

	// Output:
	// {"wire":"UE9TVCAvbXkvdGVzdC9wYXRoP3BhcmFtMT12YWx1ZTEmcGFyYW0yJm11bHRpPWZpcnN0dmFsdWUmbXVsdGk9c2Vjb25kdmFsdWUgSFRUUC8xLjENCkhvc3Q6IGV4YW1wbGUuY29tDQpVc2VyLUFnZW50OiBHby1odHRwLWNsaWVudC8xLjENCkNvbnRlbnQtTGVuZ3RoOiAxOQ0KQWNjZXB0LUVuY29kaW5nOiBnemlwLCBkZWZsYXRlDQpBY2NlcHQtTGFuZ3VhZ2U6IGVuLXVzDQpGb286IEJhcg0KRm9vOiB0d28NCg0KdGhpcyBpcyBhIHRlc3QgYm9keQ==","body":"dGhpcyBpcyBhIHRlc3QgYm9keQ==","fragment":"","headers":{"Accept-Encoding":["gzip, deflate"],"Accept-Language":["en-us"],"Foo":["Bar","two"]},"method":"POST","path":"/my/test/path","query":"param1=value1\u0026param2\u0026multi=firstvalue\u0026multi=secondvalue","queryParams":{"multi":["firstvalue","secondvalue"],"param1":["value1"],"param2":[""]},"receivedAt":3735928559,"totalRequests":1,"errors":null}
}

func TestJSONWriteRequestOnePerLine(t *testing.T) {
	out := &strings.Builder{}
	w := NewRequestWriterJSON(out)
	w.WriteRequest(newTestRequest())
	w.WriteRequest(newTestRequest())
	w.WriteRequest(newTestRequest())

	if strings.Count(out.String(), "\n") != 3 {
		t.Fail()
	}
}

func TestJSONLoggerRequestIncrementsTotal(t *testing.T) {
	w := NewRequestWriterJSON(io.Discard)
	w.WriteRequest(newTestRequest())
	w.WriteRequest(newTestRequest())
	w.WriteRequest(newTestRequest())

	if w.TotalRequests != 3 {
		t.Fail()
	}
}

// TODO: Fix this test
//func TestJSONLoggerNoPrintsBanner(t *testing.T) {
//	out := &strings.Builder{}
//	w := NewRequestWriterJSON(out)
//
//	if out.Len() > 0 {
//		t.Fail()
//	}
//}
