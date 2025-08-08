![License](https://img.shields.io/badge/license-MIT-blue)
![Version](https://img.shields.io/badge/version-v0.4-orange)
![Built with Go](https://img.shields.io/badge/built%20with-Go-informational)
![Offline First](https://img.shields.io/badge/offline-first-success)

# tmdr 

## too medical; didn't read

A fast, offline terminal tool for looking up medical acronyms. Built for engineers in healthtech.

**tmdr** gives you instant, offline access to medical acronyms — no context switching, no token burn, 

## Features

- ⚡ **Instant lookup** - inline medical acronyms ('tmdr <acronym>')
- 🎨 **Bubbles TUI** - using charm bubbles ui and lipgloss
- 🔍 **Real-time search** - Type to filter results instantly
- 🎯 **Fuzzy matching** - Handles typos gracefully
- 📚 **Zero dependencies** - Works completely offline
- 🚀 **Cross-platform** - Mac, Linux, Windows ready

## Installation

### Option 1: Use curl

```bash
curl -sSL https://tmdr.sh/install | bash
```

### Option 2: Download Pre-built Binary

Download the latest release for your platform from [GitHub Releases](https://github.com/anthony-langham/tmdr/releases):

### Option 3: Install from Source

Requires Go 1.21 or later:

```bash
git clone https://github.com/anthony-langham/tmdr.git
cd tmdr
make install  # Installs to $GOPATH/bin
```

## Usage

### Interactive Mode (Default)

```bash
tmdr  # Launch beautiful TUI using bubbles and lipgloss
```

### Search Mode
- Type to search in real-time
- Arrow keys to navigate results
- Enter to view full definition
- ESC to clear or exit

### Browse Mode
- Navigate all acronyms with arrow keys
- See full definitions instantly

## Examples

```bash
$ tmdr
# Launches interactive TUI with orange-themed interface

# In search mode, type "abg" to find:
ABG → Arterial Blood Gas
A test measuring oxygen and carbon dioxide levels in arterial blood.

# Browse shows all 107 acronyms like:
ECG → Electrocardiogram
ICU → Intensive Care Unit
MRI → Magnetic Resonance Imaging
```

## Building from Source

```bash
# Build for current platform
make build

# Build for all platforms
make dist

# Create release archives
make release

# Run tests
make test

# Clean build artifacts
make clean
```

## Development Status

Production Ready
- ✅ Interactive Terminal UI with BubbleTea
- ✅ Real-time search with fuzzy matching
- ✅ 107 medical acronyms database
- ✅ Cross-platform support
- ✅ Product feedback integration

## Contributing

Contributions welcome! Feel free to:
- Add more medical acronyms to `data/acronyms.csv`
- Report issues or suggest features
- Submit pull requests

## Feedback

We'd love to hear from you! 
- Press `f` in the app to send feedback. 
- email hello@tmdr.sh

## License

MIT License - See [LICENSE](LICENSE) for details.

***Made up north 🐝***

