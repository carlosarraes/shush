package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/carlosarraes/shush/internal/processor"
	"github.com/carlosarraes/shush/internal/types"
)

var version = "0.1.2"

func main() {
	var cli types.CLI
	kong.Parse(&cli,
		kong.Description("Remove comments from source code files"),
		kong.Vars{"version": version})

	if cli.LLM {
		showLLMGuide()
		return
	}


	gitFlags := []bool{cli.ChangesOnly, cli.Staged, cli.Unstaged}
	gitFlagCount := 0
	for _, flag := range gitFlags {
		if flag {
			gitFlagCount++
		}
	}

	if gitFlagCount > 1 {
		fmt.Fprintf(os.Stderr, "Error: git flags (--changes-only, --staged, --unstaged) are mutually exclusive\n")
		os.Exit(1)
	}

	if gitFlagCount > 0 && cli.Path != "" {
		fmt.Fprintf(os.Stderr, "Error: cannot use git flags with explicit path argument\n")
		os.Exit(1)
	}

	if gitFlagCount > 0 && cli.Recursive {
		fmt.Fprintf(os.Stderr, "Error: cannot use git flags with --recursive (git handles repository scope)\n")
		os.Exit(1)
	}


	if gitFlagCount == 0 && cli.Path == "" {
		fmt.Fprintf(os.Stderr, "Error: path argument is required\n")
		os.Exit(1)
	}

	if cli.Inline && cli.Block {
		fmt.Fprintf(os.Stderr, "Error: --inline and --block flags are mutually exclusive\n")
		os.Exit(1)
	}

	proc := processor.New(cli)
	if err := proc.Process(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func showLLMGuide() {
	fmt.Print(`# Shush CLI - LLM Guide

## Overview
Shush is a fast comment removal tool for source code files using sed under the hood.
- **Purpose**: Remove comments from source code while preserving file structure
- **Key Strength**: Processes individual files or entire directories with recursive support
- **LLM-Friendly**: Supports dry-run mode with colored preview for safe operation

## Core Commands

### Basic File Processing
` + "```bash" + `
shush file.py                    # Remove all comments from single file
shush file.js --inline          # Remove only line comments (// in JS)
shush file.c --block            # Remove only block comments (/* */ in C)
shush script.sh --dry-run       # Preview changes without modification
shush config.lua --backup       # Create backup before processing
` + "```" + `

### Directory Processing  
` + "```bash" + `
shush src/                       # Process all supported files in directory
shush . --recursive              # Process current directory recursively
shush project/ -r --verbose     # Recursive with detailed output
shush src/ --dry-run --verbose  # Preview recursive changes
` + "```" + `

### Advanced Usage
` + "```bash" + `
# Safe exploration workflow
shush project/ -r --dry-run --verbose  # 1. Preview all changes
shush project/ -r --backup             # 2. Process with backups

# Selective comment removal
shush src/ -r --inline --backup        # Remove only line comments
shush src/ -r --block --dry-run        # Preview block comment removal

# Combined operations
shush . --recursive --inline --dry-run --verbose
` + "```" + `

## Supported Languages & Comment Types

### Line Comments Only
- **Python**: ` + "`#`" + ` comments
- **Shell/Bash**: ` + "`#`" + ` comments  
- **Ruby**: ` + "`#`" + ` comments
- **Perl**: ` + "`#`" + ` comments
- **YAML**: ` + "`#`" + ` comments
- **Lua**: ` + "`--`" + ` comments

### Line + Block Comments
- **JavaScript/TypeScript**: ` + "`//`" + ` and ` + "`/* */`" + `
- **Go**: ` + "`//`" + ` and ` + "`/* */`" + `
- **C/C++**: ` + "`//`" + ` and ` + "`/* */`" + `
- **Java**: ` + "`//`" + ` and ` + "`/* */`" + `

## Processing Behavior

### Comment Removal Logic
- **Comment-only lines**: Deleted entirely (preserves line structure)
- **Inline comments**: Stripped but line kept (e.g., ` + "`code(); // comment`" + ` → ` + "`code();`" + `)
- **Block comments**: Removed (single-line or multi-line)
- **Blank lines**: Original empty lines preserved (file structure maintained)

### File Selection
- **Auto-detection**: Language determined by file extension
- **Recursive mode**: Scans subdirectories when ` + "`-r/--recursive`" + ` used
- **Supported only**: Ignores unsupported file types automatically
- **Error handling**: Continues processing other files if one fails

## Flag Combinations

### Mutually Exclusive
` + "```bash" + `
shush file.js --inline --block   # ❌ ERROR: Cannot use both
` + "```" + `

### Recommended Workflows
` + "```bash" + `
# Safe exploration
shush project/ -r --dry-run --verbose

# Production processing  
shush project/ -r --backup --verbose

# Selective processing
shush src/ -r --inline --dry-run    # Preview line comment removal
shush src/ -r --inline --backup     # Apply line comment removal
` + "```" + `

## Output Modes

### Dry-Run Preview (--dry-run)
- **Color-coded display**: Red strikethrough for deleted lines, green for kept
- **Line numbers**: Easy reference for changes
- **Summary stats**: Count of lines to be removed/kept
- **Zero risk**: No files modified

### Verbose Mode (--verbose)  
- **File discovery**: Shows which files found and processed
- **Language detection**: Displays detected language per file
- **Command execution**: Shows sed commands being run
- **Progress tracking**: File-by-file processing status

### Backup Mode (--backup)
- **Safety net**: Creates ` + "`.bak`" + ` files before modification
- **Original preservation**: Backup contains exact original content
- **Per-file basis**: Each processed file gets individual backup

## Git-Aware Processing

### Git Mode Commands
` + "```bash" + `
# Process all changes (staged + unstaged + untracked)
shush --changes-only                 # Remove comments from all changed files
shush --changes-only --dry-run       # Preview changes across entire repository

# Process only staged changes
shush --staged                       # Clean comments from staged files
shush --staged --dry-run --verbose   # Preview staged changes with details

# Process only unstaged changes  
shush --unstaged                     # Clean comments from unstaged work
shush --unstaged --inline           # Remove only line comments from unstaged files
` + "```" + `

### Git Workflow Examples
` + "```bash" + `
# Pre-commit cleanup workflow
shush --staged --dry-run             # 1. Review what will be cleaned
shush --staged --backup              # 2. Clean staged changes with backup
git commit -m "Clean implementation" # 3. Commit cleaned code

# Feature development cleanup
shush --unstaged --dry-run           # 1. Preview unstaged work cleanup  
shush --unstaged --inline            # 2. Remove only debug comments
shush --changes-only                 # 3. Clean all changes before review

# Safe exploration workflow
shush --changes-only --dry-run --verbose  # See all changes that would be made
shush --staged --backup --verbose         # Process with maximum safety
` + "```" + `

### Git Mode Behavior
- **Surgical precision**: Only processes lines that have been changed
- **Repository scope**: Automatically processes relevant files across the repo
- **Change detection**: Uses git diff to identify modified line ranges
- **Untracked files**: Processes entirely (no previous version to compare)
- **Preserves existing code**: Comments in unchanged lines remain untouched

### Git Flag Rules
- **Mutually exclusive**: Cannot combine ` + "`--staged`" + `, ` + "`--unstaged`" + `, ` + "`--changes-only`" + `
- **No explicit paths**: Git flags work on repository scope, not individual files
- **No recursive flag**: Git mode handles repository traversal automatically
- **Compatible with**: ` + "`--inline`" + `, ` + "`--block`" + `, ` + "`--dry-run`" + `, ` + "`--backup`" + `, ` + "`--verbose`" + `

### Git Error Scenarios
- **Not in repository**: Clear error message when git flags used outside git repo
- **No changes found**: Informative message when no staged/unstaged changes exist
- **Git command failures**: Graceful handling of git command errors

## Best Practices for LLM Integration

1. **Always start with dry-run** for unknown codebases
2. **Use recursive + verbose** for comprehensive analysis
3. **Create backups** before processing important code
4. **Test on small directories** before full project processing
5. **Combine flags strategically** (e.g., --recursive --dry-run --verbose)

## Error Scenarios
- **File not found**: Clear error message, continues with other files
- **Permission denied**: Skips file, continues processing
- **Unsupported extension**: Ignores file, shows in verbose mode
- **Directory not found**: Error message, exits
- **No supported files**: Error message for directories

## Command Categories by Priority
1. **Essential**: Basic file processing (` + "`shush file.py`" + `)
2. **Important**: Directory processing (` + "`shush src/ -r`" + `)
3. **Safety**: Dry-run and backup modes (` + "`--dry-run`" + `, ` + "`--backup`" + `)
4. **Selective**: Comment type filtering (` + "`--inline`" + `, ` + "`--block`" + `)

Shush excels at safe, fast comment removal with excellent preview capabilities for confident code processing.
`)
}
