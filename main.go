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

func main() {
	c := &config{}
	c.Init()
	l := loggerFromConfig(c)
	l.LogBanner()
	http.HandleFunc("/", l.Handler())
	fmt.Println(http.ListenAndServe(":8080", nil))
}
