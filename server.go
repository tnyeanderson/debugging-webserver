package main

import "io"

// Server listens for HTTP requests and logs them with a Logger.
type Server interface {
	// Init initializes a Server, perhaps by reading a config file, environment
	// variables, or prompting the user for information.
	Init()

	// Listen starts a server and writes any data it receives to out. If an error
	// is received, it should be written to errOut, followed by a newline.
	Listen(out, errOut io.Writer) error
}
