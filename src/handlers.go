package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Handler to get all files
func getFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mu.Lock()
	defer mu.Unlock()

	json.NewEncoder(w).Encode(files)
}

// Handler to add a new file entry based on URL
func addFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newFile struct {
		URL string `json:"url"`
	}

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&newFile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()

	// Create a new File entry
	file := &File{
		ID:        nextID,
		URL:       newFile.URL,
		Name:      "",
		StartedAt: time.Now(),
		Progress:  0,
	}

	// Increment the next ID for future entries
	nextID++

	// Append the new file to the slice
	files = append(files, file)

	mu.Unlock()

	log.Println("New file added")

	// Process the job immediately in a goroutine
	go downloadFile(file)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(file)
}
