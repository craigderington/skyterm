package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// RenderTimeInput renders the time input modal
func RenderTimeInput(input string, width, height int) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("cyan")).
		Bold(true).
		Padding(0, 1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("231")).
		Bold(true)

	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Faint(true)

	// Build content
	content := titleStyle.Render("Set Time") + "\n\n"
	content += labelStyle.Render("Enter time:") + "\n"
	content += inputStyle.Render(input + "█") + "\n\n"
	content += instructionStyle.Render("Format: YYYY-MM-DD HH:MM:SS") + "\n"
	content += instructionStyle.Render("    or: YYYY-MM-DD HH:MM") + "\n"
	content += instructionStyle.Render("    or: YYYY-MM-DD") + "\n\n"
	content += instructionStyle.Render("Press Enter to confirm, Esc to cancel")

	// Create modal with border
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("51")).
		Padding(1, 2).
		Width(50)

	modal := modalStyle.Render(content)

	// Center the modal
	positioned := lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
		lipgloss.WithWhitespaceChars(" "),
	)

	return positioned
}
