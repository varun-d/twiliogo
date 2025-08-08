package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

const welcomemsg string = "Welcome to a go test server. /hello and /events is active. /events is testing SSE events."

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	var addr string
	flag.StringVar(&addr, "p", ":8000", "port to connect the server to")
	flag.Parse()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/hello", hellohandle)
	http.HandleFunc("/events", rtevents)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func rtevents(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request from %s for %s", r.RemoteAddr, r.URL.Path)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	tokens := []string{"this", "is", "an", "event", "stream", "bringing", "in", "live", "data", "every", "400", "ms"}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	for _, tok := range tokens {
		// Format the event
		fmt.Fprintf(w, "data: %s\n\n", tok)
		// Flush instantly to send data
		flusher.Flush()
		// Simulate delay!
		time.Sleep(time.Millisecond * 400)
	}
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request from %s for %s", r.RemoteAddr, r.URL.Path)
	w.Header().Set("X-Custom-Header", "Custom Header Value")
	fmt.Fprintf(w, welcomemsg)
}

func hellohandle(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request from %s for %s", r.RemoteAddr, r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r.Header)
	// fmt.Fprintln(w, r.Header)
}
