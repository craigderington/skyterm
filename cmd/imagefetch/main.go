package main

import (
	"fmt"
	"os"

	"github.com/craigderington/skyterm/internal/image"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: imagefetch <object-name>")
		fmt.Println("\nExamples:")
		fmt.Println("  imagefetch Sirius")
		fmt.Println("  imagefetch \"Andromeda Galaxy\"")
		fmt.Println("  imagefetch M31")
		fmt.Println("  imagefetch Jupiter")
		os.Exit(1)
	}

	objectName := os.Args[1]

	fmt.Printf("Fetching image for: %s\n\n", objectName)

	// Create Wikipedia client
	client, err := image.NewWikipediaClient("")
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	// Fetch image info
	info, err := client.FetchImage(objectName)
	if err != nil {
		fmt.Printf("Error fetching image: %v\n", err)
		os.Exit(1)
	}

	// Display results
	fmt.Println("✓ Found image!")
	fmt.Println()
	fmt.Printf("Title:       %s\n", info.Title)
	fmt.Printf("Description: %s\n", info.Description)
	fmt.Printf("Image URL:   %s\n", info.URL)
	fmt.Printf("Dimensions:  %d × %d\n", info.Width, info.Height)
	fmt.Printf("Article:     %s\n", info.Source)
	fmt.Println()

	// Ask if user wants to download the image
	fmt.Print("Download image? (y/n): ")
	var response string
	fmt.Scanln(&response)

	if response == "y" || response == "Y" {
		fmt.Println("Downloading...")
		data, err := client.DownloadImage(info.URL)
		if err != nil {
			fmt.Printf("Error downloading image: %v\n", err)
			os.Exit(1)
		}

		// Save to file
		filename := fmt.Sprintf("%s.jpg", objectName)
		if err := os.WriteFile(filename, data, 0644); err != nil {
			fmt.Printf("Error saving image: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Saved to: %s (%d bytes)\n", filename, len(data))
	}
}
