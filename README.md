![License](https://img.shields.io/badge/license-MIT-blue)
![Built with Go](https://img.shields.io/badge/built%20with-Go-informational)
![Offline First](https://img.shields.io/badge/offline-first-success)
![Terminal Native](https://img.shields.io/badge/UX-terminal--native-yellow)

# tmdr — too medical; didn’t read
- A fast, offline terminal tool for looking up medical acronyms. 
- Built for engineers in healthtech.

---

🩺 **tmdr** (too medical; didn’t read) 
- gives you instant, offline access to medical acronyms 
— no context switching
- no token burn
- no bs

```bash
$ tmdr abg
ABG → Arterial Blood Gas
A test measuring oxygen and carbon dioxide levels in arterial blood.
```

## Features

- ⚡ **Instant lookup** - 100+ medical acronyms available offline
- 🎯 **Case-insensitive** - Type `abg`, `ABG`, or `Abg` - all work
- 📚 **No dependencies** - Works completely offline, no API calls
- 🎲 **Random mode** - Learn a new acronym with `--random`
- 🧠 **Curated definitions** - Clear, concise medical explanations

## Installation

### From Source
Requires Go 1.24 or later:

```bash
git clone https://github.com/anthony-langham/tmdr.git
cd tmdr
go build -o tmdr
./tmdr --help
```

## Usage

```bash
# Look up a specific acronym
tmdr abg
tmdr copd

# Display a random acronym (coming soon)
tmdr --random

# Show version
tmdr --version

# Show help
tmdr --help
```

## Examples

```bash
$ tmdr ecg
ECG → Electrocardiogram
Test recording electrical activity of the heart

$ tmdr icu  
ICU → Intensive Care Unit
Hospital unit for critically ill patients
```

## Development Status

Currently in active development (v0.1.0). Core functionality implemented:
- Basic CLI with acronym lookup
- 106 medical acronyms in database
- Case-insensitive search
- Help and version flags

### Coming Soon
- Fuzzy matching for typos
- Interactive TUI mode
- Daily learning mode
- More acronyms

## Contributing

Contributions welcome! Feel free to:
- Add more medical acronyms to `data/acronyms.csv`
- Report issues or suggest features
- Submit pull requests

## License

This project is licensed under the MIT License.
