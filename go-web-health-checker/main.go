package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HealthCheckResult struct {
	URL    string `json: "url"`
	Status string `json: "status"`
}

func getHealthInfo(url string, ch chan<- HealthCheckResult) error {
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get(url)
	status := "Down"
	if err != nil {
		ch <- HealthCheckResult{URL: url, Status: status}
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		status = "UP"
	}

	ch <- HealthCheckResult{URL: url, Status: status}
	return nil
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	var urls []string
	if err := json.NewDecoder(r.Body).Decode(&urls); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ch := make(chan HealthCheckResult, len(urls))
	for _, url := range urls {
		go getHealthInfo(url, ch)
	}

	results := make([]HealthCheckResult, 0, len(urls))
	for range urls {
		results = append(results, <-ch)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	fmt.Println("starting server")
	http.HandleFunc("/health-check", healthCheckHandler)
	fmt.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
