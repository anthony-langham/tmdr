# Development README

## Project Structure

```
tmdr/
├── main.go                 # Main entry point
├── go.mod                  # Go module definition
├── cmd/                    # Command-line specific code (future)
│   └── tmdr/              # CLI application
├── internal/              # Private application code
│   ├── acronym/           # Acronym domain logic
│   │   ├── acronym.go     # Core types and interfaces
│   │   └── csv_repository.go # CSV data source
│   └── cli/               # CLI utilities (future)
├── data/                  # Data files
│   └── acronyms.csv       # Medical acronyms database
└── .claude/               # Claude Code specific files
    ├── docs/              # Documentation
    ├── plans/             # Implementation plans
    └── TODO.md            # Task tracking
```

## Architecture

The project follows a clean architecture approach:

1. **Domain Layer** (`internal/acronym/`): Core business logic and types
2. **Data Layer** (`internal/acronym/csv_repository.go`): Data access implementation
3. **Presentation Layer** (`main.go`, future `cmd/`): CLI interface

## Key Design Decisions

- **Offline-first**: All data embedded/local, no network dependencies
- **Interface-based**: Repository pattern for flexible data sources
- **Minimal dependencies**: Using standard library where possible
- **Case-insensitive**: All lookups normalized to uppercase

## Building

```bash
# Build the binary
go build -o tmdr

# Run directly
go run main.go

# Run tests
go test ./...
```

## Adding Acronyms

Edit `data/acronyms.csv` with format:
```csv
acronym,definition
ABG,Arterial Blood Gas – A test measuring oxygen and CO2 levels
```

The definition format is: "Full Form – Description"