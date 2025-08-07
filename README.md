![License](https://img.shields.io/badge/license-MIT-blue)
![Version](https://img.shields.io/badge/version-v0.3-orange)
![Built with Go](https://img.shields.io/badge/built%20with-Go-informational)
![Offline First](https://img.shields.io/badge/offline-first-success)

# tmdr — too medical; didn't read

A fast, offline terminal tool for looking up medical acronyms. Built for engineers in healthtech.

🩺 **tmdr** gives you instant, offline access to medical acronyms — no context switching, no token burn, no bs.

## ✨ Interactive Terminal UI

```bash
tmdr  # Launch interactive mode (default)
```

Navigate with ease:
- **`s`** - Search acronyms in real-time
- **`b`** - Browse all acronyms
- **`f`** - Send feedback
- **`q`** - Quit

## Features

- ⚡ **Instant lookup** - medical acronyms available offline
- 🎨 **Beautiful TUI** - Orange-themed interactive interface
- 🔍 **Real-time search** - Type to filter results instantly
- 🎯 **Fuzzy matching** - Handles typos gracefully
- 📚 **Zero dependencies** - Works completely offline
- 🚀 **Cross-platform** - Mac, Linux, Windows ready

## Installation

### Option 1: Download Pre-built Binary (Recommended)

Download the latest release for your platform from [GitHub Releases](https://github.com/anthony-langham/tmdr/releases):

```bash
# macOS (Apple Silicon)
curl -L https://github.com/anthony-langham/tmdr/releases/download/v0.3/tmdr-v0.3-darwin-arm64.tar.gz | tar xz
sudo mv tmdr /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/anthony-langham/tmdr/releases/download/v0.3/tmdr-v0.3-darwin-amd64.tar.gz | tar xz
sudo mv tmdr /usr/local/bin/

# Linux
curl -L https://github.com/anthony-langham/tmdr/releases/download/v0.3/tmdr-v0.3-linux-amd64.tar.gz | tar xz
sudo mv tmdr /usr/local/bin/
```

### Option 2: Install from Source

Requires Go 1.21 or later:

```bash
git clone https://github.com/anthony-langham/tmdr.git
cd tmdr
make install  # Installs to $GOPATH/bin
```

### Option 3: Homebrew (Coming Soon)

```bash
brew tap anthony-langham/tmdr
brew install tmdr
```

## Usage

### Interactive Mode (Default)

```bash
tmdr  # Launch beautiful TUI
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

**v0.3** - Production Ready
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

We'd love to hear from you! Press `f` in the app to send feedback directly.

## License

MIT License - See [LICENSE](LICENSE) for details.

---

Made with 🧡 for engineers in healthtech
