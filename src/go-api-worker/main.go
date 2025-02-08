// worker.go
package main

import (
	"net/http"
	"sync"
	"time"
)

type Worker struct {
    client  *http.Client
    baseURL string
}

func NewWorker(baseURL string) *Worker {
    return &Worker{
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
        baseURL: baseURL,
    }
}

func (w *Worker) makeRequests(numRequests int) {
    var wg sync.WaitGroup
    semaphore := make(chan struct{}, 100) // Limit concurrent requests

    for i := 0; i < numRequests; i++ {
        wg.Add(1)
        semaphore <- struct{}{} // Acquire semaphore
        
        go func() {
            defer wg.Done()
            defer func() { <-semaphore }() // Release semaphore
            
            _, err := w.client.Get(w.baseURL + "/api")
            if err != nil {
                return
            }
        }()
    }
    
    wg.Wait()
}

func main() {
    worker := NewWorker("http://localhost:8080")
    worker.makeRequests(1000)
}
