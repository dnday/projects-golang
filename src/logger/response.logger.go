package logger

import (
	"log"
	"net/http"
)

// ResponseLogger logs a message with the HTTP request method, request URI, status code, and the message that was passed in.
func ResponseLogger(r *http.Request, s uint, msg string) {
	log.Printf("[%s] %s %d - %s", r.Method, r.RequestURI, s, msg)
}
