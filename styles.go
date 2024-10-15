// styles.go
package main

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle        lipgloss.Style
	normalTextStyle   lipgloss.Style
	cursorStyle       lipgloss.Style
	selectedStyle     lipgloss.Style
	borderStyle       lipgloss.Style
	instructionStyle  lipgloss.Style
	activeColumnStyle lipgloss.Style
)

func InitStyles(config *Config) {
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(config.Colors.Title)).
		Padding(0, 1).
		BorderStyle(lipgloss.Border{Bottom: "─"}).
		BorderBottom(true).
		BorderForeground(lipgloss.Color(config.Colors.Border))

	normalTextStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(config.Colors.NormalText))

	cursorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(config.Colors.ActiveColumnBg)).
		Background(lipgloss.Color(config.Colors.Cursor)).
		Bold(true)

	selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(config.Colors.ActiveColumnBg)).
		Background(lipgloss.Color(config.Colors.Selected)).
		Bold(true)

	borderStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(config.Colors.Border))

	instructionStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(config.Colors.Instruction))

	activeColumnStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(config.Colors.ActiveColumnBg))
}