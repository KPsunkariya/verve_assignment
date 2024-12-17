package main

import (
	"log"
	"net/http"
	"os"
	"verve_assignment/handlers"
	"verve_assignment/services"
)

func main() {
	// Logger setup
	logFile, err := setupLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func(logFile *os.File) {
		errLogFileClose := logFile.Close()
		if errLogFileClose != nil {
			log.Println(errLogFileClose)
		}
	}(logFile)

	initServices()
	go services.StartPeriodicLogger()

	// Start HTTP server
	http.HandleFunc("/api/verve/accept", handlers.AcceptHandler)
	http.HandleFunc("/api/verve/uniqueRequest", handlers.UniqueRequestHandler)

	log.Println("Starting server on :8989")
	err = http.ListenAndServe(":8989", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initServices() {
	services.InitKafkaProducer("localhost:9092")
	services.InitRedisClient("localhost:6379")
}

func ensureLogDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func setupLogger() (*os.File, error) {
	dir := "log"
	file := "app.log"
	err := ensureLogDir(dir)
	if err != nil {
		return nil, err
	}

	filePath := dir + "/" + file

	// Open or create the log file
	logFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	// Set default logger output to this file
	log.SetOutput(logFile)

	// Add optional log prefix and flags
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return logFile, nil
}
