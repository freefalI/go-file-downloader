package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	mu    sync.Mutex // Controls read and writes to files
	files = []*File{
		// {1, "http://", "nam", time.Now(), time.Now().Add(time.Minute), 20},
	}
	nextID = 1 // To generate unique IDs for each file
)

const downloadDir = "Downloads"

func main() {
	http.HandleFunc("/files", getFiles)
	http.HandleFunc("/files/add", addFile)

	log.Println("Server is running on port :8080")

	createDownloadsDir()

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	table := setupTable()

	// Start a goroutine to print progress
	go func() {
		for {
			select {
			case <-ticker.C:
				printEntries(table)
			}
		}
	}()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
