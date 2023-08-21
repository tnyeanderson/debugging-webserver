package main

import (
	"fmt"
	"net/http"
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

type config struct {
	LogFormat string
}

func loggerFromConfig(c *config) Logger {
	switch c.LogFormat {
	case "json":
		return NewJSONLogger()
	default:
		return NewDefaultLogger()
	}
}

func (c *config) Init() {
	c.LogFormat = os.Getenv("FLIES_LOG_FORMAT")
}

// handler calls WriteRequest on the Logger
func handler(l Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.WriteRequest(r)
	}
}

func main() {
	c := &config{}
	c.Init()
	l := loggerFromConfig(c)
	l.Init()
	http.HandleFunc("/", handler(l))
	fmt.Println(http.ListenAndServe(":8080", nil))
}
