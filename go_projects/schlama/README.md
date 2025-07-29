# Schlama

A better Ollama user interface - A simple and elegant CLI tool for interacting with local Ollama models.

## Overview

Schlama is a command-line interface that provides an easier way to chat with local large language models through the Ollama API. It offers a streamlined experience for model management and interaction, making it simple to switch between models, send prompts, and manage your local AI setup.

## Features

- 🤖 **Easy Model Management**: List, select, and automatically pull models
- 💬 **Simple Chat Interface**: Send prompts and get formatted responses
- 📋 **Model Discovery**: Browse available models from Ollama's registry
- ⚙️ **Configuration Management**: Persistent settings for your preferred models
- 🎨 **Markdown Rendering**: Beautiful formatted output for model responses
- 🔄 **Auto-Pull**: Automatically download models when selected if not present locally

## Prerequisites

- [Ollama](https://ollama.com/) must be installed and running on your system
- Go 1.24.5 or later (for building from source)

## Installation

### From Source

1. Clone the repository:

```bash
git clone https://github.com/HanmaDevin/schlama.git
cd schlama
```

2. Build the application:

```bash
# For Unix/Linux/macOS
make build

# For Windows
make build_win
```

3. Install (optional):

```bash
make install
```

### Pre-built Binaries

Download the latest release from the [releases page](https://github.com/HanmaDevin/schlama/releases) and place the binary in your PATH.

## Usage

### Getting Started

First, make sure Ollama is running on your system. Then you can start using Schlama:

```bash
# Show help
schlama --help

# List available models for download
schlama list

# List locally installed models
schlama local

# Select a model to use (will auto-download if not present)
schlama select llama2

# Check currently selected model
schlama model

# Send a prompt to the selected model
schlama prompt "Explain quantum computing in simple terms"
```

### Commands

| Command | Description | Example |
|---------|-------------|---------|
| `list` | List available models for download | `schlama list --limit 20` |
| `local` | Show locally installed models | `schlama local` |
| `select <model>` | Select a model to use | `schlama select llama2:7b` |
| `model` | Show currently selected model | `schlama model` |
| `prompt <text>` | Send a prompt to the selected model | `schlama prompt "Hello, world!"` |

### Flags

- `--limit, -l`: Limit the number of results when listing models (default: 50)

## Configuration

Schlama automatically creates a configuration file at `~/.config/schlama/config.yaml` with the following structure:

```yaml
prompt: "What is the meaning of life?"
model: ""
stream: false
```

The configuration stores:

- **prompt**: Default prompt (currently not used in CLI)
- **model**: Currently selected model
- **stream**: Streaming mode (for future features)

## Examples

### Basic Workflow

1. **List available models:**

```bash
schlama list
```

2. **Select a model:**

```bash
schlama select phi3
```

3. **Chat with the model:**

```bash
schlama prompt "Write a haiku about programming"
```

### Advanced Usage

**Limit model list output:**

```bash
schlama list --limit 10
```

**Check what model you're using:**

```bash
schlama model
```

**View local models:**

```bash
schlama local
```

## Development

### Building

```bash
# Build for current platform
go build -o bin/schlama .

# Or use Makefile
make build      # Unix/Linux/macOS
make build_win  # Windows
```

### Running

```bash
# Run directly
go run .

# Or use built binary
./bin/schlama        # Unix/Linux/macOS
./bin/schlama.exe    # Windows
```

### Project Structure

```text
├── cmd/                 # CLI commands
│   ├── root.go         # Root command and initialization
│   ├── list.go         # List available models
│   ├── local.go        # Show local models  
│   ├── select.go       # Select model
│   ├── model.go        # Show current model
│   └── prompt.go       # Send prompts
├── config/             # Configuration management
│   └── config.go
├── ollama/             # Ollama API integration
│   └── ollama.go
├── styles/             # UI styling
│   └── styles.go
└── main.go            # Entry point
```

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Glamour](https://github.com/charmbracelet/glamour) - Markdown rendering
- [golang.org/x/net](https://golang.org/x/net) - HTTP utilities
- [yaml.v3](https://gopkg.in/yaml.v3) - YAML configuration

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the terms specified in the [LICENSE](LICENSE) file.

## Troubleshooting

### Common Issues

#### "No model specified in config"

- Run `schlama select <model_name>` to select a model first

#### "Model not found locally"

- The tool will automatically try to download the model
- Make sure you have internet connectivity and Ollama is running

#### "Connection refused"

- Ensure Ollama is running (`ollama serve`)
- Check if Ollama is running on the default port (11434)

### Getting Help

If you encounter issues:

1. Check that Ollama is properly installed and running
2. Verify your internet connection for model downloads
3. Check the [issues page](https://github.com/HanmaDevin/schlama/issues) for known problems
4. Create a new issue if your problem isn't covered

## Acknowledgments

- Thanks to the [Ollama](https://ollama.com/) team for creating an excellent tool for running local LLMs
- Built with [Cobra](https://cobra.dev/) CLI framework
- Markdown rendering powered by [Glamour](https://github.com/charmbracelet/glamour)
