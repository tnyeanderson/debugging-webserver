// Command flies provides a webserver which produces detailed logs of each
// request it receives.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/tnyeanderson/flies"
)

const defaultBanner = `
   __ _ _           
  / _| (_)          
 | |_| |_  ___  ___ 
 |  _| | |/ _ \/ __|
 | | | | |  __/\__ \
 |_| |_|_|\___||___/
`

func getRequestWriter(out io.Writer) flies.RequestWriter {
	switch os.Getenv("FLIES_LOG_FORMAT") {
	case "json":
		return flies.NewRequestWriterJSON(out)
	case "raw":
		return flies.NewRequestWriter("\n", out)
	default:
		return flies.NewRequestWriterPretty(out)
	}
}

func main() {
	s := &flies.TCPServer{}
	s.Init()
	out := os.Stdout
	errOut := os.Stderr
	w := flies.MultiRequestWriter(errOut, getRequestWriter(out))
	fmt.Println(s.Listen(w, errOut))
}
