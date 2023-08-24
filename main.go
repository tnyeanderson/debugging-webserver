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

func main() {
	s := &TCPServer{}
	s.Init()
	out := os.Stdout
	errOut := os.Stdout
	w := MultiRequestWriter(errOut, NewRequestWriterPretty(out))
	//w := MultiRequestWriter(errOut, NewRequestWriterJSON(out), NewRequestWriter("---", out))
	//w := MultiRequestWriter(errOut, NewRequestWriter("---", out))
	fmt.Println(s.Listen(w, errOut))
}
