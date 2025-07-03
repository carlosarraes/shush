# shush 🤫

Remove comments from source code files blazingly fast using sed under the hood.

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

### Backup and Preview
```bash
# Always create backup before modifying
shush important.go --backup

# See what would be removed first
shush config.yaml --dry-run --verbose
```

## How It Works

shush uses optimized sed commands to remove comments while preserving code structure. It:
- Auto-detects language from file extension
- Builds appropriate sed patterns for the detected language
- Removes comments and empty lines in a single pass
- Preserves strings and code that might look like comments

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