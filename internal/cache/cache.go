package cache

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Cache represents a simple file-based cache with an in-memory layer
type Cache struct {
	dir           string
	inMemoryCache map[string]CacheEntry
	mu            sync.RWMutex // Mutex for protecting inMemoryCache
	expiration    time.Duration
}

// CacheEntry represents a cached endpoint summary
type CacheEntry struct {
	Summary   string    `json:"summary"`
	Timestamp time.Time `json:"timestamp"`
}

// NewCache creates a new cache instance
func NewCache(expiration time.Duration) (*Cache, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	cacheDir := filepath.Join(home, ".restapisummarizer", "cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory %s: %w", cacheDir, err)
	}

	return &Cache{
		dir:           cacheDir,
		inMemoryCache: make(map[string]CacheEntry),
		expiration:    expiration,
	}, nil
}

// Get retrieves a cached summary for the given endpoint
func (c *Cache) Get(method, path, fileHash string) (string, bool) {
	key := generateKey(method, path, fileHash)

	// Try to get from in-memory cache first
	c.mu.RLock()
	entry, ok := c.inMemoryCache[key]
	c.mu.RUnlock()

	if ok {
		if time.Since(entry.Timestamp) <= c.expiration {
			return entry.Summary, true
		}
		// Entry expired in memory, remove it
		c.mu.Lock()
		delete(c.inMemoryCache, key)
		c.mu.Unlock()
	}

	// If not in memory or expired, try to read from disk
	cachePath := filepath.Join(c.dir, key+".json")
	data, err := os.ReadFile(cachePath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error reading cache file %s: %v\n", cachePath, err)
		}
		return "", false
	}

	if err := json.Unmarshal(data, &entry); err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling cache entry from %s: %v\n", cachePath, err)
		os.Remove(cachePath) // Remove corrupted file
		return "", false
	}

	// Check if disk cache is still valid
	if time.Since(entry.Timestamp) > c.expiration {
		os.Remove(cachePath) // Remove expired file
		return "", false
	}

	// Add to in-memory cache
	c.mu.Lock()
	c.inMemoryCache[key] = entry
	c.mu.Unlock()

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
		return fmt.Errorf("failed to marshal cache entry: %w", err)
	}

	// Write to disk
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file %s: %w", cachePath, err)
	}

	// Add/update in-memory cache
	c.mu.Lock()
	c.inMemoryCache[key] = entry
	c.mu.Unlock()

	return nil
}

// Clear removes all cached entries
func (c *Cache) Clear() error {
	// Clear in-memory cache
	c.mu.Lock()
	c.inMemoryCache = make(map[string]CacheEntry)
	c.mu.Unlock()

	// Remove all files from disk
	if err := os.RemoveAll(c.dir); err != nil {
		return fmt.Errorf("failed to clear cache directory %s: %w", c.dir, err)
	}
	// Recreate the directory as RemoveAll also deletes the directory itself
	if err := os.MkdirAll(c.dir, 0755); err != nil {
		return fmt.Errorf("failed to recreate cache directory %s after clearing: %w", c.dir, err)
	}
	return nil
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