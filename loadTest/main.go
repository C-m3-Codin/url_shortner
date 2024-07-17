package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	"time"
)

const (
	tokenURL    = "http://localhost:8000/token"
	loadTestURL = "http://localhost:8000/sec/create"
	tokenData   = `{"name":"Cyril","username":"c1337_2","email":"c@c_2.com","password":"password@123"}`
	jsonData    = `{"original_url": "https://www.google.com"}`
	numRequests = 100000 // Number of requests you want to send
	numWorkers  = 10     // Number of concurrent goroutines
)

var (
	mu            sync.Mutex
	errorCount    int
	responseTimes []time.Duration
	token         string
)

func main() {
	// Get token
	var err error
	token, err = getToken()
	if err != nil {
		fmt.Println("Failed to get token:", err)
		return
	}

	var wg sync.WaitGroup
	jobs := make(chan struct{}, numRequests)

	// Create worker goroutines
	for i := 0; i < numWorkers; i++ {
		go func() {
			for range jobs {
				sendRequest()
				wg.Done()
			}
		}()
	}

	// Send requests
	for i := 0; i < numRequests; i++ {
		fmt.Println("job , ", i)
		wg.Add(1)
		jobs <- struct{}{}
	}

	close(jobs)
	wg.Wait()

	// Calculate and print statistics
	calculateAndPrintStatistics()
}

func getToken() (string, error) {
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBuffer([]byte(tokenData)))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Assuming the token is in the response body as a JSON field "token"
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error unmarshalling response: %w", err)
	}

	token, ok := result["token"].(string)
	if !ok {
		return "", fmt.Errorf("token not found in response")
	}

	return token, nil
}

func sendRequest() {
	start := time.Now()

	req, err := http.NewRequest("POST", loadTestURL, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		incrementErrorCount()
		return
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		incrementErrorCount()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		incrementErrorCount()
	}

	elapsed := time.Since(start)
	recordResponseTime(elapsed)
}

func incrementErrorCount() {
	mu.Lock()
	defer mu.Unlock()
	errorCount++
}

func recordResponseTime(d time.Duration) {
	mu.Lock()
	defer mu.Unlock()
	responseTimes = append(responseTimes, d)
}

func calculateAndPrintStatistics() {
	mu.Lock()
	defer mu.Unlock()

	fmt.Printf("Total errors: %d\n", errorCount)
	fmt.Printf("Total responses: %d\n", len(responseTimes))

	if len(responseTimes) == 0 {
		return
	}

	// Sort response times for percentile calculations
	sort.Slice(responseTimes, func(i, j int) bool {
		return responseTimes[i] < responseTimes[j]
	})

	avg := calculateAverage(responseTimes)
	p50 := calculatePercentile(responseTimes, 50)
	p80 := calculatePercentile(responseTimes, 80)
	p95 := calculatePercentile(responseTimes, 95)

	fmt.Printf("Average response time: %s\n", avg)
	fmt.Printf("50th percentile (median) response time: %s\n", p50)
	fmt.Printf("80th percentile response time: %s\n", p80)
	fmt.Printf("95th percentile response time: %s\n", p95)
}

func calculateAverage(durations []time.Duration) time.Duration {
	var total time.Duration
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}

func calculatePercentile(durations []time.Duration, percentile int) time.Duration {
	index := (percentile * len(durations)) / 100
	return durations[index]
}
