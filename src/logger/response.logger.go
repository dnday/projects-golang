package logger

import (
	"log"
	"net/http"
)

func ResponseLogger(r *http.Request, s uint, msg string) {
	log.Printf("[%s] %s %d - %s", r.Method, r.RequestURI, s, msg)
}
