package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	Wire          []byte      `json:"wire"`
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

	// Read the body into both Wire and Body
	wireBuffer := &bytes.Buffer{}
	bodyBuffer := &bytes.Buffer{}
	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, bodyBuffer))
	if err := r.Write(wireBuffer); err != nil {
		req.Errors = append(req.Errors, err.Error())
	}
	req.Body = bodyBuffer.Bytes()
	req.Wire = wireBuffer.Bytes()

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
	logSeparator()
	r.Write(os.Stdout)
	fmt.Println()
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
	fmt.Print(banner)
	fmt.Println()
	printSeparatorFullLine("+")
}

func logSeparator() {
	t := time.Now().Format(time.UnixDate)
	c := fmt.Sprintf("Total requests: %d", total)
	printSeparatorFullLine("*")
	printSeparatorMessages("-", []string{t, c})
}

func printSeparatorFullLine(char string) {
	printSeparatorMessages(char, []string{""})
}

func printSeparatorMessages(sep string, messages []string) {
	width := 80
	prefix := 3
	for _, m := range messages {
		if len(m) == 0 {
			fmt.Println(strings.Repeat(sep, width))
			continue
		}
		// prefix, space, message, space, suffix
		suffixLength := width - 3 - 1 - len(m) - 1
		fmt.Printf("%s %s %s\n", strings.Repeat(sep, prefix), m, strings.Repeat(sep, suffixLength))
	}
}

func main() {
	if lf := os.Getenv("FLIES_LOG_FORMAT"); lf != "" {
		logFormat = lf
	}
	printBanner()
	http.HandleFunc("/", handler)
	fmt.Println(http.ListenAndServe(":8080", nil))
}
