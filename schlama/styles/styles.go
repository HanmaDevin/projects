package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Catppuccin Mocha palette
	Rosewater = lipgloss.Color("#f5e0dc")
	Mauve     = lipgloss.Color("#cba6f7")
	Text      = lipgloss.Color("#cdd6f4")
	Surface0  = lipgloss.Color("#313244")
	Surface1  = lipgloss.Color("#45475a")
	Green     = lipgloss.Color("#a6e3a1")
	Maroon    = lipgloss.Color("#f38ba8")
	Sky       = lipgloss.Color("#89dceb")
	Blue      = lipgloss.Color("#89b4fa")
	Subtext0  = lipgloss.Color("#a6adc8")

	OutputStyle = lipgloss.NewStyle().
			Background(Surface0).
			Foreground(Text).
			Bold(true).
			Padding(0, 1).Render

	ErrorStyle = lipgloss.NewStyle().
			Background(Surface0).
			Foreground(Maroon).
			Bold(true).
			Padding(0, 1).Render

	FinishedStyle = lipgloss.NewStyle().
			Background(Surface0).
			Foreground(Green).
			Bold(true).
			Italic(true).
			Padding(0, 1).Render

	HintStyle = lipgloss.NewStyle().
			Foreground(Sky).
			Background(Surface0).
			Italic(true).
			Padding(0, 1).Render

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Mauve).
			Background(Surface1).
			Padding(0, 1).Render

	UserStyle = lipgloss.NewStyle().
			Foreground(Blue).Render

	AiStyle = lipgloss.NewStyle().
		Foreground(Green).Render

	RowStyle = lipgloss.NewStyle().
			Foreground(Text).
			Background(Surface0).
			Padding(0, 1).Render

	TuiBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Rosewater).
			Padding(0, 1).Render

	ChatView = lipgloss.NewStyle().
			Foreground(Text).
			Render

	HelpView = lipgloss.NewStyle().
			Foreground(Subtext0).
			Render
)
