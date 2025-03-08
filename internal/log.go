package internal

import (
	"log"
	"net/http"
	"os"
)

var hostname string

func init() {
	hostname, _ = os.Hostname()
}

func LogRequest(r *http.Request, message string, code int) {
	log.Printf("IP: %s, Method: %s, Path: %s, Status: %d, Size: %d, Message: %s",
		hostname, r.Method, r.URL.Path, code, r.ContentLength, message)
}
