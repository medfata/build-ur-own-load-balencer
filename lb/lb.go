package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

func logRequest(r *http.Request) {
	log.Printf("Received request from %s\n", r.RemoteAddr)
	log.Printf("%s %s %s\n", r.Method, r.URL.Path, r.Proto)
	log.Printf("Host: %s\n", r.Host)
	log.Printf("User-Agent: %s\n", r.Header.Get("User-Agent"))
	log.Printf("Accept: %s\n", r.Header.Get("Accept"))
	log.Printf("<------------------------------------>")
}

type Backend struct {
	URL  string
	Down uint32
}

var (
	backends = []Backend{
		{URL: "http://backend1:5433", Down: 1},
		{URL: "http://backend2:5434", Down: 1},
	}
	counter uint32
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	// Round-robin load balancing
	backendUrl, err := getNextBackend()
	if err != nil {
		http.Error(w, "No available backends", http.StatusServiceUnavailable)
		return
	}
	// Construct a new request with the same method and URL as the incoming request
	requestURL := fmt.Sprintf("%s%s", backendUrl, r.URL.Path)
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

func getNextBackend() (string, error) {
	// Filter out backends that are down
	availableBackends := make([]Backend, 0, len(backends))
	for _, backend := range backends {
		if backend.Down == 0 {
			availableBackends = append(availableBackends, backend)
		}
	}

	// Check if there are any available backends
	if len(availableBackends) == 0 {
		return "", fmt.Errorf("no available backends")
	}

	// Select the next backend using the round-robin algorithm
	idx := atomic.AddUint32(&counter, 1)
	return availableBackends[idx%uint32(len(availableBackends))].URL, nil
}

func checkBackendsHealth(period int32) {
	for {
		time.Sleep(time.Duration(period) * time.Second)

		for i := range backends {
			go func(backend *Backend) {
				resp, err := http.Get(backend.URL + "/health-check")
				if err != nil || resp.StatusCode != http.StatusOK {
					if backend.Down == 1 {
						atomic.StoreUint32((*uint32)(&backend.Down), 1) // Mark as down
					}
				} else {
					if backend.Down == 1 {
						atomic.StoreUint32((*uint32)(&backend.Down), 0) // Mark as up
					}
				}
			}(&backends[i])
		}
	}
}

func main() {
	go checkBackendsHealth(10)
	http.HandleFunc("/", handleRequest)
	log.Println("Starting server on port 80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Could not start server: %s\n", err)
	}
}
