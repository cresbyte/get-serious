package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			// Allow quitting only if not locked
			if !m.locked {
				m.quitting = true
				return m, tea.Quit
			}
			// If locked, ignore quit attempts (for now)
			return m, nil
			
		case "enter", " ":
			// Start the lock
			if !m.locked {
				m.locked = true
				m.startTime = time.Now()
				m.remaining = m.duration
			}
			return m, nil
		}
	
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	
	case tickMsg:
		// Increment animation frame
		m.frame++
		
		// Update remaining time if locked
		if m.locked {
			elapsed := time.Since(m.startTime)
			m.remaining = m.duration - elapsed
			
			// Check if time is up
			if m.remaining <= 0 {
				m.remaining = 0
				m.locked = false
				// Could add a completion message here
			}
		}
		
		// Schedule next tick
		return m, tickCmd()
	}
	
	return m, nil
}
