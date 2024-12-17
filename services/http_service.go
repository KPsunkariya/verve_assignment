package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Define the payload structure
type PostPayload struct {
	Timestamp          string `json:"timestamp"`
	UniqueRequestCount int    `json:"unique_request_count"`
}

// Function to send an HTTP POST request
func FirePOSTRequest(endpoint string, uniqueCount int) {
	// Create the payload
	payload := PostPayload{
		Timestamp:          time.Now().Format(time.RFC3339),
		UniqueRequestCount: uniqueCount,
	}

	// Serialize the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON payload: %v", err)
		return
	}

	// Send the POST request
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to send POST request: %v", err)
		return
	}
	defer resp.Body.Close()

	// Log the HTTP status code
	log.Printf("POST to %s returned status: %d", endpoint, resp.StatusCode)
}
