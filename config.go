package main

import (
	"fmt"
	"os"
)

const (
	defaultLogFormat = "text"
	defaultPort      = "8080"
)

type config struct {
	LogFormat string
	Port      string
}

func (c *config) Init() {
	c.LogFormat = os.Getenv("FLIES_LOG_FORMAT")
	c.Port = os.Getenv("FLIES_PORT")
}

func (c *config) GetLogger() Logger {
	switch c.LogFormat {
	case "json":
		return NewJSONLogger()
	default:
		return NewDefaultLogger()
	}
}

func (c *config) GetAddr() string {
	port := c.Port
	if port == "" {
		port = defaultPort
	}
	return fmt.Sprintf(":%s", port)
}
