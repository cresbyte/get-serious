package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Color scheme inspired by peaclock - bold, serious, terminal aesthetic
var (
	// Colors
	primaryColor = lipgloss.Color("#00FF00") // Matrix green
	accentColor  = lipgloss.Color("#FF0000") // Alert red
	dimColor     = lipgloss.Color("#666666") // Dim gray
	bgColor      = lipgloss.Color("#000000") // Pure black
	warningColor = lipgloss.Color("#FFFF00") // Warning yellow

	// Styles
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Align(lipgloss.Center)

	timeStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Align(lipgloss.Center)

	lockedTimeStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Align(lipgloss.Center)

	instructionStyle = lipgloss.NewStyle().
				Foreground(dimColor).
				Align(lipgloss.Center)

	statusStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true).
			Align(lipgloss.Center)
)

// View renders the TUI
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Add some top padding
	b.WriteString("\n\n\n")

	// Title with ASCII art
	title := m.renderTitle()
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n\n")

	// Status indicator
	status := m.renderStatus()
	b.WriteString(statusStyle.Render(status))
	b.WriteString("\n\n")

	// Main time display
	timeDisplay := m.renderTime()
	b.WriteString(timeDisplay)
	b.WriteString("\n\n")

	// Progress bar
	progressBar := m.renderProgressBar()
	b.WriteString(progressBar)
	b.WriteString("\n\n")

	// Instructions
	instructions := m.renderInstructions()
	b.WriteString(instructionStyle.Render(instructions))

	return b.String()
}

func (m Model) renderTitle() string {
	// ASCII art title
	return `
 ██████  ███████ ████████     ███████ ███████ ██████  ██  ██████  ██    ██ ███████ 
██       ██         ██        ██      ██      ██   ██ ██ ██    ██ ██    ██ ██      
██   ███ █████      ██        ███████ █████   ██████  ██ ██    ██ ██    ██ ███████ 
██    ██ ██         ██             ██ ██      ██   ██ ██ ██    ██ ██    ██      ██ 
 ██████  ███████    ██        ███████ ███████ ██   ██ ██  ██████   ██████  ███████ 
`
}

func (m Model) renderStatus() string {
	if m.locked {
		// Blinking effect using frame counter
		if m.frame%2 == 0 {
			return "🔒 LOCKED - NO ESCAPE 🔒"
		}
		return "▓▓ LOCKED - NO ESCAPE ▓▓"
	}
	return "⚡ READY TO LOCK ⚡"
}

func (m Model) renderTime() string {
	hours := int(m.remaining.Hours())
	minutes := int(m.remaining.Minutes()) % 60
	seconds := int(m.remaining.Seconds()) % 60

	// Large digital-style time display
	timeStr := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

	// Make it HUGE
	bigTime := m.renderBigDigits(timeStr)

	if m.locked {
		return lockedTimeStyle.Render(bigTime)
	}
	return timeStyle.Render(bigTime)
}

func (m Model) renderBigDigits(timeStr string) string {
	// Simple large digit rendering
	// In a full implementation, you'd have ASCII art for each digit
	return fmt.Sprintf(`
    ╔═══════════════════════════╗
    ║                           ║
    ║       %s       ║
    ║                           ║
    ╚═══════════════════════════╝
`, timeStr)
}

func (m Model) renderProgressBar() string {
	if !m.locked {
		return ""
	}

	width := 50
	progress := 1.0 - (float64(m.remaining) / float64(m.duration))
	filled := int(progress * float64(width))

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)

	percentage := int(progress * 100)

	barStyle := lipgloss.NewStyle().
		Foreground(accentColor).
		Align(lipgloss.Center)

	return barStyle.Render(fmt.Sprintf("[%s] %d%%", bar, percentage))
}

func (m Model) renderInstructions() string {
	if m.locked {
		return "Lock is active. Timer must reach 00:00:00 to unlock."
	}
	return "Press ENTER or SPACE to start lock • Press Q to quit"
}
