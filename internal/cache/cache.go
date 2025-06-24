package cache

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// Cache represents a simple file-based cache
type Cache struct {
	dir string
}

// CacheEntry represents a cached endpoint summary
type CacheEntry struct {
	Summary   string    `json:"summary"`
	Timestamp time.Time `json:"timestamp"`
}

// NewCache creates a new cache instance
func NewCache() (*Cache, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	
	cacheDir := filepath.Join(home, ".restapisummarizer", "cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, err
	}
	
	return &Cache{dir: cacheDir}, nil
}

// Get retrieves a cached summary for the given endpoint
func (c *Cache) Get(method, path, fileHash string) (string, bool) {
	key := generateKey(method, path, fileHash)
	cachePath := filepath.Join(c.dir, key+".json")
	
	data, err := ioutil.ReadFile(cachePath)
	if err != nil {
		return "", false
	}
	
	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return "", false
	}
	
	// Check if cache is still valid (24 hours)
	if time.Since(entry.Timestamp) > 24*time.Hour {
		os.Remove(cachePath)
		return "", false
	}
	
	return entry.Summary, true
}

// Set stores a summary in the cache
func (c *Cache) Set(method, path, fileHash, summary string) error {
	key := generateKey(method, path, fileHash)
	cachePath := filepath.Join(c.dir, key+".json")
	
	entry := CacheEntry{
		Summary:   summary,
		Timestamp: time.Now(),
	}
	
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	
	return ioutil.WriteFile(cachePath, data, 0644)
}

// Clear removes all cached entries
func (c *Cache) Clear() error {
	return os.RemoveAll(c.dir)
}

// generateKey creates a unique cache key
func generateKey(method, path, fileHash string) string {
	data := fmt.Sprintf("%s:%s:%s", method, path, fileHash)
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// HashFile generates a hash of file content
func HashFile(content string) string {
	hash := md5.Sum([]byte(content))
	return fmt.Sprintf("%x", hash)
} 