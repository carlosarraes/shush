# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Shush is a CLI tool written in Go that removes comments from source code files using sed under the hood. It supports multiple programming languages and can process individual files or directories (with optional recursive traversal). The tool emphasizes preserving file structure while providing precise comment removal control.

## Build and Development Commands

```bash
# Show all available commands
make help

# Build the binary (optimized)
make build

# Quick development build (faster, unoptimized)
make build-dev

# Run all checks (format, vet, lint, test)
make check

# Development cycle (build + quick test)
make dev

# Get comprehensive LLM-friendly usage guide
./shush --llm

# Example usage
./shush file.py
./shush src/ --recursive --dry-run --verbose
```

## Architecture

### Project Structure
- `cmd/shush/main.go` - Entry point, CLI parsing with kong, version management, LLM guide
- `internal/types/types.go` - Core type definitions (CLI struct, Language struct, BlockComment)
- `internal/processor/processor.go` - Main processing logic, file/directory handling, sed command execution, colored preview
- `internal/processor/languages.go` - Language detection and mapping (file extension → comment syntax)
- `ai_docs/` - Design documents for future features (GIT_FEATURE_PLAN.md)

### Core Flow
1. **CLI Parsing**: Kong parses arguments into `types.CLI` struct
2. **Language Detection**: File extension mapped to comment patterns via `languageMap`
3. **Processing Strategy**: 
   - Single file: Direct processing
   - Directory: Scan for supported files (recursive if `-r` flag)
   - Preview mode (`--dry-run`): Show colored diff with line numbers and counts
   - LLM mode (`--llm`): Display comprehensive usage guide
4. **Sed Command Generation**: Build sed patterns based on language and flags (`--inline`, `--block`)
5. **Execution**: Run sed commands or show preview with color-coded output using fatih/color

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

## Future Development

### Planned Git Integration
A major feature is planned for git-aware comment removal (see `ai_docs/GIT_FEATURE_PLAN.md`):
- `--staged`: Remove comments only from staged changes
- `--unstaged`: Remove comments only from unstaged changes  
- `--changes-only`: Process all git changes (staged + unstaged + untracked)

This would enable surgical comment removal from only the lines you've changed, preserving existing codebase comments.

## Development Notes

### Common Development Workflow
1. `make dev` - Quick build and test cycle
2. `make check` - Run all quality checks before committing
3. `make release` - Prepare for release (runs full check suite)
4. Update version in `cmd/shush/main.go`, then `make tag && make push`

### Adding New Features
- Version management: Update `var version` in `cmd/shush/main.go`
- CLI flags: Add to `types.CLI` struct with appropriate kong tags
- Language support: Extend `languageMap` in `languages.go`
- For LLM integration patterns: Reference existing `--llm` implementation

### Makefile Targets
The Makefile provides comprehensive build, test, and release automation. Use `make help` to see all available commands including cross-platform builds, code quality checks, and release management.