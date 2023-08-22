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

func getLogger(logFormat string) Logger {
	switch logFormat {
	case "json":
		return NewJSONLogger()
	default:
		return NewDefaultLogger()
	}
}

func getServer(protocol string) Server {
	switch protocol {
	case "tcp":
		return &tcpServer{}
	default:
		return &defaultServer{}
	}
}

func main() {
	protocol := os.Getenv("FLIES_PROTOCOL")
	s := getServer(protocol)
	s.Init()
	s.GetLogger().Init()
	fmt.Println(s.Listen())
}
