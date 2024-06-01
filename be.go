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

    fmt.Fprintf(w,"Hello From Backend Server \n")
}

func main() {
    http.HandleFunc("/", handler)
    log.Println("Starting backeend server on port 5433")
    err := http.ListenAndServe(":5433", nil)
    if err != nil {
		fmt.Println(err)
        log.Fatalf("Could not start backend server: %s\n", err)
    }
}
