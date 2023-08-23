package main

import (
	"fmt"
	"os"
)

const defaultBanner = `
   __ _ _           
  / _| (_)          
 | |_| |_  ___  ___ 
 |  _| | |/ _ \/ __|
 | | | | |  __/\__ \
 |_| |_|_|\___||___/
`

func getServer(protocol string) Server {
	switch protocol {
	case "tcp":
		return &TCPServer{}
	default:
		return &HTTPServer{}
	}
}

func getLogger(logFormat string) Logger {
	switch logFormat {
	case "json":
		return NewJSONLogger()
	default:
		return NewPrettyLogger()
	}
}

func main() {
	protocol := os.Getenv("FLIES_PROTOCOL")
	logFormat := os.Getenv("FLIES_LOG_FORMAT")
	l := getLogger(logFormat)
	l.Init()
	s := getServer(protocol)
	s.Init(l)
	fmt.Println(s.Listen())
}
