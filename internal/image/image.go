package image

import (
	"fmt"
)

// Fetcher is an interface for fetching images
type Fetcher interface {
	// FetchImage retrieves image information for the given object name
	FetchImage(objectName string) (*WikipediaImageInfo, error)

	// DownloadImage downloads the actual image data from a URL
	DownloadImage(imageURL string) ([]byte, error)
}

// GetImageForObject is a convenience function to get image info for an astronomical object
// It tries Wikipedia as the primary source
func GetImageForObject(objectName string, cacheDir string) (*WikipediaImageInfo, error) {
	client, err := NewWikipediaClient(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create Wikipedia client: %w", err)
	}

	return client.FetchImage(objectName)
}
