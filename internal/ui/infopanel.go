package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/craigderington/skyterm/internal/astro"
	"github.com/craigderington/skyterm/internal/catalog"
	"github.com/craigderington/skyterm/internal/image"
)

// ObjectInfo holds information about a selected object for display
type ObjectInfo struct {
	Type         string
	Name         string
	Star         *catalog.Star
	Planet       *astro.Planet
	DeepSky      *catalog.MessierObject
	ImageInfo    *image.WikipediaImageInfo
	ImageData    string // Rendered image for terminal
	ImageLoading bool   // True while fetching image
	ImageError   error  // Error if image fetch failed
}

// RenderInfoPanel renders an information panel for the selected object
func RenderInfoPanel(info *ObjectInfo, observer *astro.Observer, width, height int) string {
	if info == nil {
		return ""
	}

	selected := info

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("cyan")).
		Bold(true).
		Padding(0, 1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("green")).
		Width(15)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("231"))

	var content string

	// Image status indicator
	if selected.ImageLoading {
		loadingStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("yellow")).
			Italic(true)
		content += loadingStyle.Render("Loading image...") + "\n\n"
	} else if selected.ImageInfo != nil && selected.ImageData != "" {
		// Show that an image is available with viewing instruction
		infoStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("cyan")).
			Bold(true)
		viewStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("green"))
		content += infoStyle.Render("📷 Image available") + "\n"
		content += viewStyle.Render("Press 'v' to view fullscreen") + "\n\n"
	} else if selected.ImageError != nil {
		// Show error message
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("red")).
			Faint(true)
		content += errorStyle.Render("(Image unavailable)") + "\n\n"
	}

	switch selected.Type {
	case "star":
		if selected.Star == nil {
			return ""
		}
		s := selected.Star

		content += titleStyle.Render(s.Name) + "\n\n"
		content += labelStyle.Render("Type:") + valueStyle.Render("Star") + "\n"
		content += labelStyle.Render("Magnitude:") + valueStyle.Render(fmt.Sprintf("%.2f", s.Magnitude)) + "\n"
		content += labelStyle.Render("Spectral Type:") + valueStyle.Render(string(s.SpectralType)) + "\n"
		content += "\n"
		content += labelStyle.Render("RA:") + valueStyle.Render(astro.FormatRA(s.RA)) + "\n"
		content += labelStyle.Render("Dec:") + valueStyle.Render(astro.FormatDec(s.Dec)) + "\n"
		content += labelStyle.Render("Altitude:") + valueStyle.Render(fmt.Sprintf("%.1f°", s.Altitude)) + "\n"
		content += labelStyle.Render("Azimuth:") + valueStyle.Render(fmt.Sprintf("%.1f°", s.Azimuth)) + "\n"
		content += "\n"

		// Calculate rise/set/transit times
		rst := astro.CalculateRiseSetTransit(s.RA, s.Dec, observer, astro.CurrentTime())
		if rst.NeverRises {
			content += labelStyle.Render("Visibility:") + valueStyle.Render("Never rises") + "\n"
		} else if rst.Circumpolar {
			content += labelStyle.Render("Visibility:") + valueStyle.Render("Circumpolar") + "\n"
			content += labelStyle.Render("Transit:") + valueStyle.Render(rst.Transit.Format("15:04")) + "\n"
		} else {
			content += labelStyle.Render("Rises:") + valueStyle.Render(astro.FormatTime(rst.Rise)) + "\n"
			content += labelStyle.Render("Transit:") + valueStyle.Render(rst.Transit.Format("15:04")) + "\n"
			content += labelStyle.Render("Sets:") + valueStyle.Render(astro.FormatTime(rst.Set)) + "\n"
		}

	case "planet":
		if selected.Planet == nil {
			return ""
		}
		p := selected.Planet

		content += titleStyle.Render(p.Name) + "\n\n"
		content += labelStyle.Render("Type:") + valueStyle.Render(string(p.BodyType)) + "\n"
		content += labelStyle.Render("Magnitude:") + valueStyle.Render(fmt.Sprintf("%.1f", p.Magnitude)) + "\n"
		content += "\n"
		content += labelStyle.Render("RA:") + valueStyle.Render(astro.FormatRA(p.RA)) + "\n"
		content += labelStyle.Render("Dec:") + valueStyle.Render(astro.FormatDec(p.Dec)) + "\n"
		content += labelStyle.Render("Altitude:") + valueStyle.Render(fmt.Sprintf("%.1f°", p.Altitude)) + "\n"
		content += labelStyle.Render("Azimuth:") + valueStyle.Render(fmt.Sprintf("%.1f°", p.Azimuth)) + "\n"

	case "deepsky":
		if selected.DeepSky == nil {
			return ""
		}
		d := selected.DeepSky

		content += titleStyle.Render(d.Name) + "\n"
		if d.CommonName != "" {
			content += valueStyle.Render(d.CommonName) + "\n"
		}
		content += "\n"
		content += labelStyle.Render("Type:") + valueStyle.Render(d.Type) + "\n"
		content += labelStyle.Render("Magnitude:") + valueStyle.Render(fmt.Sprintf("%.1f", d.Magnitude)) + "\n"
		content += "\n"
		content += labelStyle.Render("RA:") + valueStyle.Render(astro.FormatRA(d.RA)) + "\n"
		content += labelStyle.Render("Dec:") + valueStyle.Render(astro.FormatDec(d.Dec)) + "\n"
		content += labelStyle.Render("Altitude:") + valueStyle.Render(fmt.Sprintf("%.1f°", d.Altitude)) + "\n"
		content += labelStyle.Render("Azimuth:") + valueStyle.Render(fmt.Sprintf("%.1f°", d.Azimuth)) + "\n"
	}

	// Add close instruction
	closeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Faint(true)
	content += "\n" + closeStyle.Render("Press 'i' again to close")

	// Create panel with border for text content
	panelWidth := 40

	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("51")).
		Padding(1, 2).
		Width(panelWidth)

	panel := panelStyle.Render(content)

	// Position text panel in top-right corner
	positioned := lipgloss.Place(
		width,
		height,
		lipgloss.Right,
		lipgloss.Top,
		panel,
		lipgloss.WithWhitespaceChars(" "),
	)

	return positioned
}
