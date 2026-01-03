package render

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

// RenderGrid draws a coordinate grid overlay showing altitude and azimuth
func RenderGrid(canvas *Canvas, centerAlt, centerAz, fov float64) {
	gridStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("242"))

	// Draw altitude lines (horizontal)
	for alt := -90.0; alt <= 90.0; alt += 15.0 {
		// Sample points along this altitude circle
		for az := 0.0; az < 360.0; az += 5.0 {
			x, y, visible := Project(alt, az, centerAlt, centerAz, fov, canvas.Width, canvas.Height)
			if visible {
				cell := canvas.Cells[y][x]
				if cell.Char == ' ' {
					canvas.Set(x, y, '·', gridStyle)
				}
			}
		}

		// Label altitude lines
		az := centerAz
		x, y, visible := Project(alt, az, centerAlt, centerAz, fov, canvas.Width, canvas.Height)
		if visible {
			label := fmt.Sprintf("%.0f°", alt)
			for i, ch := range label {
				if x+i < canvas.Width {
					canvas.Set(x+i, y, ch, labelStyle)
				}
			}
		}
	}

	// Draw azimuth lines (vertical from horizon to zenith)
	for az := 0.0; az < 360.0; az += 15.0 {
		// Sample points along this azimuth meridian
		for alt := -90.0; alt <= 90.0; alt += 5.0 {
			x, y, visible := Project(alt, az, centerAlt, centerAz, fov, canvas.Width, canvas.Height)
			if visible {
				cell := canvas.Cells[y][x]
				if cell.Char == ' ' {
					canvas.Set(x, y, '·', gridStyle)
				}
			}
		}

		// Label azimuth lines at horizon
		alt := 0.0
		x, y, visible := Project(alt, az, centerAlt, centerAz, fov, canvas.Width, canvas.Height)
		if visible {
			label := fmt.Sprintf("%.0f°", az)
			for i, ch := range label {
				if y+i < canvas.Height {
					canvas.Set(x, y+i, ch, labelStyle)
				}
			}
		}
	}

	// Draw cardinal directions at horizon
	cardinals := map[float64]string{
		0.0:   "N",
		90.0:  "E",
		180.0: "S",
		270.0: "W",
	}

	cardinalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("yellow")).
		Bold(true)

	for az, label := range cardinals {
		x, y, visible := Project(0.0, az, centerAlt, centerAz, fov, canvas.Width, canvas.Height)
		if visible {
			canvas.Set(x, y, rune(label[0]), cardinalStyle)
		}
	}
}
