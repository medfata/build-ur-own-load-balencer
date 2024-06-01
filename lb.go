package main

import (
    "fmt"
    "log"
    "net/http"
	"io/ioutil"
)
func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request from %s\n", r.RemoteAddr)
    log.Printf("%s %s %s\n", r.Method, r.URL.Path, r.Proto)
    log.Printf("Host: %s\n", r.Host)
    log.Printf("User-Agent: %s\n", r.Header.Get("User-Agent"))
    log.Printf("Accept: %s\n", r.Header.Get("Accept"))

	// Read the incoming request body
	//bodyBytes, err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	http.Error(w, "Internal server error", http.StatusInternalServerError)
	//	return
	//}

	// Construct a new request with the same method and URL as the incoming request
	requestURL := fmt.Sprintf("http://%s:%d%s", "localhost", 5433, "/")
	newReq, err := http.NewRequest(r.Method, requestURL, nil)
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

	// Send the new request to the destination server with port 5433
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
		log.Printf("Error while copyin the response from the destination server: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(respBytes)
}

func main() {
    http.HandleFunc("/", handleRequest)
    log.Println("Starting server on port 5432")
    err := http.ListenAndServe(":5432", nil)
    if err != nil {
		fmt.Println(err)
        log.Fatalf("Could not start server: %s\n", err)
    }
}
