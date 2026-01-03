package render

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/craigderington/skyterm/internal/catalog"
)

// RenderConstellations draws constellation lines on the canvas
func RenderConstellations(canvas *Canvas, stars []catalog.Star, constellations []catalog.Constellation, centerAlt, centerAz, fov float64) {
	// Create a map of star names to stars for quick lookup
	starMap := make(map[string]catalog.Star)
	for _, star := range stars {
		starMap[star.Name] = star
	}

	lineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	lineChar := '─'

	for _, constellation := range constellations {
		for _, line := range constellation.Lines {
			star1, ok1 := starMap[line.Star1Name]
			star2, ok2 := starMap[line.Star2Name]

			if !ok1 || !ok2 {
				continue
			}

			// Skip if same star (placeholder lines)
			if line.Star1Name == line.Star2Name {
				continue
			}

			// Project both stars
			x1, y1, visible1 := Project(star1.Altitude, star1.Azimuth, centerAlt, centerAz, fov, canvas.Width, canvas.Height)
			x2, y2, visible2 := Project(star2.Altitude, star2.Azimuth, centerAlt, centerAz, fov, canvas.Width, canvas.Height)

			// Only draw if both stars are visible
			if !visible1 || !visible2 {
				continue
			}

			// Draw line using Bresenham's algorithm
			drawLine(canvas, x1, y1, x2, y2, lineChar, lineStyle)
		}
	}
}

// RenderConstellationLabels draws constellation name labels
func RenderConstellationLabels(canvas *Canvas, stars []catalog.Star, labels []catalog.ConstellationLabel, centerAlt, centerAz, fov float64) {
	starMap := make(map[string]catalog.Star)
	for _, star := range stars {
		starMap[star.Name] = star
	}

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("51")). // Bright cyan
		Bold(true)

	for _, label := range labels {
		star, ok := starMap[label.StarName]
		if !ok {
			continue
		}

		x, y, visible := Project(star.Altitude, star.Azimuth, centerAlt, centerAz, fov, canvas.Width, canvas.Height)
		if !visible {
			continue
		}

		// Offset label to the right and slightly up from the star
		labelX := x + 2
		labelY := y - 1

		// Check if label position is valid
		if labelY < 0 || labelY >= canvas.Height {
			// Try placing it below instead
			labelY = y + 1
			if labelY >= canvas.Height {
				continue
			}
		}

		// Draw label, checking each character position
		for i, ch := range label.Name {
			xPos := labelX + i
			if xPos >= 0 && xPos < canvas.Width {
				canvas.Set(xPos, labelY, ch, labelStyle)
			}
		}
	}
}

// drawLine uses Bresenham's line algorithm to draw a line between two points
func drawLine(canvas *Canvas, x1, y1, x2, y2 int, char rune, style lipgloss.Style) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)

	sx := -1
	if x1 < x2 {
		sx = 1
	}

	sy := -1
	if y1 < y2 {
		sy = 1
	}

	err := dx - dy

	x, y := x1, y1

	for {
		// Set pixel, but don't overwrite stars
		cell := canvas.Cells[y][x]
		if cell.Char == ' ' || cell.Char == char {
			canvas.Set(x, y, char, style)
		}

		if x == x2 && y == y2 {
			break
		}

		e2 := 2 * err

		if e2 > -dy {
			err -= dy
			x += sx
		}

		if e2 < dx {
			err += dx
			y += sy
		}

		// Safety check to prevent infinite loops
		if x < 0 || x >= canvas.Width || y < 0 || y >= canvas.Height {
			break
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
