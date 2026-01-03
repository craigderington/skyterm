package render

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/craigderington/skyterm/internal/astro"
)

// RenderPlanets draws planets on the canvas
func RenderPlanets(canvas *Canvas, planets *astro.PlanetarySystem, centerAlt, centerAz, fov float64) {
	if planets == nil {
		return
	}

	// Planet colors and symbols
	planetStyles := map[string]struct {
		char  rune
		color lipgloss.Color
	}{
		"Sun":     {'☉', lipgloss.Color("226")}, // Bright yellow
		"Moon":    {'☽', lipgloss.Color("250")}, // Light gray
		"Mercury": {'☿', lipgloss.Color("249")}, // Gray
		"Venus":   {'♀', lipgloss.Color("230")}, // Yellowish-white
		"Mars":    {'♂', lipgloss.Color("196")}, // Red
		"Jupiter": {'♃', lipgloss.Color("215")}, // Orange-white
		"Saturn":  {'♄', lipgloss.Color("229")}, // Pale yellow
		"Uranus":  {'♅', lipgloss.Color("117")}, // Pale cyan
		"Neptune": {'♆', lipgloss.Color("27")},  // Blue
	}

	for _, planet := range planets.AllPlanets() {
		// Project planet to screen coordinates
		x, y, visible := Project(planet.Altitude, planet.Azimuth, centerAlt, centerAz, fov, canvas.Width, canvas.Height)
		if !visible {
			continue
		}

		style, ok := planetStyles[planet.Name]
		if !ok {
			continue
		}

		// Render planet symbol with bold
		planetStyle := lipgloss.NewStyle().
			Foreground(style.color).
			Bold(true)

		canvas.Set(x, y, style.char, planetStyle)
	}
}

// RenderPlanetLabels draws planet name labels
func RenderPlanetLabels(canvas *Canvas, planets *astro.PlanetarySystem, centerAlt, centerAz, fov float64) {
	if planets == nil {
		return
	}

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("yellow")).
		Bold(true)

	for _, planet := range planets.AllPlanets() {
		// Project planet to screen coordinates
		x, y, visible := Project(planet.Altitude, planet.Azimuth, centerAlt, centerAz, fov, canvas.Width, canvas.Height)
		if !visible {
			continue
		}

		// Position label to the right of the planet
		labelX := x + 2
		labelY := y

		// Check bounds
		if labelY < 0 || labelY >= canvas.Height {
			continue
		}

		// Draw label
		for i, ch := range planet.Name {
			xPos := labelX + i
			if xPos >= 0 && xPos < canvas.Width {
				canvas.Set(xPos, labelY, ch, labelStyle)
			}
		}
	}
}
