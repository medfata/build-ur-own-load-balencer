package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
)

func logRequest(r *http.Request) {
	log.Printf("Received request from %s\n", r.RemoteAddr)
	log.Printf("%s %s %s\n", r.Method, r.URL.Path, r.Proto)
	log.Printf("Host: %s\n", r.Host)
	log.Printf("User-Agent: %s\n", r.Header.Get("User-Agent"))
	log.Printf("Accept: %s\n", r.Header.Get("Accept"))
	log.Printf("<------------------------------------>")
}

var (
	backends = []string{
		"http://backend1:5433",
		"http://backend2:5434",
	}
	counter uint32
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	// Round-robin load balancing
	backend := getNextBackend()

	// Construct a new request with the same method and URL as the incoming request
	requestURL := fmt.Sprintf("%s%s", backend, r.URL.Path)
	newReq, err := http.NewRequest(r.Method, requestURL, r.Body)
	if err != nil {
		log.Printf("Error while constructing new request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Copy the headers from the incoming request to the new request
	newReq.Header = make(http.Header)
	for k, v := range r.Header {
		newReq.Header[k] = v
	}

	// Send the new request to the destination server
	resp, err := http.DefaultClient.Do(newReq)
	if err != nil {
		log.Printf("Error while sending the request to the destination: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response from the destination server to the outgoing response
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error while copying the response from the destination server: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(respBytes)
}

func getNextBackend() string {
	idx := atomic.AddUint32(&counter, 1)
	return backends[idx%uint32(len(backends))]
}

func main() {
	http.HandleFunc("/", handleRequest)
	log.Println("Starting server on port 80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Could not start server: %s\n", err)
	}
}
