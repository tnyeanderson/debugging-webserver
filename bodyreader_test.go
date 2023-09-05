package flies

import (
	"fmt"
	"io"
	"os"
)

func ExampleBodyReader() {
	req := newTestRequest()
	req.Body = NewBodyReader(req.Body)

	fmt.Println("Read body once:")
	io.Copy(os.Stdout, req.Body)
	req.Body.Close()
	fmt.Println()

	fmt.Println("Read body again:")
	io.Copy(os.Stdout, req.Body)
	req.Body.Close()
	fmt.Println()

	// Output:
	// Read body once:
	// this is a test body
	// Read body again:
	// this is a test body
}
