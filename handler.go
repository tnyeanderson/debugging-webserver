package main

import "net/http"

// Handler returns an http.Handler which calls WriteRequest on the Logger.
func Handler(l Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.WriteRequest(r)
	}
}
