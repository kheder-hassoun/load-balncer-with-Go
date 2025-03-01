package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ServerStats keeps track of how many times each server replied
type ServerStats struct {
	mu    sync.Mutex
	stats map[string]int
}

// NewServerStats initializes a new ServerStats instance
func NewServerStats() *ServerStats {
	return &ServerStats{
		stats: make(map[string]int),
	}
}

// Increment increments the count for a specific server
func (s *ServerStats) Increment(server string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats[server]++
}

// Print prints the statistics
func (s *ServerStats) Print() {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Println("Server Statistics:")
	for server, count := range s.stats {
		fmt.Printf("%s: %d requests\n", server, count)
	}
}

// sendRequest sends a request to the load balancer and extracts the server name from the response
func sendRequest(url string, stats *ServerStats, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error sending request: %v\n", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return
	}

	// Extract the server name from the HTML response
	html := string(body)
	serverName := extractServerName(html)
	if serverName == "" {
		log.Println("Could not extract server name from response")
		return
	}

	// Update statistics
	stats.Increment(serverName)
	fmt.Printf("Request replied by: %s\n", serverName)
}

// extractServerName extracts the server name from the HTML response
func extractServerName(html string) string {
	// Look for the server name in the HTML
	startTag := `<h1 style="color:#00b4ff; text-align: center; font-style: italic" id="server_name">`
	endTag := `</h1>`

	start := strings.Index(html, startTag)
	if start == -1 {
		return ""
	}
	start += len(startTag)

	end := strings.Index(html[start:], endTag)
	if end == -1 {
		return ""
	}

	return html[start : start+end]
}

func main() {
	// URL of the load balancer
	loadBalancerURL := "http://localhost:9095"

	// Number of requests to send
	numRequests := 10

	// Create a wait group to wait for all requests to complete
	var wg sync.WaitGroup

	// Initialize server statistics
	stats := NewServerStats()

	// Send requests in separate goroutines
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go sendRequest(loadBalancerURL, stats, &wg)
		time.Sleep(100 * time.Millisecond) // Add a small delay between requests
	}

	// Wait for all requests to complete
	wg.Wait()

	// Print statistics
	stats.Print()
}
