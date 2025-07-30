package tui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/HanmaDevin/schlama/config"
	"github.com/HanmaDevin/schlama/ollama"
	"github.com/HanmaDevin/schlama/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	model        string
	chat         []string
	textinput    textinput.Model
	help         string
	width        int
	height       int
	scrollOffset int // Track scroll position
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
		model:        body.Model,
		chat:         []string{},
		textinput:    ti,
		help:         "Enter: Send | Ctrl+C: Quit | /help: Show help | ↑↓: Scroll",
		width:        80, // Default width
		height:       24, // Default height
		scrollOffset: 0,  // Start at bottom (most recent)
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
		return "[help] Available commands:\n/select <model> - Select a model\n/list - List all models\n/local - List local models"
	case "/select":
		if len(parts) < 2 {
			return listSelectableModels()
		}
		modelName := strings.Join(parts[1:], " ")
		return fmt.Sprintf("[select] Selected model: %s", modelName)
	case "/list":
		return listAvailableModels()
	case "/local":
		localOutput := getLocalModelsForDisplay()
		if localOutput == "No models found" {
			return "[local] No local models found.\nUse 'ollama pull <model>' to install models first."
		}
		if localOutput == "Could not retrieve local models" {
			return "[local] Error: Could not retrieve local models.\nMake sure Ollama is running."
		}
		return "[local] Installed models:\n\n" + localOutput
	default:
		return fmt.Sprintf("Unknown command: %s", command)
	}
}

func listSelectableModels() string {
	// Get local models (already installed)
	localModels := getLocalModelNames()

	if len(localModels) == 0 {
		return "[select] No local models found.\nUse 'ollama pull <model>' to install models first.\n\nAvailable models to pull:\n" + getPopularModels()
	}

	var result strings.Builder
	result.WriteString("[select] Choose a model:\n\n")
	result.WriteString("Local models (ready to use):\n")

	for i, model := range localModels {
		result.WriteString(fmt.Sprintf("  %d. %s\n", i+1, model))
	}

	result.WriteString("\nUsage: /select <model_name>\n")
	result.WriteString("Example: /select " + localModels[0])

	return result.String()
}

func getLocalModelNames() []string {
	cmd := exec.Command("ollama", "list")
	out, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return []string{}
	}

	var modelNames []string
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			// Extract just the model name (first field)
			modelNames = append(modelNames, fields[0])
		}
	}

	return modelNames
}

func getLocalModelsForDisplay() string {
	cmd := exec.Command("ollama", "list")
	out, err := cmd.Output()
	if err != nil {
		return "Could not retrieve local models"
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return "No models found"
	}

	var result strings.Builder

	// Just list the model names, one per line
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) > 0 {
			result.WriteString("• " + fields[0] + "\n")
		}
	}

	return strings.TrimRight(result.String(), "\n")
}

func getPopularModels() string {
	popularModels := []string{
		"llama3.2:latest",
		"qwen2.5:latest",
		"deepseek-r1:latest",
		"codellama:latest",
		"mistral:latest",
	}

	var result strings.Builder
	for _, model := range popularModels {
		result.WriteString("  • " + model + "\n")
	}

	return result.String()
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
		case "up":
			// Scroll up (show older messages)
			if m.scrollOffset < len(m.chat)-1 {
				m.scrollOffset++
			}
			return m, nil
		case "down":
			// Scroll down (show newer messages)
			if m.scrollOffset > 0 {
				m.scrollOffset--
			}
			return m, nil
		case "pgup":
			// Page up (scroll up by 5 lines)
			m.scrollOffset = min(m.scrollOffset+5, len(m.chat)-1)
			return m, nil
		case "pgdown":
			// Page down (scroll down by 5 lines)
			m.scrollOffset = max(m.scrollOffset-5, 0)
			return m, nil
		case "enter":
			input := m.textinput.Value()
			if strings.HasPrefix(input, "/") {
				result := handleCommand(input)

				// Clear chat and show command result
				m.chat = []string{result}
				m.scrollOffset = 0 // Reset scroll position

				// Handle model selection
				if strings.HasPrefix(input, "/select ") {
					modelName := strings.TrimSpace(strings.TrimPrefix(input, "/select "))
					if modelName != "" {
						// Check if the model is present in the local models (exactly like select.go)
						if !ollama.IsModelPresent(modelName) {
							m.chat = []string{fmt.Sprintf("[select] Model not found locally. Pulling model: %s\nThis may take a while depending on the model size...", modelName)}

							// Attempt to pull the model
							err := ollama.PullModel(modelName)
							if err != nil {
								m.chat = []string{fmt.Sprintf("[select] Error pulling model: %s\n\nHere is a list of available local models:\n%s", err.Error(), getLocalModelsForDisplay())}
								return m, nil
							}
							// Model was successfully pulled, continue with selection
						}

						// At this point, model should exist locally, so select it
						body := config.ReadConfig()
						if body == nil {
							body = ollama.NewOllamaModel()
						}

						cfg := config.Config{
							Model: modelName,
						}

						if err := config.WriteConfig(cfg); err != nil {
							m.chat = []string{fmt.Sprintf("[select] Error saving configuration: %s", err)}
						} else {
							m.model = modelName
							m.chat = []string{fmt.Sprintf("[select] Current Model: %s", modelName)}
						}
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
				m.scrollOffset = 0 // Reset to bottom when new message arrives
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

	// If we have content, fill from the top with scroll offset
	if len(allLines) > 0 {
		// Calculate the start index based on scroll offset
		startIdx := 0
		if len(allLines) > chatHeight {
			// Normal case: more content than can fit
			if m.scrollOffset == 0 {
				// At bottom (most recent)
				startIdx = len(allLines) - chatHeight
			} else {
				// Scrolled up from bottom
				startIdx = len(allLines) - chatHeight - m.scrollOffset
				if startIdx < 0 {
					startIdx = 0
				}
			}
		}

		// Copy lines to display, starting from the calculated index
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

	// Add scroll indicator
	scrollInfo := ""
	if len(allLines) > chatHeight {
		if m.scrollOffset == 0 {
			scrollInfo = " (Latest)"
		} else {
			scrollInfo = fmt.Sprintf(" (↑%d)", m.scrollOffset)
		}
	}

	chatModel := styles.HelpView(fmt.Sprintf("Agent: %s%s", m.model, scrollInfo))
	return fmt.Sprintf("\n%s\n%s\n\n%s\n\n%s", chatModel, chatView, inputView, helpView)
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
