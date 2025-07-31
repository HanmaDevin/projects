# 🦙 Schlama

> **Because talking to llamas should be this easy!** 🚀

A delightfully simple CLI and TUI for chatting with your local Ollama models. No more wrestling with complex commands - just pure llama magic! ✨

## 🎯 What Does It Do?

Schlama makes your local AI models as easy to use as ordering coffee:

- 🗨️ **Interactive TUI Chat** - Like WhatsApp, but for llamas
- ⚡ **Lightning CLI** - One-liners for quick AI tasks  
- 🔄 **Smart Model Management** - Auto-downloads what you need
- 🎨 **Beautiful Output** - Markdown rendering that doesn't hurt your eyes

## 🚀 Quick Start

```bash
# Install Ollama first: https://ollama.com
# Then build Schlama:
git clone https://github.com/HanmaDevin/schlama.git
cd schlama && make build

# Start the interactive chat
./bin/schlama tui

# Or use CLI mode
./bin/schlama select llama3.2
./bin/schlama prompt "What's the meaning of life?"
```

## 🎮 Two Ways to Play

### 🖥️ TUI Mode (Recommended)

```bash
schlama tui
```

- **Interactive chat interface** with scrolling
- **Built-in commands** (`/help`, `/select`, `/local`, `/list`)
- **Dynamic sizing** - adapts to your terminal
- **Persistent conversations** - until you exit

### ⚡ CLI Mode (Power Users)

```bash
schlama list              # Browse all available models
schlama local             # See what's installed  
schlama select phi3       # Pick your model (auto-downloads!)
schlama model             # Check current selection
schlama prompt "Hi AI!"   # Chat away
```

## 🎨 TUI Commands

| Command | What It Does | Example |
|---------|--------------|---------|
| `/help` | Shows available commands | `/help` |
| `/select model` | Switch models (pulls if needed) | `/select llama3.2` |
| `/local` | List installed models | `/local` |
| `/list` | Browse all available models | `/list` |
| **Just type** | Chat with current model | `Write me a poem` |

**🎯 Pro Tips:**

- Use `↑↓` arrows to scroll through chat history
- `Ctrl+C` to exit
- Models auto-download when selected - grab a coffee! ☕

## 📁 Project Structure

```text
schlama/
├── 🎯 cmd/          # CLI commands
├── ⚙️  config/      # Settings management  
├── 🦙 ollama/       # API magic
├── 🎨 styles/       # Pretty colors
├── 🖥️  tui/         # Interactive interface
└── 📦 bin/          # Your built binary
```

## 🔧 Configuration

Lives at `~/.config/schlama/config.yaml`:

```yaml
model: "llama3.2"    # Your currently selected model
```

> **Note:** The config only stores the selected model. Prompts are handled throughout the app. Stream is always off.

## ⚠️ Troubleshooting

### "Ollama is not running"

```bash
ollama serve  # Start Ollama first!
```

### "No models found"

```bash
schlama select llama3.2  # Downloads automatically
```

### TUI looks weird?

- Resize your terminal
- Try a different terminal emulator

## 🛠️ Building & Development

```bash
# Build for your platform
make build          # Unix/Linux/macOS  
make build_win      # Windows

# Quick development
go run . tui        # Test TUI
go run . prompt "test"  # Test CLI
```

## 🎉 Features

- ✅ **Interactive TUI** with real-time chat
- ✅ **CLI mode** for automation & scripting
- ✅ **Auto-model downloading**
- ✅ **Beautiful Catppuccin theme**
- ✅ **Markdown rendering**
- ✅ **Scrollable chat history**
- ✅ **Model management**
- ✅ **Configuration persistence**

## 🤝 Contributing

Got ideas? Found bugs? Want to make llamas even more awesome?

1. Fork it 🍴
2. Branch it (`git checkout -b feature/llama-superpowers`)
3. Commit it (`git commit -m 'Add llama telepathy'`)
4. Push it (`git push origin feature/llama-superpowers`)
5. PR it! 🚀

## 📜 License

See [LICENSE](LICENSE) file. TL;DR: Be nice, have fun! 🎈

## 🙏 Thanks

- **[Ollama](https://ollama.com/)** - For making local AI accessible
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - For the amazing TUI framework  
- **[Cobra](https://cobra.dev/)** - For CLI superpowers
- **You!** - For giving Schlama a try 🎉

---

**Happy llamaing!** 🦙✨
