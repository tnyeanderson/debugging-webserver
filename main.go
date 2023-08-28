package main

import (
	"fmt"
	"io"
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

func getRequestWriter(out io.Writer) RequestWriter {
	switch os.Getenv("FLIES_LOG_FORMAT") {
	case "json":
		return NewRequestWriterJSON(out)
	case "raw":
		return NewRequestWriter("\n", out)
	default:
		return NewRequestWriterPretty(out)
	}
}

func main() {
	s := &TCPServer{}
	s.Init()
	out := os.Stdout
	errOut := os.Stderr
	w := MultiRequestWriter(errOut, getRequestWriter(out))
	fmt.Println(s.Listen(w, errOut))
}
