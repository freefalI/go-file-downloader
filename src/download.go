package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

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
	fileName := strconv.Itoa(fileObject.ID) + "-" + path.Base(response.Request.URL.Path)
	fileObject.mu.Lock()
	fileObject.Name = fileName
	fileObject.mu.Unlock()
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
		fileObject.mu.Lock()
		fileObject.Size = int(totalBytes) / 1024 / 1024
		fileObject.mu.Unlock()
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

		fileObject.mu.Lock()
		if sizeKnown {
			fileObject.Progress = int(float64(downloadedBytes) / float64(totalBytes) * 100)
		} else {
			fileObject.Size = int(downloadedBytes) / 1024 / 1024 //to MB
		}
		fileObject.mu.Unlock()
	}
	fileObject.mu.Lock()
	fileObject.IsDone = true
	fileObject.FinishedAt = time.Now()
	fileObject.mu.Unlock()

	fmt.Printf("\nDownloaded file: %s\n", filePath)
}
