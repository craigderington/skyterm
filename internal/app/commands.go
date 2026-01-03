package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TickMsg is sent on every tick to update the time
type TickMsg time.Time

// tickCmd returns a command that sends a tick message
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
