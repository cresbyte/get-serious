package main

import (
	"fmt"
	"os"

	"get-serious/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	program := tea.NewProgram(tui.NewModel(), tea.WithAltScreen(),
		tea.WithMouseCellMotion())

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running this program")
		os.Exit(1)
	}
}
