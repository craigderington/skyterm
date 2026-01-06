package image

import (
	"testing"
)

func TestWikipediaClient_FetchImage(t *testing.T) {
	// Skip in CI or when offline
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewWikipediaClient("")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tests := []struct {
		name       string
		objectName string
		wantError  bool
	}{
		{
			name:       "Sirius",
			objectName: "Sirius",
			wantError:  false,
		},
		{
			name:       "Andromeda Galaxy",
			objectName: "Andromeda Galaxy",
			wantError:  false,
		},
		{
			name:       "M31 (Messier notation)",
			objectName: "M31",
			wantError:  false,
		},
		{
			name:       "Orion Nebula",
			objectName: "Orion Nebula",
			wantError:  false,
		},
		{
			name:       "Jupiter",
			objectName: "Jupiter",
			wantError:  false,
		},
		{
			name:       "Alpha Centauri",
			objectName: "Alpha Centauri",
			wantError:  false,
		},
		{
			name:       "Nonexistent object",
			objectName: "XYZ999ZZZ",
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := client.FetchImage(tt.objectName)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error for %q, got nil", tt.objectName)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error for %q: %v", tt.objectName, err)
				return
			}

			if info == nil {
				t.Errorf("expected image info for %q, got nil", tt.objectName)
				return
			}

			if info.URL == "" {
				t.Errorf("expected non-empty URL for %q", tt.objectName)
			}

			if info.Title == "" {
				t.Errorf("expected non-empty title for %q", tt.objectName)
			}

			if info.Source == "" {
				t.Errorf("expected non-empty source URL for %q", tt.objectName)
			}

			t.Logf("Found image for %q:", tt.objectName)
			t.Logf("  Title: %s", info.Title)
			t.Logf("  URL: %s", info.URL)
			t.Logf("  Size: %dx%d", info.Width, info.Height)
			t.Logf("  Description: %s", info.Description)
			t.Logf("  Source: %s", info.Source)
		})
	}
}

func TestGenerateNameVariants(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{
			input: "M31",
			want:  []string{"M31", "M31 (star)", "M31 (constellation)", "M31 (galaxy)", "M31 (nebula)", "M31 (planet)", "Messier 31"},
		},
		{
			input: "NGC224",
			want:  []string{"NGC224", "NGC224 (star)", "NGC224 (constellation)", "NGC224 (galaxy)", "NGC224 (nebula)", "NGC224 (planet)", "NGC 224"},
		},
		{
			input: "Alpha Centauri",
			want:  []string{"Alpha Centauri", "Alpha Centauri (star)", "Alpha Centauri (constellation)", "Alpha Centauri (galaxy)", "Alpha Centauri (nebula)", "Alpha Centauri (planet)", "α Centauri"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := generateNameVariants(tt.input)

			// Check that all expected variants are present
			for _, expected := range tt.want {
				found := false
				for _, variant := range got {
					if variant == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected variant %q not found in results", expected)
				}
			}
		})
	}
}

func TestImageCache(t *testing.T) {
	// Create temporary cache directory
	tmpDir := t.TempDir()

	cache, err := NewImageCache(tmpDir)
	if err != nil {
		t.Fatalf("failed to create cache: %v", err)
	}

	// Test cache miss
	_, err = cache.Get("test-object")
	if err == nil {
		t.Error("expected cache miss error, got nil")
	}

	// Test cache set and get
	testInfo := &WikipediaImageInfo{
		URL:         "https://example.com/image.jpg",
		Title:       "Test Object",
		Description: "A test astronomical object",
		Width:       1024,
		Height:      768,
		Source:      "https://en.wikipedia.org/wiki/Test",
	}

	if err := cache.Set("test-object", testInfo); err != nil {
		t.Fatalf("failed to set cache entry: %v", err)
	}

	retrieved, err := cache.Get("test-object")
	if err != nil {
		t.Fatalf("failed to get cache entry: %v", err)
	}

	if retrieved.URL != testInfo.URL {
		t.Errorf("URL mismatch: got %q, want %q", retrieved.URL, testInfo.URL)
	}

	if retrieved.Title != testInfo.Title {
		t.Errorf("Title mismatch: got %q, want %q", retrieved.Title, testInfo.Title)
	}

	// Test cache size
	size, err := cache.Size()
	if err != nil {
		t.Fatalf("failed to get cache size: %v", err)
	}

	if size != 1 {
		t.Errorf("expected cache size 1, got %d", size)
	}

	// Test cache clear
	if err := cache.Clear(); err != nil {
		t.Fatalf("failed to clear cache: %v", err)
	}

	size, err = cache.Size()
	if err != nil {
		t.Fatalf("failed to get cache size after clear: %v", err)
	}

	if size != 0 {
		t.Errorf("expected cache size 0 after clear, got %d", size)
	}
}
