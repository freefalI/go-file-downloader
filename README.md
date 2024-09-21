# Go File Downloader

## Overview

This project is a simple file downloader implemented in Go. It allows users to add files for download via an HTTP API and tracks the progress of each download.

## Features

- **HTTP API**: 
  - Add new files to download via a POST request.
  - Retrieve the list of all files and their download statuses.
  
- **Progress Tracking**: 
  - Real-time tracking of download progress, including percentage completion and estimated time remaining.

- **Concurrency**: 
  - Utilizes goroutines to handle multiple downloads simultaneously without blocking the main application.

- **Graceful Handling of Missing Content-Length**: 
  - Supports downloading files even when the `Content-Length` header is not provided by the server.

## Topics Practiced

During the development of this project, several key topics in Go were practiced:

1. **Goroutines and Concurrency**:
   - Implementing concurrent downloads using goroutines.
   - Managing shared state safely with mutexes.

2. **HTTP Server Development**:
   - Creating an HTTP server that handles requests and responses.
   - Implementing RESTful API endpoints for adding and retrieving files.

3. **JSON Encoding/Decoding**:
   - Working with JSON data structures for API communication.

4. **File I/O Operations**:
   - Handling file creation and writing downloaded content to disk.
   - Managing file paths using the `path` and `filepath` packages.

5. **Error Handling**:
   - Implementing robust error handling for network requests and file operations.

6. **Table Rendering in Console**:
   - Using external libraries (e.g., `tablewriter`) to format and display information in the console.

7. **Time Management**:
   - Utilizing the `time` package for tracking download start and finish times.

## Getting Started

### Prerequisites

- Go (version 1.23.1 or later)
- Internet connection for downloading files

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/freefalI/go-file-downloader.git
   cd go-file-downloader
   ```

2. Build the application:

   ```bash
   go build cmd/main.go
   ```

3. Run the application:

   ```bash
   ./go-file-downloader
   ```

### Usage

- To add a new file for download, send a POST request to `/files/add` with a JSON body containing the URL:

    ```json
    {
        "url": "https://ash-speed.hetzner.com/100MB.bin"
    }
    ```

- To retrieve the status of all downloads, send a GET request to `/files`.

### Example Requests

Using `curl`, you can interact with the API as follows:

- Add a new file:

    ```bash
    curl -X POST http://localhost:8080/files/add -d '{"url": "https://ash-speed.hetzner.com/100MB.bin"}' -H "Content-Type: application/json"
    ```

- Get all files:

    ```bash
    curl http://localhost:8080/files
    ```

### [Postman collection](file-downloader.postman_collection.json)

## License

This project is licensed under the MIT License.
