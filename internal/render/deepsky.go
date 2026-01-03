package render

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/craigderington/skyterm/internal/catalog"
)

// RenderDeepSkyObjects draws Messier objects on the canvas
func RenderDeepSkyObjects(canvas *Canvas, objects []catalog.MessierObject, centerAlt, centerAz, fov, magLimit float64) {
	// Object type symbols and colors
	typeStyles := map[string]struct {
		char  rune
		color lipgloss.Color
	}{
		"Galaxy":              {'◈', lipgloss.Color("141")}, // Purple
		"Nebula":              {'◇', lipgloss.Color("213")}, // Pink
		"Supernova Remnant":   {'✸', lipgloss.Color("196")}, // Red
		"Globular Cluster":    {'◉', lipgloss.Color("220")}, // Gold
		"Open Cluster":        {'◌', lipgloss.Color("117")}, // Light blue
		"Planetary Nebula":    {'◎', lipgloss.Color("48")},  // Cyan
	}

	for _, obj := range objects {
		// Skip if too dim
		if obj.Magnitude > magLimit+3 {
			continue
		}

		// Project to screen
		x, y, visible := Project(obj.Altitude, obj.Azimuth, centerAlt, centerAz, fov, canvas.Width, canvas.Height)
		if !visible {
			continue
		}

		// Get style for object type
		style, ok := typeStyles[obj.Type]
		if !ok {
			// Default style
			style.char = '◆'
			style.color = lipgloss.Color("245")
		}

		objStyle := lipgloss.NewStyle().Foreground(style.color)

		canvas.Set(x, y, style.char, objStyle)
	}
}

// RenderDeepSkyLabels draws labels for Messier objects
func RenderDeepSkyLabels(canvas *Canvas, objects []catalog.MessierObject, centerAlt, centerAz, fov, magLimit float64) {
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("magenta")).
		Faint(true)

	for _, obj := range objects {
		// Skip if too dim
		if obj.Magnitude > magLimit+3 {
			continue
		}

		// Project to screen
		x, y, visible := Project(obj.Altitude, obj.Azimuth, centerAlt, centerAz, fov, canvas.Width, canvas.Height)
		if !visible {
			continue
		}

		// Create label (M number)
		label := fmt.Sprintf("M%d", obj.Number)

		// Position label
		labelX := x + 2
		labelY := y

		if labelY < 0 || labelY >= canvas.Height {
			continue
		}

		// Draw label
		for i, ch := range label {
			xPos := labelX + i
			if xPos >= 0 && xPos < canvas.Width {
				canvas.Set(xPos, labelY, ch, labelStyle)
			}
		}
	}
}
