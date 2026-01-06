package main

import (
	"fmt"
	"os"

	"github.com/craigderington/skyterm/internal/image"
)

func main() {
	fmt.Println("=== Terminal Capability Detection ===")
	fmt.Println()

	// Show environment variables
	fmt.Println("Environment Variables:")
	fmt.Printf("  TERM           = %q\n", os.Getenv("TERM"))
	fmt.Printf("  TERM_PROGRAM   = %q\n", os.Getenv("TERM_PROGRAM"))
	fmt.Printf("  KITTY_WINDOW_ID = %q\n", os.Getenv("KITTY_WINDOW_ID"))
	fmt.Println()

	// Detect capability
	renderer := image.NewTerminalRenderer()
	capability := renderer.GetCapability()

	fmt.Println("Detected Capability:")
	fmt.Printf("  %s\n", renderer.CapabilityString())
	fmt.Println()

	// Show what this means
	switch capability {
	case 3: // CapabilityKitty
		fmt.Println("✓ Kitty Graphics Protocol ENABLED")
		fmt.Println("  You should see high-resolution images!")
	case 4: // CapabilityITerm2
		fmt.Println("✓ iTerm2 Inline Images ENABLED")
		fmt.Println("  You should see high-resolution images!")
	case 2: // CapabilitySixel
		fmt.Println("✓ Sixel Graphics ENABLED")
		fmt.Println("  You should see good quality images!")
	case 1: // CapabilityUnicode
		fmt.Println("⚠ Unicode Block Art (fallback mode)")
		fmt.Println("  Images will be pixelated")
		fmt.Println()
		fmt.Println("For better image quality:")
		fmt.Println("  - Use Kitty terminal (https://sw.kovidgoyal.net/kitty/)")
		fmt.Println("  - Use iTerm2 on macOS")
		fmt.Println("  - Use a terminal with Sixel support")
	default:
		fmt.Println("⚠ No image support detected")
	}
	fmt.Println()

	// Test with sample object
	if len(os.Args) > 1 {
		objectName := os.Args[1]
		fmt.Printf("Testing image fetch for: %s\n", objectName)
		fmt.Println("Fetching from Wikipedia...")

		info, err := image.GetImageForObject(objectName, "")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Found: %s\n", info.Title)
		fmt.Printf("  URL: %s\n", info.URL)
		fmt.Printf("  Size: %dx%d\n", info.Width, info.Height)
		fmt.Println()

		// Try to render (small size for testing)
		fmt.Println("Attempting to render image (10x10 for test)...")
		rendered, err := renderer.RenderImageFromURL(info.URL, 10, 10)
		if err != nil {
			fmt.Printf("Render error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("=== Rendered Image ===")
		fmt.Print(rendered)
		fmt.Println("=== End Image ===")
	} else {
		fmt.Println("Usage: termtest <object-name>")
		fmt.Println("Example: termtest Sirius")
	}
}
