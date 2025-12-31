# ShellNotes

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

**A lightning-fast, terminal-based note-taking application for developers**

[Features](#features) • [Installation](#installation) • [Usage](#usage) • [Architecture](#architecture) • [Contributing](#contributing)

</div>

---

## Overview

ShellNotes is a minimalist, keyboard-driven note management system built for developers who live in the terminal. Built with Go and the Bubble Tea framework, it provides a seamless note-taking experience without leaving your command line.

## Features

- 📝 **Intuitive Note Management** - Create, edit, and delete notes with simple keyboard shortcuts
- 🔍 **Instant Search** - Real-time filtering across all notes
- 💾 **Smart Persistence** - Automatic file management with manual save control
- 🎨 **Clean UI** - Modern terminal interface powered by Bubble Tea and Lip Gloss
- 🔒 **Privacy First** - Local-only storage, plain Markdown files in `~/.ShellNotes`
- ⚡ **Zero Dependencies** - Single binary with no external runtime requirements

## Demo

![Untitled video - Made with Clipchamp (1) (1)](https://github.com/user-attachments/assets/b526a154-d75c-4170-ba8f-dbf12d52499d)

## Installation

### Using Go Install

```bash
go install github.com/yourusername/ShellNotes@latest
```

### From Source

```bash
git clone https://github.com/yourusername/ShellNotes.git
cd ShellNotes
go build -o shellnotes
sudo mv shellnotes /usr/local/bin/  # Optional: Install globally
```

## Usage

### Launch

```bash
shellnotes
```

### Keyboard Shortcuts

| Key | Action | Context |
|-----|--------|---------|
| `Ctrl+N` | Create new note | Any |
| `Ctrl+L` | Open note list | Any |
| `Ctrl+S` | Save current note | Editor |
| `Ctrl+D` | Delete note | List view |
| `Esc` | Cancel/Go back | Any |
| `Ctrl+Q` | Quit | Any |
| `/` | Filter notes | List view |

### Quick Start

1. Press `Ctrl+N` to create a note
2. Type your content
3. Press `Ctrl+S` to save
4. Press `Ctrl+L` to view all notes

## Architecture

ShellNotes is built using the **Elm Architecture** (Model-View-Update pattern) via the Bubble Tea framework, providing a clean separation of concerns and predictable state management.

### Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                     Bubble Tea Runtime                   │
│                   (Event Loop Manager)                   │
└─────────────────────────────────────────────────────────┘
                            │
                ┌───────────┴───────────┐
                ▼                       ▼
        ┌──────────────┐        ┌──────────────┐
        │   Update()   │        │    View()    │
        │              │        │              │
        │ • Handles    │        │ • Renders    │
        │   messages   │        │   UI state   │
        │ • Updates    │◄───────┤ • Returns    │
        │   model      │        │   string     │
        │ • Returns    │        │              │
        │   commands   │        │              │
        └──────────────┘        └──────────────┘
                │
                ▼
        ┌──────────────┐
        │    Model     │
        │              │
        │ • App state  │
        │ • File refs  │
        │ • UI state   │
        └──────────────┘
                │
                ▼
        ┌──────────────┐
        │ File System  │
        │              │
        │ ~/.ShellNotes│
        │   *.md files │
        └──────────────┘
```

### Key Components

**1. Model (State)**
```go
type model struct {
    newFileInput           textinput.Model  // File name input
    createFileInputVisible bool             // UI state flag
    currentFile            *os.File         // Active file handle
    noteTextArea           textarea.Model   // Text editor
    list                   list.Model       // Note list
    showingList            bool             // View state
    isFileUnsaved          bool             // Tracking flag
}
```

**2. Update (Logic)**
- Receives keyboard messages and window events
- Updates model state based on user actions
- Handles file I/O operations (create, read, write, delete)
- Returns updated model and optional commands

**3. View (Presentation)**
- Renders current state to terminal string
- Displays appropriate UI based on model state
- Uses Lip Gloss for styling and layout

### Design Patterns

- **Elm Architecture**: Unidirectional data flow ensures predictable state
- **Component-based UI**: Modular Bubble Tea components (textinput, textarea, list)
- **Single Responsibility**: Each function handles one concern
- **Immutable Updates**: Model changes return new state rather than mutating

### Tech Stack

- **Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework implementing Elm Architecture
- **Components**: [Bubbles](https://github.com/charmbracelet/bubbles) - Pre-built TUI components (text input, textarea, list)
- **Styling**: [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal layout and styling
- **Language**: Go 1.21+ with standard library for file operations

## Development

```bash
# Clone and setup
git clone https://github.com/yourusername/ShellNotes.git
cd ShellNotes
go mod download

# Run in development
go run main.go

# Build
go build -o shellnotes
```

## Roadmap

- [ ] Confirmation dialogs for destructive actions
- [ ] Note categories and tags
- [ ] Full-text search
- [ ] Export to PDF/HTML
- [ ] Syntax highlighting
- [ ] Custom vault locations

## Contributing

Contributions welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

1. Fork the repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add: AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

Built with [Charm](https://charm.sh/) libraries and inspired by terminal-based workflows.

---

<div align="center">

**[⬆ back to top](#shellnotes)**

Made with ❤️ using Go and Bubble Tea

</div>
