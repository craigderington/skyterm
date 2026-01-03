package render

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/craigderington/skyterm/internal/catalog"
)

// RenderStarLabels renders labels for bright stars
func RenderStarLabels(canvas *Canvas, stars []catalog.Star, centerAlt, centerAz, fov, magLimit float64) {
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("white")).
		Faint(true)

	for _, star := range stars {
		// Only label bright stars (magnitude < 2.5) or famous stars
		if star.Magnitude > 2.5 {
			continue
		}

		// Skip if star itself is not visible
		if star.Magnitude > magLimit {
			continue
		}

		// Project star to screen
		x, y, visible := Project(star.Altitude, star.Azimuth, centerAlt, centerAz, fov, canvas.Width, canvas.Height)
		if !visible {
			continue
		}

		// Position label to the right and slightly up
		labelX := x + 2
		labelY := y

		// Check bounds
		if labelY < 0 || labelY >= canvas.Height {
			continue
		}

		// Draw label
		for i, ch := range star.Name {
			xPos := labelX + i
			if xPos >= 0 && xPos < canvas.Width {
				canvas.Set(xPos, labelY, ch, labelStyle)
			}
		}
	}
}
