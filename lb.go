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

    fmt.Fprintf(w, "Hello, you've hit %s\n", r.URL.Path)
}

func main() {
    http.HandleFunc("/", handler)
    log.Println("Starting server on port 5432")
    err := http.ListenAndServe(":5432", nil)
    if err != nil {
		fmt.Println(err)
        log.Fatalf("Could not start server: %s\n", err)
    }
}
