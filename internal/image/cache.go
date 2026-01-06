package image

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	cacheDirPerm  = 0755
	cacheFilePerm = 0644
	cacheVersion  = 1
)

// ImageCache manages cached image metadata
type ImageCache struct {
	dir string
}

// cacheEntry represents a cached image info entry
type cacheEntry struct {
	Version   int                 `json:"version"`
	Timestamp time.Time           `json:"timestamp"`
	ObjectName string             `json:"object_name"`
	ImageInfo *WikipediaImageInfo `json:"image_info"`
}

// NewImageCache creates a new image cache
func NewImageCache(dir string) (*ImageCache, error) {
	if dir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		dir = filepath.Join(homeDir, ".cache", "skyterm", "images")
	}

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(dir, cacheDirPerm); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &ImageCache{dir: dir}, nil
}

// Get retrieves image info from cache
func (c *ImageCache) Get(objectName string) (*WikipediaImageInfo, error) {
	cacheFile := c.getCacheFilePath(objectName)

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("cache miss")
		}
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache entry: %w", err)
	}

	// Check cache version
	if entry.Version != cacheVersion {
		return nil, fmt.Errorf("cache version mismatch")
	}

	// Check if cache is stale (older than 30 days)
	if time.Since(entry.Timestamp) > 30*24*time.Hour {
		return nil, fmt.Errorf("cache entry is stale")
	}

	return entry.ImageInfo, nil
}

// Set stores image info in cache
func (c *ImageCache) Set(objectName string, info *WikipediaImageInfo) error {
	entry := cacheEntry{
		Version:    cacheVersion,
		Timestamp:  time.Now(),
		ObjectName: objectName,
		ImageInfo:  info,
	}

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache entry: %w", err)
	}

	cacheFile := c.getCacheFilePath(objectName)
	if err := os.WriteFile(cacheFile, data, cacheFilePerm); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

// getCacheFilePath returns the cache file path for an object
func (c *ImageCache) getCacheFilePath(objectName string) string {
	// Sanitize object name for use as filename
	sanitized := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, objectName)

	return filepath.Join(c.dir, sanitized+".json")
}

// Clear removes all cached entries
func (c *ImageCache) Clear() error {
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			filePath := filepath.Join(c.dir, entry.Name())
			if err := os.Remove(filePath); err != nil {
				return fmt.Errorf("failed to remove cache file %s: %w", filePath, err)
			}
		}
	}

	return nil
}

// Size returns the number of cached entries
func (c *ImageCache) Size() (int, error) {
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		return 0, fmt.Errorf("failed to read cache directory: %w", err)
	}

	count := 0
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			count++
		}
	}

	return count, nil
}
