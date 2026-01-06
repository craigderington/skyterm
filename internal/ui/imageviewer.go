package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderImageViewer renders a fullscreen image viewer
func RenderImageViewer(objectInfo *ObjectInfo, width, height int) string {
	if objectInfo == nil || objectInfo.ImageData == "" {
		return ""
	}

	// Create header with object name
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("cyan")).
		Bold(true).
		Background(lipgloss.Color("235")).
		Width(width).
		Align(lipgloss.Center).
		Padding(0, 1)

	header := headerStyle.Render(objectInfo.Name)

	// Create description section if available
	var description string
	if objectInfo.ImageInfo != nil && objectInfo.ImageInfo.Description != "" {
		descStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("250")).
			Width(width - 4).
			Align(lipgloss.Center).
			Padding(0, 2)
		description = descStyle.Render(objectInfo.ImageInfo.Description)
	}

	// Create footer with instructions
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Background(lipgloss.Color("235")).
		Width(width).
		Align(lipgloss.Center).
		Padding(0, 1)

	footer := footerStyle.Render("Press ESC, V, or Q to close")

	// Calculate center position for image
	// Split image into lines to count height
	imageLines := strings.Split(objectInfo.ImageData, "\n")
	imageHeight := len(imageLines)
	imageWidth := 0
	for _, line := range imageLines {
		if len(line) > imageWidth {
			imageWidth = len(line)
		}
	}

	// Center the image
	// Account for header, description (if present), and footer
	descLines := 0
	if description != "" {
		descLines = len(strings.Split(description, "\n")) + 1 // +1 for extra spacing
	}

	topMargin := (height - imageHeight - 3 - descLines) / 2 // -3 for header and footer
	if topMargin < 0 {
		topMargin = 0
	}

	leftMargin := (width - imageWidth) / 2
	if leftMargin < 0 {
		leftMargin = 0
	}

	// Build the view
	var view strings.Builder

	// Add header
	view.WriteString(header + "\n")

	// Add description if available
	if description != "" {
		view.WriteString(description + "\n\n")
	}

	// Add top margin
	for i := 0; i < topMargin; i++ {
		view.WriteString("\n")
	}

	// Add image with left margin
	// Note: Don't filter empty lines - they may contain escape sequences
	for _, line := range imageLines {
		view.WriteString(strings.Repeat(" ", leftMargin))
		view.WriteString(line)
		view.WriteString("\n")
	}

	// Add bottom margin
	remainingLines := height - topMargin - imageHeight - 2
	for i := 0; i < remainingLines && i >= 0; i++ {
		view.WriteString("\n")
	}

	// Add footer
	view.WriteString(footer)

	// Add Wikipedia source info if available
	if objectInfo.ImageInfo != nil {
		infoStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Faint(true).
			Width(width).
			Align(lipgloss.Center)

		sourceInfo := fmt.Sprintf("Source: Wikipedia (%s) | Image: %dx%d",
			objectInfo.ImageInfo.Title,
			objectInfo.ImageInfo.Width,
			objectInfo.ImageInfo.Height,
		)

		view.WriteString("\n")
		view.WriteString(infoStyle.Render(sourceInfo) + "\n")
	}

	return view.String()
}
