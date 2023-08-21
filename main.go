package main

import (
	"fmt"
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

func main() {
	// Initialize server
	s := &defaultServer{}
	s.Init()
	s.Logger = s.GetLogger()
	s.Logger.Init()
	fmt.Println(s.Listen())
}
