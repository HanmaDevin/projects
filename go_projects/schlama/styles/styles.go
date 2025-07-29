package styles

import "github.com/charmbracelet/lipgloss"

var (
	OutputStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#7287fd")).Bold(true)
	MsgStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#04a5e5")).Faint(true)
	ErrorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#d20f39")).Bold(true)
	HighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#40a02b")).Bold(true).Italic(true)
)
