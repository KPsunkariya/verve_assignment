package handlers

import (
	"log"
	"net/http"
	"verve_assignment/services"
)

func UniqueRequestHandler(w http.ResponseWriter, r *http.Request) {
	writeOk(w, "ok")
	return
}

func AcceptHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	endpoint := r.URL.Query().Get("endpoint")

	// Validate the ID
	if id == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	log.Println("input id: ", id)

	if !services.IsUniqueID(id) {
		log.Printf("Duplicate ID: %s", id)
		writeOk(w, "ok")
		return
	}

	//// Fire HTTP request if endpoint exists
	//if endpoint != "" {
	//	go services.FireHTTPRequest(endpoint)
	//}

	// If endpoint is provided, fire a POST request
	if endpoint != "" {
		uniqueCount := services.GetUniqueCount() // Retrieve current unique count
		go services.FirePOSTRequest(endpoint, uniqueCount)
	}
	writeOk(w, "ok")
}

func writeOk(w http.ResponseWriter, status string) {
	w.WriteHeader(http.StatusOK)
	write, err := w.Write([]byte(status))
	if err != nil {
		log.Println("error while write response", err)
		return
	}

	log.Println("write:", write)
}
