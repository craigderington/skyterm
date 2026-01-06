package image

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	wikipediaAPIURL = "https://en.wikipedia.org/w/api.php"
	userAgent       = "skyterm/1.0 (https://github.com/craigderington/skyterm)"
	requestTimeout  = 10 * time.Second
)

// WikipediaClient handles fetching images from Wikipedia
type WikipediaClient struct {
	httpClient *http.Client
	cache      *ImageCache
}

// NewWikipediaClient creates a new Wikipedia API client
func NewWikipediaClient(cacheDir string) (*WikipediaClient, error) {
	cache, err := NewImageCache(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}

	return &WikipediaClient{
		httpClient: &http.Client{
			Timeout: requestTimeout,
		},
		cache: cache,
	}, nil
}

// WikipediaImageInfo contains information about an image from Wikipedia
type WikipediaImageInfo struct {
	URL         string
	Title       string
	Description string
	Width       int
	Height      int
	Source      string // Wikipedia article URL
}

// FetchImage fetches an image for the given astronomical object
// Returns the image info, or an error if not found
func (w *WikipediaClient) FetchImage(objectName string) (*WikipediaImageInfo, error) {
	// Check cache first
	if cachedInfo, err := w.cache.Get(objectName); err == nil {
		return cachedInfo, nil
	}

	// Search Wikipedia for the article
	pageTitle, err := w.searchArticle(objectName)
	if err != nil {
		return nil, fmt.Errorf("failed to find Wikipedia article: %w", err)
	}

	// Get the main image from the article
	imageInfo, err := w.getPageImage(pageTitle)
	if err != nil {
		return nil, fmt.Errorf("failed to get image from article: %w", err)
	}

	// Add article URL
	imageInfo.Source = fmt.Sprintf("https://en.wikipedia.org/wiki/%s", url.PathEscape(pageTitle))

	// Cache the result
	if err := w.cache.Set(objectName, imageInfo); err != nil {
		// Log but don't fail on cache errors
		fmt.Printf("Warning: failed to cache image info: %v\n", err)
	}

	return imageInfo, nil
}

// DownloadImage downloads the actual image data
func (w *WikipediaClient) DownloadImage(imageURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	return data, nil
}

// searchArticle searches Wikipedia for an article matching the object name
func (w *WikipediaClient) searchArticle(objectName string) (string, error) {
	// Try with common astronomical suffixes/prefixes FIRST
	// This helps avoid disambiguation pages (e.g., M31 -> Messier 31)
	variants := generateNameVariants(objectName)
	for _, variant := range variants {
		// Skip the original name on first pass if it looks like it might be ambiguous
		if variant == objectName && w.mightBeAmbiguous(objectName) {
			continue
		}

		if exists, err := w.articleExists(variant); err == nil && exists {
			// Check if this is a disambiguation page
			if isDisambig, _ := w.isDisambiguationPage(variant); isDisambig {
				continue // Try next variant
			}
			return variant, nil
		}
	}

	// If variants didn't work, try exact match as fallback
	if exists, err := w.articleExists(objectName); err == nil && exists {
		if isDisambig, _ := w.isDisambiguationPage(objectName); !isDisambig {
			return objectName, nil
		}
	}

	// Fall back to search API
	params := url.Values{
		"action":     {"query"},
		"list":       {"search"},
		"srsearch":   {objectName},
		"format":     {"json"},
		"srlimit":    {"1"},
		"srprop":     {""},
		"utf8":       {"1"},
	}

	var result struct {
		Query struct {
			Search []struct {
				Title string `json:"title"`
			} `json:"search"`
		} `json:"query"`
	}

	if err := w.apiRequest(params, &result); err != nil {
		return "", err
	}

	if len(result.Query.Search) == 0 {
		return "", fmt.Errorf("no Wikipedia article found for %q", objectName)
	}

	return result.Query.Search[0].Title, nil
}

// mightBeAmbiguous returns true if the object name might lead to a disambiguation page
func (w *WikipediaClient) mightBeAmbiguous(name string) bool {
	// Messier objects in M## format are often disambiguation pages
	if strings.HasPrefix(strings.ToUpper(name), "M") && len(name) <= 4 {
		return true
	}
	return false
}

// isDisambiguationPage checks if a page is a disambiguation page
func (w *WikipediaClient) isDisambiguationPage(title string) (bool, error) {
	params := url.Values{
		"action": {"query"},
		"titles": {title},
		"prop":   {"pageprops"},
		"format": {"json"},
	}

	var result struct {
		Query struct {
			Pages map[string]struct {
				PageProps struct {
					Disambiguation string `json:"disambiguation"`
				} `json:"pageprops"`
			} `json:"pages"`
		} `json:"query"`
	}

	if err := w.apiRequest(params, &result); err != nil {
		return false, err
	}

	for _, page := range result.Query.Pages {
		// If the page has the "disambiguation" property, it's a disambiguation page
		return page.PageProps.Disambiguation != "", nil
	}

	return false, nil
}

// articleExists checks if an article with the exact title exists
func (w *WikipediaClient) articleExists(title string) (bool, error) {
	params := url.Values{
		"action": {"query"},
		"titles": {title},
		"format": {"json"},
	}

	var result struct {
		Query struct {
			Pages map[string]struct {
				Missing bool `json:"missing,omitempty"`
			} `json:"pages"`
		} `json:"query"`
	}

	if err := w.apiRequest(params, &result); err != nil {
		return false, err
	}

	for _, page := range result.Query.Pages {
		return !page.Missing, nil
	}

	return false, nil
}

// getPageImage retrieves the main image from a Wikipedia article
func (w *WikipediaClient) getPageImage(pageTitle string) (*WikipediaImageInfo, error) {
	// First try with redirects enabled to handle cases like M31 -> Andromeda Galaxy
	params := url.Values{
		"action":     {"query"},
		"titles":     {pageTitle},
		"prop":       {"pageimages|pageterms"},
		"piprop":     {"original|name"},
		"redirects":  {"1"}, // Follow redirects
		"format":     {"json"},
		"utf8":       {"1"},
	}

	var result struct {
		Query struct {
			Pages map[string]struct {
				Title    string `json:"title"`
				Original struct {
					Source string `json:"source"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"original"`
				PageImage string `json:"pageimage"`
				Terms     struct {
					Description []string `json:"description"`
				} `json:"terms"`
			} `json:"pages"`
			Redirects []struct {
				From string `json:"from"`
				To   string `json:"to"`
			} `json:"redirects"`
		} `json:"query"`
	}

	if err := w.apiRequest(params, &result); err != nil {
		return nil, err
	}

	for _, page := range result.Query.Pages {
		// If no original image, try to get the first image from the page
		if page.Original.Source == "" {
			// Try alternative approach: get images from page
			imageURL, err := w.getFirstPageImage(page.Title)
			if err != nil {
				return nil, fmt.Errorf("no image found for article %q", pageTitle)
			}

			description := ""
			if len(page.Terms.Description) > 0 {
				description = page.Terms.Description[0]
			}

			return &WikipediaImageInfo{
				URL:         imageURL,
				Title:       page.Title,
				Description: description,
				Width:       0, // Unknown
				Height:      0, // Unknown
			}, nil
		}

		description := ""
		if len(page.Terms.Description) > 0 {
			description = page.Terms.Description[0]
		}

		return &WikipediaImageInfo{
			URL:         page.Original.Source,
			Title:       page.Title,
			Description: description,
			Width:       page.Original.Width,
			Height:      page.Original.Height,
		}, nil
	}

	return nil, fmt.Errorf("no pages found in response")
}

// getFirstPageImage gets the first image from a Wikipedia page
func (w *WikipediaClient) getFirstPageImage(pageTitle string) (string, error) {
	params := url.Values{
		"action":  {"query"},
		"titles":  {pageTitle},
		"prop":    {"images"},
		"imlimit": {"10"}, // Get first 10 images
		"format":  {"json"},
	}

	var result struct {
		Query struct {
			Pages map[string]struct {
				Images []struct {
					Title string `json:"title"`
				} `json:"images"`
			} `json:"pages"`
		} `json:"query"`
	}

	if err := w.apiRequest(params, &result); err != nil {
		return "", err
	}

	for _, page := range result.Query.Pages {
		if len(page.Images) == 0 {
			return "", fmt.Errorf("no images found")
		}

		// Get the URL of the first image (skip common icons/logos)
		for _, img := range page.Images {
			// Skip common non-content images
			imgTitle := strings.ToLower(img.Title)
			if strings.Contains(imgTitle, "commons-logo") ||
				strings.Contains(imgTitle, "wikidata") ||
				strings.Contains(imgTitle, "edit-clear") ||
				strings.Contains(imgTitle, "question_book") {
				continue
			}

			// Get the actual image URL
			imageURL, err := w.getImageURL(img.Title)
			if err == nil {
				return imageURL, nil
			}
		}
	}

	return "", fmt.Errorf("no suitable image found")
}

// getImageURL gets the URL of an image file from Wikipedia
func (w *WikipediaClient) getImageURL(imageTitle string) (string, error) {
	params := url.Values{
		"action": {"query"},
		"titles": {imageTitle},
		"prop":   {"imageinfo"},
		"iiprop": {"url"},
		"format": {"json"},
	}

	var result struct {
		Query struct {
			Pages map[string]struct {
				ImageInfo []struct {
					URL string `json:"url"`
				} `json:"imageinfo"`
			} `json:"pages"`
		} `json:"query"`
	}

	if err := w.apiRequest(params, &result); err != nil {
		return "", err
	}

	for _, page := range result.Query.Pages {
		if len(page.ImageInfo) > 0 {
			return page.ImageInfo[0].URL, nil
		}
	}

	return "", fmt.Errorf("no image URL found")
}

// apiRequest makes a request to the Wikipedia API
func (w *WikipediaClient) apiRequest(params url.Values, result interface{}) error {
	reqURL := wikipediaAPIURL + "?" + params.Encode()

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// generateNameVariants generates common variations of astronomical object names
func generateNameVariants(name string) []string {
	variants := []string{}

	// Handle Messier objects FIRST (M31 -> Messier 31)
	// This avoids disambiguation pages like "M31"
	if strings.HasPrefix(strings.ToUpper(name), "M") && len(name) > 1 {
		if num := name[1:]; num != "" {
			// Try "Messier X" first as it's more specific
			variants = append(variants, "Messier "+num)
		}
	}

	// Add original name
	variants = append(variants, name)

	// Handle NGC objects
	if strings.HasPrefix(strings.ToUpper(name), "NGC") {
		variants = append(variants, strings.ReplaceAll(name, "NGC", "NGC "))
	}

	// Handle IC objects (Index Catalogue)
	if strings.HasPrefix(strings.ToUpper(name), "IC") && len(name) > 2 {
		variants = append(variants, strings.ReplaceAll(name, "IC", "IC "))
	}

	// Handle Greek letters (alpha -> α)
	greekMap := map[string]string{
		"alpha": "α", "beta": "β", "gamma": "γ", "delta": "δ",
		"epsilon": "ε", "zeta": "ζ", "eta": "η", "theta": "θ",
		"iota": "ι", "kappa": "κ", "lambda": "λ", "mu": "μ",
		"nu": "ν", "xi": "ξ", "omicron": "ο", "pi": "π",
		"rho": "ρ", "sigma": "σ", "tau": "τ", "upsilon": "υ",
		"phi": "φ", "chi": "χ", "psi": "ψ", "omega": "ω",
	}
	lowerName := strings.ToLower(name)
	for latin, greek := range greekMap {
		if strings.Contains(lowerName, latin) {
			variants = append(variants, strings.ReplaceAll(name, latin, greek))
			variants = append(variants, strings.ReplaceAll(name, strings.Title(latin), greek))
		}
	}

	// Add parenthetical clarifications (last resort)
	variants = append(variants, name+" (star)")
	variants = append(variants, name+" (constellation)")
	variants = append(variants, name+" (galaxy)")
	variants = append(variants, name+" (nebula)")
	variants = append(variants, name+" (planet)")
	variants = append(variants, name+" (cluster)")

	return variants
}
