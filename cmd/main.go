package main

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
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
	Size       int       `json:"size"`
	IsDone     bool      `json:"is_done"`
}

var (
	mu    sync.Mutex
	files = []*File{
		// {1, "http://", "nam", time.Now(), time.Now().Add(time.Minute), 20},
	}
	nextID = 1 // To generate unique IDs for each file
)

const downloadDir = "Downloads"

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

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(file)
}

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

func setupTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "URL", "Name", "Started at", "Finished at", "Size", "Progress", "Done"})
	table.SetBorder(true)

	return table
}

func printEntries(table *tablewriter.Table) {
	table.ClearRows()

	for _, file := range files {
		var finishedAt string
		if !file.FinishedAt.IsZero() {
			finishedAt = file.FinishedAt.Format(time.DateTime)
		}

		table.Append([]string{
			fmt.Sprintf("%d", file.ID),
			file.URL,
			file.Name,
			file.StartedAt.Format(time.DateTime),
			finishedAt,
			fmt.Sprintf("%dMb", file.Size),
			fmt.Sprintf("%d%%", file.Progress),
			fmt.Sprintf("%t", file.IsDone),
		})
	}

	clearTerminal()
	table.Render()
}

func clearTerminal() {
	os.Stdout.WriteString("\033[H\033[2J")
}

func createDownloadsDir() {
	// Create a "Downloads" directory in the current working directory
	err := os.MkdirAll(downloadDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating Downloads directory: %v\n", err)
		return
	}
}

func downloadFile(fileObject *File) {
	log.Printf("File with id:%d is processing", fileObject.ID)

	// Get the data
	response, err := http.Get(fileObject.URL)
	if err != nil {
		log.Printf("failed to download file: %s", err)
		return
	}
	defer response.Body.Close()

	// Check if the response status is OK
	if response.StatusCode != http.StatusOK {
		log.Printf("failed to download file: %s", response.Status)
		return
	}

	// Create the file
	fileName := path.Base(response.Request.URL.Path)
	fileObject.Name = fileName
	filePath := filepath.Join(downloadDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("failed to create file: %s", err)
		return
	}
	defer file.Close()

	// Get the content length
	totalBytes := response.ContentLength
	fmt.Println(totalBytes)
	var downloadedBytes int64

	sizeKnown := totalBytes > 0
	// Create a buffer to hold chunks of data
	buffer := make([]byte, 32*1024) // 32 KB buffer

	if sizeKnown {
		fileObject.Size = int(totalBytes) / 1024 / 1024
	}
	// Download the file
	for {
		// Read a chunk from the response body
		bytesRead, err := response.Body.Read(buffer)
		if err != nil && err != io.EOF {
			log.Printf("failed to read response: %s", err)
			return
		}
		if bytesRead == 0 {
			break // End of file
		}

		// Write the chunk to the file
		_, err = file.Write(buffer[:bytesRead])
		if err != nil {
			log.Printf("failed to write to file: %s", err)
			return
		}

		// Update the downloaded bytes count
		downloadedBytes += int64(bytesRead)
		if !sizeKnown {
			fileObject.Size = int(downloadedBytes) / 1024 / 1024
		}
		if sizeKnown {
			fileObject.Progress = int(float64(downloadedBytes) / float64(totalBytes) * 100)
		}
	}
	fileObject.IsDone = true
	fileObject.FinishedAt = time.Now()
	
	fmt.Printf("\nDownloaded file: %s\n", filePath)
}
