// main.go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	config, err := LoadConfig("config.toml")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	InitStyles(config)

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
