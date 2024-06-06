package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"sync"
)

func sendRequests(numRequests int, filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var wg sync.WaitGroup

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(requestNum int) {
			defer wg.Done()

			resp, err := http.Get("http://localhost")
			if err != nil {
				fmt.Fprintf(file, "Request %d: Error - %v\n", requestNum, err)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintf(file, "Request %d: Error reading response body - %v\n", requestNum, err)
				return
			}

			fmt.Fprintf(file, "Request %d: Status Code - %d, Body - %s\n", requestNum, resp.StatusCode, string(body))
		}(i + 1)
	}

	wg.Wait()
}

func main() {
	numCPU := runtime.NumCPU()
	fmt.Printf("Number of logical CPUs: %d\n", numCPU)
	numRequests := 500
	filename := "test-resp.txt"
	sendRequests(numRequests, filename)
	fmt.Printf("All %d requests have been sent. Responses and errors are logged in %s.\n", numRequests, filename)
}
