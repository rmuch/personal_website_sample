package main

import (
	"log"
	"net/http"
)

// A http.Handler object to log requests before pushing them onwards.
type LogFilter struct {
	inner http.Handler
}

func (l *LogFilter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if .inner is nil.
	if l.inner == nil {
		log.Fatal("LogFilter proxy object (.inner) undefined.")
	}

	// Log the request.
	// TODO: Implement SQL logging, option to disable logging.

	if r.Header["X-Forwarded-For"] != nil {
		log.Printf("Serving a request %s via %s for %s", r.RequestURI, r.RemoteAddr, r.Header["X-Forwarded-For"])
	} else {
		log.Printf("Serving a request %s for %s", r.RequestURI, r.RemoteAddr)
	}

	// Forward the request.
	l.inner.ServeHTTP(w, r)
}

// Create a new log filter.
func NewLogFilter(inner http.Handler) *LogFilter {
	return &LogFilter{inner}
}
