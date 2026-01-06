# Image Package

This package provides Wikipedia image fetching functionality for astronomical objects in skyterm.

## Features

- **Wikipedia API Integration**: Fetches images from Wikipedia articles
- **Smart Name Resolution**: Handles multiple naming formats:
  - Messier objects (M31 → Messier 31 → Andromeda Galaxy)
  - NGC/IC catalog objects (NGC224 → NGC 224)
  - Greek letter designations (Alpha Centauri → α Centauri)
  - Common name variants (star, galaxy, nebula, planet, cluster)
- **Disambiguation Avoidance**: Automatically skips disambiguation pages
- **Redirect Following**: Follows Wikipedia redirects automatically
- **Local Caching**: Caches image metadata locally to reduce API calls
- **Fallback Image Detection**: If main image isn't available, finds first suitable image on page

## Usage

### Basic Example

```go
import "github.com/craigderington/skyterm/internal/image"

// Create a client
client, err := image.NewWikipediaClient("")
if err != nil {
    log.Fatal(err)
}

// Fetch image info
info, err := client.FetchImage("Sirius")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Image URL: %s\n", info.URL)
fmt.Printf("Title: %s\n", info.Title)
fmt.Printf("Description: %s\n", info.Description)

// Download the actual image
imageData, err := client.DownloadImage(info.URL)
if err != nil {
    log.Fatal(err)
}
```

### Convenience Function

```go
// Quick one-liner to get image info
info, err := image.GetImageForObject("M31", "")
```

## Cache

Images are cached in `~/.cache/skyterm/images/` by default. Cache entries:
- Are stored as JSON files
- Include timestamp and version info
- Expire after 30 days
- Can be cleared with `cache.Clear()`

## Custom Cache Directory

```go
client, err := image.NewWikipediaClient("/custom/cache/path")
```

## Testing

```bash
# Run unit tests (fast)
go test ./internal/image/ -v -short

# Run integration tests (requires internet)
go test ./internal/image/ -v

# Test with example program
go build ./cmd/imagefetch
./imagefetch "Andromeda Galaxy"
./imagefetch "M42"
./imagefetch "Jupiter"
```

## WikipediaImageInfo Structure

```go
type WikipediaImageInfo struct {
    URL         string  // Direct URL to image file
    Title       string  // Wikipedia article title
    Description string  // Short description from Wikipedia
    Width       int     // Image width in pixels (0 if unknown)
    Height      int     // Image height in pixels (0 if unknown)
    Source      string  // Wikipedia article URL
}
```

## Supported Object Types

The fetcher works well with:
- **Stars**: Sirius, Betelgeuse, Alpha Centauri, etc.
- **Planets**: Jupiter, Mars, Saturn, etc.
- **Galaxies**: Andromeda Galaxy, M31, NGC 224, etc.
- **Nebulae**: Orion Nebula, M42, Crab Nebula, etc.
- **Star Clusters**: Pleiades, M45, Hyades, etc.
- **Constellations**: Orion, Ursa Major, etc.

## Error Handling

Common errors:
- `"no Wikipedia article found"` - Object name not recognized, try alternate name
- `"no image found for article"` - Article exists but has no images
- `"cache miss"` - Not an error, just means not in cache
- `"cache entry is stale"` - Cache entry too old, will refetch

## Performance

- First fetch: ~500ms-2s (network call to Wikipedia API)
- Cached fetch: <1ms (read from local JSON)
- Cache size: ~300-500 bytes per entry

## Future Enhancements

Potential improvements:
- Support for other image sources (NASA, ESO, etc.)
- Image quality preferences (thumbnail, medium, full resolution)
- Batch fetching for multiple objects
- Pre-warming cache with common objects
- Image format preferences (JPG, PNG, SVG)
