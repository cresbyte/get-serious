package main

import (
	"fmt"
	"os"
	"os/exec"

	"get-serious/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if os.Getenv("GET_SERIOUS_DETACHED") != "1" {
		exe, err := os.Executable()
		if err == nil {
			// Try common terminals with small window geometry
			terminals := []struct {
				cmd  string
				args []string
			}{
				{"alacritty", []string{"--dimensions", "80", "24", "-e", exe}},
				{"xterm", []string{"-geometry", "80x24", "-e", exe}},
				{"gnome-terminal", []string{"--geometry=80x24", "--", exe}},
				{"konsole", []string{"--dimensions", "80x24", "-e", exe}},
			}

			for _, t := range terminals {
				if _, err := exec.LookPath(t.cmd); err == nil {
					cmd := exec.Command(t.cmd, t.args...)
					cmd.Env = append(os.Environ(), "GET_SERIOUS_DETACHED=1")
					if err := cmd.Start(); err == nil {
						return // Exit original process
					}
				}
			}
		}
	}

	program := tea.NewProgram(tui.NewModel(), tea.WithAltScreen(),
		tea.WithMouseCellMotion())

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running this program")
		os.Exit(1)
	}
}
