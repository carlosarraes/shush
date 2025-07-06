# shush 🤫

Remove comments from source code files blazingly fast using in-memory processing. Features git-aware processing, smart comment preservation, and Claude Code integration.

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

# Keep comment-only lines as empty lines (indentation always preserved)
shush script.py --preserve-lines

# Verbose output
shush app.go --verbose

# Git-aware processing (only process changed lines)
shush --staged                    # Clean comments from staged changes
shush --unstaged                  # Clean comments from unstaged changes  
shush --changes-only              # Clean comments from all changes

# Configuration and hooks management
shush --create-config             # Create .shush.toml configuration
shush --config                    # Show current configuration
shush --install-hook              # Install Claude Code hooks (user-wide)
shush --install-hook -s project   # Install project-specific hooks
shush --hook-status               # Check hook installation status

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
# Comment filtering
--inline       Remove only line comments
--block        Remove only block comments

# Processing modes
-r, --recursive    Process directories recursively
--dry-run          Show what would be removed without making changes
--backup           Create backup files before modification
--verbose          Show detailed output
--preserve-lines   Keep comment-only lines as empty lines (indentation always preserved)

# Git-aware flags
--changes-only Remove comments only from git changes (staged + unstaged + untracked)
--staged       Remove comments only from staged git changes
--unstaged     Remove comments only from unstaged git changes

# Configuration management
--config       Show current configuration and location
--create-config Create example .shush.toml configuration file

# Claude Code hooks
--install-hook   Install Claude Code hooks for automatic comment cleanup
--uninstall-hook Uninstall Claude Code hooks  
--list-hooks     List current Claude Code hooks configuration
--hook-status    Check if shush hooks are installed
-s, --hook-scope Hook scope: 'project' for local, default for user-wide

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
    # Comment-only line
    print("Hello")  # Inline comment

# After running: shush example.py (default - comment-only lines deleted)
def hello():
    print("Hello")

# After running: shush example.py --preserve-lines (comment-only lines kept as empty)
def hello():
    
    print("Hello")

# Note: Code indentation is ALWAYS preserved in both modes
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

### Line Structure Control
```bash
# Default: Comment-only lines are deleted entirely
shush script.py                      # Clean removal

# Keep comment-only lines as empty lines (preserves line numbers)
shush script.py --preserve-lines     # Useful for debugging/line references

# Code indentation is ALWAYS preserved regardless of flag
shush python_code.py --preserve-lines --dry-run
```

### Backup and Preview
```bash
# Always create backup before modifying
shush important.go --backup

# See what would be removed first
shush config.yaml --dry-run --verbose
```

## Comment Preservation Configuration 🎯

Shush supports smart comment preservation through `.shush.toml` configuration files:

```bash
# Create example configuration file
shush --create-config

# Check current configuration  
shush --config
```

### Configuration File (`.shush.toml`)

```toml
# Patterns to preserve in comments (supports wildcards with *)
preserve = [
    "TODO:",
    "FIXME:",
    "HACK:",
    "XXX:",
    "@ts-ignore",
    "@ts-expect-error",
    "eslint-",
    "prettier-ignore",
    "pylint:",
    "mypy:",
    "type: ignore",
    "*IMPORTANT*",   # Wildcard: preserves any comment containing IMPORTANT
    "*DEBUG*",       # Wildcard: preserves any comment containing DEBUG
]
```

### Configuration Discovery

Shush searches for configuration in this order:
1. `.shush.toml` (current directory)
2. `.shush.toml` (git repository root)  
3. `~/.config/.shush.toml` (global user config)

## Claude Code Integration 🤖

Seamless integration with Claude Code via hooks - **available now**:

```bash
# Install automatic comment cleanup after Claude modifies files
shush --install-hook               # User-wide (all projects)
shush --install-hook -s project    # Project-specific only

# Manage hooks
shush --hook-status                # Check installation status
shush --list-hooks                 # Show all configured hooks
shush --uninstall-hook             # Remove hooks

# Hook scope conflict detection
# Prevents duplicate execution when both user and project hooks exist
```

Once installed, comments will be automatically cleaned whenever Claude Code uses Write, Edit, or MultiEdit tools. Respects your `.shush.toml` configuration for comment preservation!

## How It Works

shush uses optimized processing to remove comments while preserving code structure:

- **Language Detection**: Auto-detects language from file extension
- **String-Aware Parsing**: Preserves URLs and strings that contain comment markers (e.g., `"https://example.com"`)
- **Git-Aware Processing**: Only processes changed lines for surgical precision
- **Smart Preservation**: Configurable comment preservation via `.shush.toml` patterns
- **Unified In-Memory Processing**: 
  - String-aware parsing for all file operations
  - Line-based processing with comment preservation for all modes
- **Claude Code Integration**: Automatic cleanup via PostToolUse hooks

## Building from Source

```bash
git clone https://github.com/carlosarraes/shush.git
cd shush
go build -o shush
```

## Requirements

- Linux or macOS (x86_64 or ARM64)

## License

MIT