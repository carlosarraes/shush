# shush 🤫

Remove comments from source code files blazingly fast using sed under the hood. Features git-aware processing and Claude Code integration.

## Installation

### Quick Install (Linux/macOS)

```bash
curl -sSf https://raw.githubusercontent.com/carlosarraes/shush/main/install.sh | sh
```

### Manual Download

Download the binary for your platform from the [releases page](https://github.com/carlosarraes/shush/releases).

## Usage

```bash
# Remove all comments from a file
shush file.py

# Remove all comments from a directory
shush src/

# Process directories recursively
shush src/ --recursive

# Remove only line comments (// or #)
shush file.js --inline

# Remove only block comments (/* */)
shush file.c --block

# Preview what would be removed (dry run)
shush script.sh --dry-run

# Create a backup before modifying
shush config.lua --backup

# Verbose output
shush app.go --verbose

# Git-aware processing (only process changed lines)
shush --staged                    # Clean comments from staged changes
shush --unstaged                  # Clean comments from unstaged changes  
shush --changes-only              # Clean comments from all changes

# Combine flags for complex operations
shush src/ --recursive --inline --dry-run --verbose
```

## Supported Languages

| Language | Line Comments | Block Comments |
|----------|--------------|----------------|
| Python   | `#`          | -              |
| JavaScript/TypeScript | `//` | `/* */` |
| Go       | `//`         | `/* */`        |
| C/C++    | `//`         | `/* */`        |
| Java     | `//`         | `/* */`        |
| Lua      | `--`         | -              |
| Shell/Bash | `#`        | -              |
| Ruby     | `#`          | -              |
| Perl     | `#`          | -              |
| YAML     | `#`          | -              |

## Options

```
--inline       Remove only line comments
--block        Remove only block comments
-r, --recursive Process directories recursively
--dry-run      Show what would be removed without making changes
--backup       Create backup files before modification
--verbose      Show detailed output

# Git-aware flags
--changes-only Remove comments only from git changes (staged + unstaged + untracked)
--staged       Remove comments only from staged git changes
--unstaged     Remove comments only from unstaged git changes

# Utility flags
--version      Show version information
--llm          Show LLM-friendly usage guide
--help         Show help message
```

## Examples

### Python
```bash
# Before
# This is a comment
def hello():
    print("Hello")  # Inline comment

# After running: shush example.py
def hello():
    print("Hello")
```

### JavaScript
```bash
# Remove only line comments, preserve block comments
shush app.js --inline

# Remove only block comments, preserve line comments
shush app.js --block
```

### Directory Processing
```bash
# Process all supported files in a directory
shush src/ --verbose

# Process directories recursively
shush . --recursive --dry-run

# Process only specific comment types in entire project
shush . --recursive --inline --backup
```

### Git-Aware Processing
```bash
# Clean comments from staged changes before commit
shush --staged --dry-run          # Preview changes
shush --staged --backup           # Apply with backup

# Clean comments from current work
shush --unstaged --inline         # Remove only line comments
shush --changes-only              # Clean all changes (staged + unstaged + untracked)

# Pre-commit workflow
shush --staged --dry-run          # 1. Review what will be cleaned  
shush --staged                    # 2. Clean staged changes
git commit -m "Clean code"        # 3. Commit cleaned code
```

### Backup and Preview
```bash
# Always create backup before modifying
shush important.go --backup

# See what would be removed first
shush config.yaml --dry-run --verbose
```

## Claude Code Integration 🤖

*Coming soon in v0.2.0* - Seamless integration with Claude Code via hooks:

```bash
# Install automatic comment cleanup after Claude modifies files
shush --install-hooks              # User-wide (all projects)
shush --install-hooks project      # Project-specific only

# Check hook status  
shush --hooks-status               # See current configuration
```

Once installed, comments will be automatically cleaned whenever Claude Code uses Write, Edit, or MultiEdit tools. No manual intervention required!

## How It Works

shush uses optimized sed commands to remove comments while preserving code structure. It:
- Auto-detects language from file extension
- Builds appropriate sed patterns for the detected language
- **Git-aware processing**: Only processes changed lines for surgical precision
- Preserves strings and code that might look like comments
- **Claude Code integration**: Automatic cleanup via PostToolUse hooks

## Building from Source

```bash
git clone https://github.com/carlosarraes/shush.git
cd shush
go build -o shush
```

## Requirements

- sed (available on all Unix-like systems)
- Linux or macOS (x86_64 or ARM64)

## License

MIT