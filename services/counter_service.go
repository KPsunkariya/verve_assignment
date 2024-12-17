package services

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	mu             sync.Mutex
	uniqueRequests = make(map[string]struct{})
	kafkaTopic     = "unique-request-counts"
)

func GetUniqueCount() int {
	return len(uniqueRequests)
}

//func IsUniqueID(id string) bool {
//	mu.Lock()         // Lock the map
//	defer mu.Unlock() // Ensure unlock happens after this function
//
//	// Check if the ID already exists
//	if _, exist := uniqueRequests[id]; exist {
//		return false // ID is not unique
//	}
//	// Add the ID to the map
//	uniqueRequests[id] = struct{}{}
//	return true // ID is unique
//}

func StartPeriodicLogger() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		// Lock the map to count unique IDs
		mu.Lock()
		count := len(uniqueRequests)
		// Log the count
		log.Printf("Unique requests in the last minute: %d", count)
		message := `{"timestamp": "` + time.Now().Format(time.RFC3339) + `", "unique_request_count": ` + strconv.Itoa(count) + `}`
		// Send message to Kafka
		SendToKafka(kafkaTopic, message)

		// Clear the map for the next minute
		uniqueRequests = make(map[string]struct{})
		mu.Unlock()
	}
}

func FireHTTPRequest(endpoint string) {
	resp, err := http.Get(endpoint)
	if err != nil {
		log.Printf("Failed to reach endpoint: %v", err)
		return
	}
	log.Printf("HTTP GET to %s returned status: %d", endpoint, resp.StatusCode)
}
