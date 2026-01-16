package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Color scheme inspired by terminal defaults
var (
	// Colors using ANSI indices for maximum compatibility
	primaryColor = lipgloss.AdaptiveColor{Light: "2", Dark: "10"} // Green
	accentColor  = lipgloss.AdaptiveColor{Light: "1", Dark: "9"}  // Red
	dimColor     = lipgloss.AdaptiveColor{Light: "8", Dark: "7"}  // Gray/Dim
	warningColor = lipgloss.AdaptiveColor{Light: "3", Dark: "11"} // Yellow

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

	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)
)

// View renders the TUI
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var content string
	if m.state == stateSetup {
		content = m.setupView()
	} else if m.state == stateSites {
		content = m.sitesView()
	} else {
		content = m.timerView()
	}

	// Wrap in container with double border
	// We subtract border width (2) and padding (4) from width/height
	return containerStyle.
		Width(m.width - 6).
		Height(m.height - 4).
		Render(content)
}

func (m Model) sitesView() string {
	var b strings.Builder

	// Title with ASCII art
	title := m.renderTitle()
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n\n")

	// Header
	header := "Step 2: Select sites to block"
	b.WriteString(timeStyle.Render(header))
	b.WriteString("\n\n")

	// List of sites
	for i, site := range m.selectableSites {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := "[ ]"
		if site.selected {
			checked = "[x]"
		}

		itemStyle := lipgloss.NewStyle().Foreground(primaryColor).Align(lipgloss.Left).PaddingLeft(4)
		if m.cursor == i {
			itemStyle = itemStyle.Foreground(accentColor).Bold(true)
		}

		b.WriteString(itemStyle.Render(fmt.Sprintf("%s %s %s", cursor, checked, site.name)) + "\n")
	}

	// Virtual "Start" button
	startCursor := " "
	if m.cursor == len(m.selectableSites) {
		startCursor = ">"
	}
	startStyle := lipgloss.NewStyle().Foreground(warningColor).Align(lipgloss.Left).PaddingLeft(4).Bold(true)
	if m.cursor == len(m.selectableSites) {
		startStyle = startStyle.Background(primaryColor).Foreground(lipgloss.AdaptiveColor{Light: "15", Dark: "0"})
	}
	b.WriteString("\n")
	b.WriteString(startStyle.Render(fmt.Sprintf("%s [ START SESSION ]", startCursor)) + "\n\n")

	// Custom Input
	if m.showCustomInput {
		b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(m.siteInput.View()))
		b.WriteString("\n\n")
	}

	// Instructions
	instructions := "↑/↓: Navigate • SPACE: Toggle • ENTER: Select/Start • CTRL+C: Quit"
	b.WriteString(instructionStyle.Render(instructions))

	return b.String()
}

func (m Model) setupView() string {
	var b strings.Builder

	// Title with ASCII art
	title := m.renderTitle()
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n\n")

	// Header
	header := "Step 1: How long do you want to get serious?"
	b.WriteString(timeStyle.Render(header))
	b.WriteString("\n\n")

	// Input
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(m.input.View()))
	b.WriteString("\n\n")

	// Instructions
	instructions := "Press ENTER to continue • Press CTRL+C to quit"
	b.WriteString(instructionStyle.Render(instructions))

	return b.String()
}

func (m Model) timerView() string {
	var b strings.Builder

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

var bigDigits = map[rune][]string{
	'0': {
		"██████",
		"██  ██",
		"██  ██",
		"██  ██",
		"██████",
	},
	'1': {
		"    ██",
		"    ██",
		"    ██",
		"    ██",
		"    ██",
	},
	'2': {
		"██████",
		"    ██",
		"██████",
		"██    ",
		"██████",
	},
	'3': {
		"██████",
		"    ██",
		"██████",
		"    ██",
		"██████",
	},
	'4': {
		"██  ██",
		"██  ██",
		"██████",
		"    ██",
		"    ██",
	},
	'5': {
		"██████",
		"██    ",
		"██████",
		"    ██",
		"██████",
	},
	'6': {
		"██████",
		"██    ",
		"██████",
		"██  ██",
		"██████",
	},
	'7': {
		"██████",
		"    ██",
		"    ██",
		"    ██",
		"    ██",
	},
	'8': {
		"██████",
		"██  ██",
		"██████",
		"██  ██",
		"██████",
	},
	'9': {
		"██████",
		"██  ██",
		"██████",
		"    ██",
		"██████",
	},
	':': {
		"      ",
		"  ██  ",
		"      ",
		"  ██  ",
		"      ",
	},
	' ': {
		"      ",
		"      ",
		"      ",
		"      ",
		"      ",
	},
}

func (m Model) renderTime() string {
	hours := int(m.remaining.Hours())
	minutes := int(m.remaining.Minutes()) % 60
	seconds := int(m.remaining.Seconds()) % 60

	// Format time string
	timeStr := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

	// Render large block digits
	bigTime := m.renderBigDigits(timeStr)

	style := timeStyle
	if m.locked {
		style = lockedTimeStyle
	}

	return style.Render(bigTime)
}

func (m Model) renderBigDigits(timeStr string) string {
	height := 5
	var result []string
	for i := 0; i < height; i++ {
		var line strings.Builder
		for _, char := range timeStr {
			if patterns, ok := bigDigits[char]; ok {
				line.WriteString(patterns[i])
				line.WriteString("  ") // Spacing between digits
			}
		}
		result = append(result, line.String())
	}

	// Join lines and center the whole block
	block := strings.Join(result, "\n")
	return lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(block)
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
