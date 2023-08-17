package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const banner = `
   __ _ _           
  / _| (_)          
 | |_| |_  ___  ___ 
 |  _| | |/ _ \/ __|
 | | | | |  __/\__ \
 |_| |_|_|\___||___/
`

const sep = "-----------------------------"

type request struct {
	*http.Request
}

func (r *request) String() (out string) {
	query := r.URL.RawQuery
	if query != "" {
		query = "?" + query
	}
	fragment := r.URL.RawFragment
	if fragment != "" {
		fragment = "#" + fragment
	}
	out += fmt.Sprintf("Received: %s\n\n", time.Now().Format(time.UnixDate))
	out += fmt.Sprintf("%s %s%s#%s\n", r.Method, r.URL.Path, query, fragment)
	for key, values := range r.Header {
		for _, value := range values {
			out += fmt.Sprintf("%s: %s\n", key, value)
		}
	}
	out += fmt.Sprintf("\n")
	body, _ := io.ReadAll(r.Body)
	out += string(body)
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	req := request{r}
	fmt.Println(sep)
	fmt.Println(req.String())
	fmt.Println(sep)
}

func main() {
	http.HandleFunc("/", handler)

	log.Println(banner)
	fmt.Println(sep)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
