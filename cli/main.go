// Command flies provides a webserver which produces detailed logs of each
// request it receives.
package main

import (
	"html/template"
	"io"
	"log"
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
	switch os.Getenv("FLIES_FORMAT") {
	case "json":
		return flies.NewRequestWriterJSON(out)
	case "raw":
		return flies.NewRequestWriter("\n", out)
	case "template":
		templateFile := os.Getenv("FLIES_TEMPLATE")
		if templateFile == "" {
			log.Fatal("FLIES_TEMPLATE_FILE must be set if FLIES_FORMAT=template")
		}
		templateText, err := os.ReadFile(templateFile)
		if err != nil {
			log.Fatal(err.Error())
		}
		tmpl, err := template.New("").Parse(string(templateText))
		if err != nil {
			log.Fatal(err.Error())
		}
		return flies.NewRequestWriterTemplate(out, tmpl)
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
	err := s.Listen(w, errOut)
	if err != nil {
		log.Fatal(err.Error())
	}
}
