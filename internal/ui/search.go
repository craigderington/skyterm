package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderSearchBox renders a search input box
func RenderSearchBox(query string, width, height int) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("cyan")).
		Bold(true)

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("green"))

	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("231")).
		Background(lipgloss.Color("236"))

	var content string
	content += titleStyle.Render("Search Objects") + "\n\n"
	content += promptStyle.Render("Enter object name: ")
	content += inputStyle.Render(query + "█") + "\n\n"
	content += lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press Enter to search, Esc to cancel")

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("cyan")).
		Padding(1, 2).
		Width(50)

	box := boxStyle.Render(content)

	// Center the search box
	centered := lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)

	return centered
}

// SearchObjects performs a search across all object types
func SearchObjects(query string, stars []string, planets []string, deepsky []string) (objectType, name string, found bool) {
	query = strings.ToLower(strings.TrimSpace(query))

	if query == "" {
		return "", "", false
	}

	// Search stars
	for _, star := range stars {
		if strings.Contains(strings.ToLower(star), query) {
			return "star", star, true
		}
	}

	// Search planets
	for _, planet := range planets {
		if strings.Contains(strings.ToLower(planet), query) {
			return "planet", planet, true
		}
	}

	// Search deep sky
	for _, obj := range deepsky {
		if strings.Contains(strings.ToLower(obj), query) {
			return "deepsky", obj, true
		}
	}

	return "", "", false
}
