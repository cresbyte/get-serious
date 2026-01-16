package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// tickMsg will be sent on every timer tick hopefully
type tickMsg time.Time

type sessionState int

const (
	stateSetup sessionState = iota
	stateSites
	stateTimer
)

type selectableSite struct {
	name     string
	selected bool
}

type Model struct {
	state           sessionState
	duration        time.Duration
	remaining       time.Duration
	startTime       time.Time
	locked          bool
	quitting        bool
	width           int
	height          int
	frame           int
	input           textinput.Model
	siteInput       textinput.Model
	selectableSites []selectableSite
	cursor          int
	showCustomInput bool
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "25m"
	ti.Focus()
	ti.CharLimit = 10
	ti.Width = 20

	si := textinput.New()
	si.Placeholder = "Enter custom site..."
	si.CharLimit = 50
	si.Width = 30

	defaultSites := []selectableSite{
		{name: "google.com", selected: false},
		{name: "youtube.com", selected: false},
		{name: "facebook.com", selected: false},
		{name: "instagram.com", selected: false},
		{name: "x.com", selected: false},
		{name: "reddit.com", selected: false},
		{name: "linkedin.com", selected: false},
		{name: "netflix.com", selected: false},
		{name: "Add Custom...", selected: false},
	}

	return Model{
		state:           stateSetup,
		duration:        0,
		remaining:       0,
		locked:          false,
		frame:           0,
		input:           ti,
		siteInput:       si,
		selectableSites: defaultSites,
		cursor:          0,
		showCustomInput: false,
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
