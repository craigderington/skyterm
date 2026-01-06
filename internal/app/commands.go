package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/craigderington/skyterm/internal/image"
)

// TickMsg is sent on every tick to update the time
type TickMsg time.Time

// tickCmd returns a command that sends a tick message
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// ImageFetchedMsg is sent when an image has been fetched
type ImageFetchedMsg struct {
	ObjectName   string
	ImageInfo    *image.WikipediaImageInfo
	RenderedData string
	Error        error
}

// fetchImageCmd fetches and renders an image for an object
func fetchImageCmd(objectName string) tea.Cmd {
	return func() tea.Msg {
		// Create Wikipedia client
		client, err := image.NewWikipediaClient("")
		if err != nil {
			return ImageFetchedMsg{
				ObjectName: objectName,
				Error:      err,
			}
		}

		// Fetch image info
		info, err := client.FetchImage(objectName)
		if err != nil {
			return ImageFetchedMsg{
				ObjectName: objectName,
				Error:      err,
			}
		}

		// Download image
		imgData, err := client.DownloadImage(info.URL)
		if err != nil {
			return ImageFetchedMsg{
				ObjectName: objectName,
				ImageInfo:  info,
				Error:      err,
			}
		}

		// Render for terminal
		renderer := image.NewTerminalRenderer()
		// Use small size for info panel to avoid overflow
		// Max 24 chars wide, 8 rows tall (fits nicely above text info)
		width, height := image.CalculateAspectRatio(info.Width, info.Height, 24, 8)
		rendered, err := renderer.Render(imgData, width, height)
		if err != nil {
			return ImageFetchedMsg{
				ObjectName: objectName,
				ImageInfo:  info,
				Error:      err,
			}
		}

		return ImageFetchedMsg{
			ObjectName:   objectName,
			ImageInfo:    info,
			RenderedData: rendered,
			Error:        nil,
		}
	}
}
