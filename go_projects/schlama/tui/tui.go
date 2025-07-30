package tui

import (
	"fmt"
	"strings"

	"github.com/HanmaDevin/schlama/config"
	"github.com/HanmaDevin/schlama/ollama"
	"github.com/HanmaDevin/schlama/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	model     string
	chat      []string
	textinput textinput.Model
	help      string
	width     int
	height    int
}

func Start() error {
	p := tea.NewProgram(initialModel())
	_, err := p.Run()
	return err
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Ask a question or type a command..."
	ti.Focus()
	ti.CharLimit = 256 * 2048
	ti.Width = 40

	body := config.ReadConfig()

	return model{
		model:     body.Model,
		chat:      []string{},
		textinput: ti,
		help:      "Enter: Send | Ctrl+C: Quit | /help: Show help",
		width:     80, // Default width
		height:    24, // Default height
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func handleCommand(input string) string {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "Invalid command"
	}

	command := parts[0]

	switch command {
	case "/help":
		return "[help] Available commands:\n/list - List all available models\n/local - List local models\n/select <model> - Select a model"
	case "/list":
		return listAvailableModels()
	case "/local":
		return ollama.ListLocalModels()
	case "/select":
		if len(parts) < 2 {
			return "[select] Usage: /select <model_name>\nUse /list to see available models"
		}
		modelName := strings.Join(parts[1:], " ")
		return fmt.Sprintf("[select] Selected model: %s", modelName)
	default:
		return fmt.Sprintf("Unknown command: %s", command)
	}
}

func listAvailableModels() string {
	models := ollama.ListModels()
	if len(models) == 0 {
		return "[list] No models available"
	}

	var result strings.Builder
	result.WriteString("[list] Available models:\n")

	for i, model := range models {
		if i >= 20 { // Limit to first 20 models to avoid overwhelming the chat
			result.WriteString("... and more (total: ")
			result.WriteString(fmt.Sprintf("%d", len(models)))
			result.WriteString(" models)")
			break
		}

		result.WriteString("• ")
		result.WriteString(model.Name)

		if len(model.Sizes) > 0 {
			result.WriteString(" (")
			result.WriteString(strings.Join(model.Sizes, ", "))
			result.WriteString(")")
		}
		result.WriteString("\n")
	}

	return result.String()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			input := m.textinput.Value()
			if strings.HasPrefix(input, "/") {
				result := handleCommand(input)
				m.chat = append(m.chat, result)

				// Handle model selection
				if strings.HasPrefix(input, "/select ") {
					modelName := strings.TrimSpace(strings.TrimPrefix(input, "/select"))
					if modelName != "" {
						m.model = modelName
						// Update config if needed
						body := config.ReadConfig()
						body.Model = modelName
						// Note: You might want to save this back to config
					}
				}
			} else {
				body := config.ReadConfig()
				body.Prompt = input
				response, err := ollama.GetResponse(body)
				if err != nil {
					m.chat = append(m.chat, fmt.Sprintf("Error: %s", err))
				} else {
					m.chat = append(m.chat, fmt.Sprintf("You: %s\nAI: %s", input, response))
				}
			}
			m.textinput.SetValue("")
		}
	}
	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	// Calculate chat dimensions based on terminal size
	chatWidth := min(m.width-6, 80) // Leave some margin, cap at 80
	if chatWidth < 40 {
		chatWidth = 40 // Minimum width
	}

	// Update text input width to match chat width
	m.textinput.Width = chatWidth - 4 // Account for border padding

	// Calculate available height for chat (subtract space for input, help, model info, and borders)
	reservedHeight := 8 // Space for input box, help text, model info, and borders
	availableHeight := m.height - reservedHeight

	minChatHeight := 10
	maxChatHeight := max(availableHeight, minChatHeight)

	// Calculate how many lines we actually need
	var allLines []string
	for _, message := range m.chat {
		wrappedLines := wrapText(message, chatWidth)
		allLines = append(allLines, wrappedLines...)
	}

	// Determine the optimal chat height based on content and available space
	chatHeight := minChatHeight
	if len(allLines) > minChatHeight && len(allLines) <= maxChatHeight {
		chatHeight = len(allLines)
	} else if len(allLines) > maxChatHeight {
		chatHeight = maxChatHeight
	}

	// Create a slice to hold the final display lines
	displayLines := make([]string, chatHeight)

	// Fill with empty strings initially
	for i := range displayLines {
		displayLines[i] = strings.Repeat(" ", chatWidth)
	}

	// If we have content, fill from the top
	if len(allLines) > 0 {
		// Take the most recent lines that fit in our display
		startIdx := 0
		if len(allLines) > chatHeight {
			startIdx = len(allLines) - chatHeight
		}

		// Copy lines to display, starting from the top
		for i, line := range allLines[startIdx:] {
			if i < chatHeight {
				// Pad the line to exactly chatWidth
				if len(line) < chatWidth {
					displayLines[i] = line + strings.Repeat(" ", chatWidth-len(line))
				} else {
					displayLines[i] = line[:chatWidth]
				}
			}
		}
	}

	chatContent := strings.Join(displayLines, "\n")
	chatView := styles.TuiBorder(chatContent)
	inputView := styles.TuiBorder(m.textinput.View())
	helpView := styles.HelpView(m.help)
	chatModel := styles.HelpView(fmt.Sprintf("Agent: %s", m.model))
	return fmt.Sprintf("\n%s\n%s\n\n%s\n\n%s", chatModel, chatView, inputView, helpView)
}

// Helper functions for min/max
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// wrapText wraps text to fit within the specified width, preserving newlines
func wrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	// First split by newlines to preserve line breaks
	textLines := strings.Split(text, "\n")
	var allLines []string

	for _, line := range textLines {
		if strings.TrimSpace(line) == "" {
			// Add empty line for blank lines
			allLines = append(allLines, "")
			continue
		}

		words := strings.Fields(line)
		if len(words) == 0 {
			allLines = append(allLines, "")
			continue
		}

		var currentLine strings.Builder

		for _, word := range words {
			// If adding this word would exceed the width, start a new line
			if currentLine.Len() > 0 && currentLine.Len()+1+len(word) > width {
				allLines = append(allLines, currentLine.String())
				currentLine.Reset()
			}

			// Add word to current line
			if currentLine.Len() > 0 {
				currentLine.WriteString(" ")
			}
			currentLine.WriteString(word)
		}

		// Add the last line if it has content
		if currentLine.Len() > 0 {
			allLines = append(allLines, currentLine.String())
		}
	}

	// Handle case where we have no lines
	if len(allLines) == 0 {
		allLines = []string{""}
	}

	return allLines
}
