package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// LoggingMiddleware logs HTTP requests including method, path, and POST body
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log basic request info
		slog.Info(fmt.Sprintf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr))

		// Log POST body if present
		if r.Method == "POST" && r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				slog.Info(fmt.Sprintf("Error reading body: %v", err))
			} else {
				// Log the body
				slog.Info(fmt.Sprintf("Body: %s", string(bodyBytes)))
				// Restore the body for the next handler
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		next.ServeHTTP(w, r)

		slog.Info(fmt.Sprintf("Request completed in %v", time.Since(start)))
	})
}
