package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
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

var total int64 = 0

// Can be either "text" or "json"
var logFormat = "text"

type request struct {
	Body          []byte      `json:"body"`
	Fragment      string      `json:"fragment"`
	Headers       http.Header `json:"headers"`
	Method        string      `json:"method"`
	Path          string      `json:"path"`
	Query         string      `json:"query"`
	QueryParams   url.Values  `json:"queryParams"`
	ReceivedAt    int64       `json:"receivedAt"`
	TotalRequests int64       `json:"totalRequests"`
	Errors        []string    `json:"errors"`
}

func newRequest(r *http.Request) *request {
	req := request{}
	req.Method = r.Method
	req.Path = r.URL.Path
	req.Query = r.URL.RawQuery
	req.Fragment = r.URL.RawFragment
	req.TotalRequests = total
	req.ReceivedAt = time.Now().Unix()
	req.Headers = r.Header

	body, err := io.ReadAll(r.Body)
	if err != nil {
		req.Errors = append(req.Errors, err.Error())
	}
	req.Body = body

	queryParams, err := url.ParseQuery(req.Query)
	if err != nil {
		req.Errors = append(req.Errors, err.Error())
	}
	req.QueryParams = queryParams

	return &req
}

func toJSON(r *http.Request) ([]byte, error) {
	req := newRequest(r)
	return json.Marshal(req)
}

func prettyPrint(r *http.Request) (out string) {
	query := r.URL.RawQuery
	if query != "" {
		query = "?" + query
	}
	fragment := r.URL.RawFragment
	if fragment != "" {
		fragment = "#" + fragment
	}
	out += fmt.Sprintf("Total requests: %d\n", total)
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
	total++
	switch logFormat {
	case "json":
		printJSON(r)
	default:
		printText(r)
	}
}

func printText(r *http.Request) {
	fmt.Println(sep)
	fmt.Println(prettyPrint(r))
	fmt.Println(sep)
}

func printJSON(r *http.Request) {
	b, _ := toJSON(r)
	fmt.Println(string(b))
}

func printBanner() {
	// Should not be printed in JSON mode
	if logFormat == "json" {
		return
	}
	fmt.Println(banner)
	fmt.Println(sep)
}

func main() {
	if lf := os.Getenv("FLIES_LOG_FORMAT"); lf != "" {
		logFormat = lf
	}
	printBanner()
	http.HandleFunc("/", handler)
	fmt.Println(http.ListenAndServe(":8080", nil))
}
