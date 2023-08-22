package main

// Server listens for HTTP requests and logs them with a Logger.
type Server interface {
	// Init initializes a server, perhaps by reading a config file, environment
	// variables, or prompting the user for information.
	Init()

	// GetLogger returns the Logger for the Server.
	GetLogger() Logger

	// Listen starts a server and sets it up to log requests.
	Listen() error
}
