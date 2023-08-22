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

func main() {
	protocol := os.Getenv("FLIES_PROTOCOL")
	s := getServer(protocol)
	s.Init()
	s.GetLogger().Init()
	fmt.Println(s.Listen())
}
