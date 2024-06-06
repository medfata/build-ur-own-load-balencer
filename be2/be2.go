package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request from %s\n", r.RemoteAddr)
	log.Printf("%s %s %s\n", r.Method, r.URL.Path, r.Proto)
	log.Printf("Host: %s\n", r.Host)
	log.Printf("User-Agent: %s\n", r.Header.Get("User-Agent"))
	log.Printf("Accept: %s\n", r.Header.Get("Accept"))
	htmlContent := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="utf-8">
		<title>Index Page</title>
	</head>
	<body>
		Hello from the web server running on port 5434.
	</body>
	</html>`
	// Set the Content-Type header to indicate that the response contains HTML
	w.Header().Set("Content-Type", "text/html")
	// Write the HTML content to the response writer
	fmt.Fprintf(w, htmlContent)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/health-check", healthCheckHandler)
	log.Println("Starting backeend server on port 5434")
	err := http.ListenAndServe(":5434", nil)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Could not start backend server: %s\n", err)
	}
}
