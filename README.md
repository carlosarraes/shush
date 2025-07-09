# shush 🤫
**S**loppily **H**ushing **U**nwanted **S**ource-code **H**eavy (comments)

Remove comments from source code files blazingly fast. Features Claude Code integration, smart comment preservation, and git-aware processing. Supports [40+ file extensions](SUPPORTED_LANGUAGES.md) across most popular programming languages.

## Installation

### Quick Install (Linux/macOS)
```bash
curl -sSf https://raw.githubusercontent.com/carlosarraes/shush/main/install.sh | sh
```

### Manual Download
Download the binary for your platform from the [releases page](https://github.com/carlosarraes/shush/releases).

## 🤖 Claude Code Integration

Automatic comment cleanup whenever Claude Code modifies files:

```bash
# Install hooks
shush --install-hook              # User-wide (all projects)
shush --install-hook --project    # Project-specific only

# Manage hooks
shush --hook-status               # Check installation status
shush --list-hooks                # Show all configured hooks
shush --uninstall-hook            # Remove hooks
```

Once installed, comments are automatically cleaned whenever Claude Code uses Write, Edit, or MultiEdit tools. Respects your `.shush.toml` configuration for comment preservation.

## ⚙️ Configuration (.shush.toml)

Smart comment preservation through configuration:

```bash
shush --create-config    # Create example configuration
shush --config          # Show current configuration
```

### Configuration File Example
```toml
# Patterns to preserve in comments (supports wildcards with *)
preserve = [
    "TODO:",
    "FIXME:",
    "@ts-ignore",
    "eslint-",
    "*IMPORTANT*",   # Wildcard: preserves any comment containing IMPORTANT
    "*DEBUG*",       # Wildcard: preserves any comment containing DEBUG
]

# Number of context lines to show in preview mode (default: 3)
context_lines = 3
```

### Configuration Discovery
Shush searches for configuration in this order:
1. `.shush.toml` (current directory)
2. `.shush.toml` (git repository root)  
3. `~/.config/.shush.toml` (global user config)

## 🚫 File Exclusion (.shushignore)

Exclude files and directories from processing:

```bash
# Create .shushignore file
echo "*.tmp" > .shushignore       # Ignore all .tmp files
echo "build/" >> .shushignore     # Ignore build directory
echo "test*.js" >> .shushignore   # Ignore test files
echo "!important.js" >> .shushignore  # But keep important.js
```

**Ignore File Locations:**
- `.shushignore` (project root or current directory)
- `~/.config/.shushignore` (global user ignore patterns)

## Usage

### Basic Operations
```bash
# Remove comments from file/directory
shush file.py
shush src/ --recursive

# Preview changes
shush script.sh --dry-run
shush script.sh --dry-run --context-lines 5

# Comment type filtering
shush file.js --inline          # Only line comments
shush file.c --block            # Only block comments

# Backup and preserve options
shush config.lua --backup
shush script.py --preserve-lines  # Keep comment-only lines as empty
```

### Git-Aware Processing
```bash
# Process only changed lines
shush --staged                   # Clean staged changes
shush --unstaged                 # Clean unstaged changes  
shush --changes-only             # Clean all changes (staged + unstaged + untracked)
```

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
--preserve-lines   Keep comment-only lines as empty lines
-c, --context-lines Number of context lines to show in preview mode

# Git-aware flags
--changes-only Remove comments only from git changes
--staged       Remove comments only from staged git changes
--unstaged     Remove comments only from unstaged git changes

# Configuration
--config       Show current configuration and location
--create-config Create example .shush.toml configuration file

# Claude Code hooks
--install-hook   Install Claude Code hooks for automatic comment cleanup
--uninstall-hook Uninstall Claude Code hooks  
--list-hooks     List current Claude Code hooks configuration
--hook-status    Check if shush hooks are installed
--project        Use project scope for hook operations (default: user-wide)

# Utility
--version      Show version information
--llm          Show LLM-friendly usage guide
--help         Show help message
```

## Examples

### Python
```python
# Before
# This is a comment
def hello():
    # Comment-only line
    print("Hello")  # Inline comment

# After: shush example.py
def hello():
    print("Hello")

# After: shush example.py --preserve-lines
def hello():
    
    print("Hello")
```

### Git Workflow
```bash
# Preview and clean staged changes
shush --staged --dry-run          # 1. Review what will be cleaned  
shush --staged                    # 2. Clean staged changes
git commit -m "Clean code"        # 3. Commit cleaned code
```

## How It Works

- **Language Detection**: Auto-detects language from file extension
- **String-Aware Parsing**: Preserves URLs and strings containing comment markers
- **Git-Aware Processing**: Only processes changed lines for surgical precision
- **Smart Preservation**: Configurable comment preservation via `.shush.toml` patterns
- **Claude Code Integration**: Automatic cleanup via PostToolUse hooks

## Requirements

- Linux or macOS (x86_64 or ARM64)

## License

MIT