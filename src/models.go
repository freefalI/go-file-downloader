package main

import (
	"sync"
	"time"
)

type File struct {
	ID         int        `json:"id"`
	URL        string     `json:"url"`
	Name       string     `json:"name"`
	StartedAt  time.Time  `json:"started_at"`
	FinishedAt time.Time  `json:"finished_at"`
	Progress   int        `json:"progress"`
	Size       int        `json:"size"` //in MB
	IsDone     bool       `json:"is_done"`
	mu         sync.Mutex // Mutex for this specific file
}
