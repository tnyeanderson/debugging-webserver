package main

import (
	"fmt"
	"net/http"
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
	// Initialize config
	c := &config{}
	c.Init()

	// Initialize logger
	l := c.GetLogger()
	l.Init()

	// Start server
	http.HandleFunc("/", Handler(l))
	fmt.Println(http.ListenAndServe(c.GetAddr(), nil))
}
