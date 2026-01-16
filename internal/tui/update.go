package tui

import (
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.state {
	case stateSetup:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit
			case "enter":
				val := m.input.Value()
				if val == "" {
					return m, nil
				}

				// Try parsing as minutes if it's just a number
				if _, err := strconv.Atoi(val); err == nil {
					val += "m"
				}

				d, err := time.ParseDuration(val)
				if err != nil {
					m.input.SetValue("")
					return m, nil
				}
				m.duration = d
				m.remaining = d
				m.state = stateSites
				m.siteInput.Focus()
				return m, nil
			}
		}

		m.input, cmd = m.input.Update(msg)
		return m, cmd

	case stateSites:
		if m.showCustomInput {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.String() {
				case "esc":
					m.showCustomInput = false
					return m, nil
				case "enter":
					val := strings.TrimSpace(m.siteInput.Value())
					if val != "" {
						// Add the custom site to the list (just before "Add Custom...")
						newSite := selectableSite{name: val, selected: true}
						m.selectableSites = append(m.selectableSites[:len(m.selectableSites)-1], newSite, m.selectableSites[len(m.selectableSites)-1])
						m.siteInput.SetValue("")
						m.showCustomInput = false
					}
					return m, nil
				}
			}
			m.siteInput, cmd = m.siteInput.Update(msg)
			return m, cmd
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.selectableSites) {
					m.cursor++
				}
			case " ":
				if m.cursor < len(m.selectableSites) {
					m.selectableSites[m.cursor].selected = !m.selectableSites[m.cursor].selected
				}
			case "enter":
				if m.cursor == len(m.selectableSites) {
					// "Start Session" button (virtual)
					return m, tickCmd()
				}
				if m.cursor < len(m.selectableSites) {
					if m.selectableSites[m.cursor].name == "Add Custom..." {
						m.showCustomInput = true
						m.siteInput.Focus()
					} else {
						// Optionally transition if they press enter on "Start"
						m.state = stateTimer
						return m, tickCmd()
					}
				}
			}
		}
		return m, nil

	case stateTimer:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "ctrl+c", "esc":
				if !m.locked {
					m.quitting = true
					return m, tea.Quit
				}
				return m, nil

			case "enter", " ":
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
			m.frame++
			if m.locked {
				elapsed := time.Since(m.startTime)
				m.remaining = m.duration - elapsed
				if m.remaining <= 0 {
					m.remaining = 0
					m.locked = false
					m.state = stateSetup
					m.input.SetValue("")
					m.input.Focus()
					m.cursor = 0
					return m, nil
				}
			}
			return m, tickCmd()
		}
	}

	return m, nil
}
