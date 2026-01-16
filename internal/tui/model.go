package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// tickMsg will be sent on every timer tick hopefully
type tickMsg time.Time

type Model struct {
	duration  time.Duration
	remaining time.Duration
	startTime time.Time
	locked    bool
	quitting  bool
	width     int
	height    int
	frame     int
}

func NewModel() Model {
	// I want the default to be 25 minutes - some Pomodoro nonsense
	duration := 1 * time.Minute

	return Model{
		duration:  duration,
		remaining: duration,
		locked:    false,
		frame:     0,
	}
}

// This is supposed to return a command that sends a tick message every second
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Initialize the model I think its required by BB tea
func (m Model) Init() tea.Cmd {
	// tea.Batch does not guarantee results, that f
	return tea.Batch(tea.EnterAltScreen, tickCmd())
}
