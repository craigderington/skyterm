package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// RenderHelp returns a help screen with keybindings
func RenderHelp(width, height int) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("cyan")).
		Bold(true)

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("yellow")).
		Bold(true)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("green")).
		Width(18)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Width(40)

	// Helper function to create a help line with aligned columns
	line := func(key, desc string) string {
		return "  " + keyStyle.Render(key) + " " + descStyle.Render(desc)
	}

	var help string
	help += titleStyle.Render("skyterm - Terminal Astronomy") + "\n\n"

	help += sectionStyle.Render("Navigation") + "\n"
	help += line("↑/k, ↓/j", "Pan up/down") + "\n"
	help += line("←/h, →/l", "Pan left/right") + "\n"
	help += line("K/J/H/L", "Fast pan") + "\n"
	help += line("+/-", "Zoom in/out") + "\n"
	help += line("0", "Reset view") + "\n\n"

	help += sectionStyle.Render("Cardinal Directions") + "\n"
	help += line("n/s/e/w", "North/South/East/West") + "\n"
	help += line("z", "Zenith (straight up)") + "\n\n"

	help += sectionStyle.Render("Display Toggles") + "\n"
	help += line("g", "Toggle coordinate grid") + "\n"
	help += line("C", "Toggle constellation lines") + "\n"
	help += line("N", "Toggle constellation names") + "\n"
	help += line("p", "Toggle planets (Sun, Moon, planets)") + "\n"
	help += line("P", "Toggle planet labels") + "\n"
	help += line("d", "Toggle deep sky objects (Messier)") + "\n"
	help += line("S", "Toggle star labels (bright stars)") + "\n"
	help += line("m", "Cycle magnitude limit (3/4/5/6)") + "\n\n"

	help += sectionStyle.Render("Object Interaction") + "\n"
	help += line("Enter", "Select nearest object to center") + "\n"
	help += line("i", "Toggle info panel for selected") + "\n"
	help += line("c", "Center view on selected object") + "\n"
	help += line("f", "Follow selected object (lock view)") + "\n"
	help += line("/", "Search for object by name") + "\n\n"

	help += sectionStyle.Render("Time Controls") + "\n"
	help += line("Space", "Pause/resume time flow") + "\n"
	help += line("[ / ]", "Step time backward/forward") + "\n"
	help += line("{ / }", "Fast step time backward/forward") + "\n"
	help += line("T", "Jump to current time (now)") + "\n"
	help += line("t", "Set custom time") + "\n\n"

	help += sectionStyle.Render("General") + "\n"
	help += line("?", "Toggle this help screen") + "\n"
	help += line("q, Ctrl+C", "Quit application") + "\n\n\n"

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Faint(true)
	help += footerStyle.Render("Stars update in real-time based on your location and current time.")

	// Center the help text
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("cyan"))

	return style.Render(help)
}
