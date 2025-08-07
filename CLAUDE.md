# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**tmdr** (too medical; didn't read) is a fast, offline terminal tool for looking up medical acronyms, built for engineers in healthtech. It's a Go-based CLI application that provides instant, offline access to medical acronyms without context switching or API dependencies.

## Project Status

This project is actively being developed. Current features implemented:
- ✅ Interactive Terminal UI with BubbleTea framework
- ✅ Search-as-you-type functionality with fuzzy matching
- ✅ Browse all acronyms with navigation
- ✅ Product-Market Fit feedback form with Huh forms
- ✅ Orange-themed UI consistent with tmdr branding
- ✅ ESC key navigation with proper back/exit behavior
- ✅ Mailto-based feedback submission (zero infrastructure)

## Technology Stack

- **Language**: Go
- **CLI Framework**: Starting with `flag`, potentially migrating to Cobra
- **TUI Framework**: Charmbracelet BubbleTea (for v0.2+)
- **Styling**: Lip Gloss (for terminal styling)
- **Type**: Terminal/CLI application
- **Architecture**: Offline-first design with embedded data

## Development Setup

Go module is initialized. To build and run:
```bash
# Build the binary
go build -o tmdr

# Run directly
go run main.go

# Run with flags
go run main.go --version
go run main.go --help
```

## Implementation Plan

The project follows a phased approach (see `.claude/plans/plan-v002.md`):

### Phase 1: MVP (v0.1)
- Basic CLI with acronym lookup
- Local CSV data source
- Case-insensitive search
- Essential flags (--help, --version, --random)

### Phase 2: Enhanced Features (v0.2)
- Fuzzy matching for typos
- Terminal UI with BubbleTea
- Learning modes (--daily, --learn)

### Phase 3: Production Ready (v1.0)
- Comprehensive testing
- Cross-platform builds
- Homebrew/Scoop distribution
- User-submitted terms

## Data Format

Acronyms stored in CSV format:
```csv
acronym,definition
ABG,Arterial Blood Gas – A test measuring oxygen and CO2 levels.
HIV,Human Immunodeficiency Virus – A virus that attacks the immune system.
```

## Commands

```bash
# Launch interactive TUI (default)
tmdr

# Navigation in TUI
- Press 's' to search for acronyms
- Press 'b' to browse all acronyms  
- Press 'f' to provide feedback
- Press 'h' or 't' to return home
- Press 'q' or Ctrl+C to quit
- ESC key: context-aware back/exit

# Search mode
- Type to search in real-time
- Arrow keys to navigate results
- Enter to select
- ESC to clear search or exit

# Feedback form
- Navigate questions with Tab/Shift+Tab
- Select answers with Enter/Space
- ESC to go back (exits at first question)
- Form submits via mailto to hello@tmdr.sh
```