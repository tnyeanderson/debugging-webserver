package flies

import "io"

// Server listens for HTTP requests and logs them with a Logger.
type Server interface {
	// Init initializes a Server, perhaps by reading a config file, environment
	// variables, or prompting the user for information.
	Init()

	// Listen writes any data received by the server to rawWriter, and writes any
	// errors (followed by a newline) to errWriter. The request is parsed, and
	// passed to reqWriter.
	Listen(errWriter, rawWriter io.Writer, reqWriter RequestWriter) error
}
