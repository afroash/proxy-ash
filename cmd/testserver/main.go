package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		
		// Add response headers to help verify proxy behavior
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Server", "TestServer")
		
		// Send a response with timestamp to help verify latency
		response := fmt.Sprintf(`{
			"message": "Hello from test server!",
			"path": "%s",
			"timestamp": "%s"
		}`, r.URL.Path, time.Now().Format(time.RFC3339Nano))
		
		w.Write([]byte(response))
	})

	log.Printf("Starting test server on :9090")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
