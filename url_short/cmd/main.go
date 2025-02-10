package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/its-kos/assignment1/api"
	"github.com/its-kos/assignment1/storage"
	"github.com/its-kos/assignment1/types"
)

const (
	dbFilePath      = "db.json"
	cleanUpInterval = 30 		// Seconds
)

func main() {

	var dbLock sync.Mutex

	go func() {
		for {
			time.Sleep(cleanUpInterval * time.Second) // Run every X seconds
			removeExpiredURLs(&dbLock)                // Call cleanup function
		}
	}()

	server := api.NewApiServer(":8000", storage.NewMemoryStorage(&dbLock))
	log.Fatal(server.Run())
}

// Reads, filters, and writes back valid URLs to file
func removeExpiredURLs(mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()

	// Check - if file does not exist we skip
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		slog.Info("File does not exist, skipping...")
		return
	}

	// Open db file
	file, err := os.OpenFile(dbFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		slog.Error(fmt.Sprint("Error opening file:", err))
		return
	}

	// Read file in array urls
	var urls []types.URL
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&urls); err != nil && err.Error() != "EOF" {
		slog.Error(fmt.Sprint("Error decoding json:", err))
		file.Close()
		return
	}
	file.Close()

	// Filter out expired URLs
	now := time.Now()
	var validURLs []types.URL
	var expiredURLs []types.URL
	for _, url := range urls {
		expiry := url.CreatedAt.Add(time.Duration(url.TimeToLive) * time.Second)
		if now.Before(expiry) {
			validURLs = append(validURLs, url)
		} else {
			expiredURLs = append(expiredURLs, url)
		}
	}

	// Open the file again to write valid urls with O_TRUNC flag to clear file contents
	file, err = os.OpenFile(dbFilePath, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		slog.Error(fmt.Sprint("Error opening file:", err))
		return
	}
	defer file.Close()

	// Dump valid urls back in file
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(validURLs); err != nil {
		slog.Error(fmt.Sprint("Error writing updated json", err))
	}

	slog.Info(fmt.Sprint(len(expiredURLs), " url entries removed"), slog.Any("expired urls", expiredURLs))
}
