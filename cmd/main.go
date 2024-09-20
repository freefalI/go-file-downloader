package main

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type File struct {
	ID         int       `json:"id"`
	URL        string    `json:"url"`
	Name       string    `json:"name"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	Progress   int       `json:"progress"`
}

var (
	mu    sync.Mutex
	files = []File{
		{1, "http://", "nam", time.Now(), time.Now().Add(time.Minute), 20},
	}
	nextID = 2 // To generate unique IDs for each file
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
	file := File{
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

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(file)
}

func main() {
	http.HandleFunc("/files", getFiles)
	http.HandleFunc("/files/add", addFile)

	log.Println("Server is running on port :8080")

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

func setupTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "URL", "Name", "Started at", "Finished at", "Progress"})
	table.SetBorder(true)

	return table
}

func printEntries(table *tablewriter.Table) {
	table.ClearRows()

	for _, file := range files {
		table.Append([]string{
			fmt.Sprintf("%d", file.ID),
			file.URL,
			file.Name,
			file.StartedAt.Format(time.RFC1123),
			file.FinishedAt.Format(time.RFC1123),
			fmt.Sprintf("%d%%", file.Progress),
		})
	}

	clearTerminal()
	table.Render()
}

func clearTerminal() {
	os.Stdout.WriteString("\033[H\033[2J")
}
