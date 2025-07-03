# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Shush is a CLI tool written in Go that removes comments from source code files using sed under the hood. It supports multiple programming languages and can process individual files or directories (with optional recursive traversal).

## Build and Development Commands

```bash
# Build the binary
go build -o shush ./cmd/shush

# Run the tool
./shush file.py
./shush src/ --recursive --dry-run

# Test with dry-run to see preview
./shush file.js --dry-run --verbose
```

## Architecture

### Project Structure
- `cmd/shush/main.go` - Entry point, CLI parsing with kong, version management
- `internal/types/types.go` - Core type definitions (CLI struct, Language struct, BlockComment)
- `internal/processor/processor.go` - Main processing logic, file/directory handling, sed command execution
- `internal/processor/languages.go` - Language detection and mapping (file extension → comment syntax)

### Core Flow
1. **CLI Parsing**: Kong parses arguments into `types.CLI` struct
2. **Language Detection**: File extension mapped to comment patterns via `languageMap`
3. **Processing Strategy**: 
   - Single file: Direct processing
   - Directory: Scan for supported files (recursive if `-r` flag)
   - Preview mode (`--dry-run`): Show colored diff without changes
4. **Sed Command Generation**: Build sed patterns based on language and flags (`--inline`, `--block`)
5. **Execution**: Run sed commands or show preview with color-coded output

### Key Design Principles
- **Comment Removal Logic**: 
  - Comment-only lines are deleted entirely
  - Inline comments are stripped but lines preserved
  - Original blank lines remain untouched (preserves file structure)
- **Language Support**: Extensible via `languageMap` in `languages.go`
- **Sed Integration**: Leverages existing sed for performance and reliability

### Comment Processing Rules
- Line comments (`//`, `#`, `--`): Remove entire line if comment-only, strip inline comments
- Block comments (`/* */`): Remove single-line or multi-line blocks
- Mutual exclusion: `--inline` and `--block` flags cannot be used together
- Structure preservation: Never remove intentional blank lines

## Release Process

- Version is managed in `cmd/shush/main.go` (`var version`)
- GitHub Actions builds cross-platform binaries (Linux/macOS, x86_64/ARM64) on tag push
- Release workflow uses artifact pattern to avoid race conditions
- Tag format: `v0.1.1` triggers automatic release with all platform binaries

## Language Extension

To add new language support, update `languageMap` in `internal/processor/languages.go`:

```go
"ext": {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
```