package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs HTTP requests including method, path, and POST body
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log basic request info
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Log POST body if present
		if r.Method == "POST" && r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading body: %v", err)
			} else {
				// Log the body
				log.Printf("Body: %s", string(bodyBytes))
				// Restore the body for the next handler
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		next.ServeHTTP(w, r)

		log.Printf("Request completed in %v", time.Since(start))
	})
}
