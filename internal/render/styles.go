package render

import "github.com/charmbracelet/lipgloss"

// GetStarStyle returns an enhanced style for a star based on its magnitude and spectral type
func GetStarStyle(magnitude float64, spectralType rune, colorByType bool) lipgloss.Style {
	style := lipgloss.NewStyle()

	if colorByType {
		// Color by spectral type
		color := getColorForSpectralType(spectralType)
		style = style.Foreground(color)
	} else {
		// Default white
		style = style.Foreground(lipgloss.Color("231"))
	}

	// Brightest stars get bold
	if magnitude < 0.5 {
		style = style.Bold(true)
	}

	return style
}

// GetStarChar returns an appropriate character for a star based on its magnitude
func GetStarChar(magnitude float64) rune {
	return getCharForMagnitude(magnitude)
}
