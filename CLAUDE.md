# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Shush is a CLI tool written in Go that removes comments from source code files while preserving important comments through configuration. It supports multiple programming languages and can process individual files or directories with optional recursive traversal. The tool emphasizes preserving file structure while providing precise comment removal control.

**Key Features**:
- **Git-aware processing**: Surgical precision targeting only changed lines  
- **Claude Code integration**: Automatic comment cleanup via hooks (✅ available now)
- **Smart comment preservation**: Configurable patterns via .shush.toml with wildcard support
- **String-aware parsing**: Preserves URLs and strings containing comment markers
- **Dual processing modes**: Traditional sed-based + in-memory line-based for git operations

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

# Git-aware examples
./shush --staged --dry-run
./shush --changes-only

# Configuration and hooks examples
./shush --create-config
./shush --config
./shush --install-hook
./shush --hook-status
```

## Architecture

### Project Structure
- `cmd/shush/main.go` - Entry point, CLI parsing with kong, version management, LLM guide, hook commands
- `internal/types/types.go` - Core type definitions (CLI struct with hook flags, Language struct, BlockComment)
- `internal/processor/processor.go` - Main processing logic, routing between sed and git modes
- `internal/processor/git_processor.go` - Git-aware processing with line-based comment removal + config integration
- `internal/processor/languages.go` - Language detection and mapping (file extension → comment syntax)
- `internal/git/` - Git operations: repo detection, diff parsing, line range extraction
- `internal/config/` - TOML configuration system with wildcard pattern matching
- `internal/hooks/` - Claude Code hooks integration with conflict detection
- `ai_docs/` - Design documents (GIT_FEATURE_PLAN.md, HOOKS_FEATURE_PLAN.md)

### Core Flow
1. **CLI Parsing**: Kong parses arguments into `types.CLI` struct with hook and git flag validation
2. **Mode Detection**: Route to hook commands, git-aware processing, or traditional processing
3. **Hook Commands** (`--install-hook`, `--hook-status`, etc.):
   - Settings file management for user-wide and project-specific scopes
   - Cross-scope conflict detection and prevention
   - JSON configuration merging with existing Claude Code hooks
4. **Git Mode** (`--staged`, `--unstaged`, `--changes-only`):
   - Configuration loading from .shush.toml files
   - Git repository detection and file change analysis
   - Diff parsing to extract precise line ranges
   - String-aware comment removal with pattern-based preservation
5. **Traditional Mode** (file/directory paths):
   - Language detection from file extension
   - Sed command generation and execution
   - Directory scanning (recursive if `-r` flag)
6. **Output**: Color-coded preview with preserved comment indicators, totals summary, or file modification

### Key Design Principles
- **Dual Processing Architecture**: 
  - Traditional: sed-based for entire files/directories
  - Git-aware: in-memory line-based for surgical precision with config integration
- **String-Aware Processing**: 
  - Preserves URLs and strings containing comment markers (e.g., `"https://example.com"`)
  - Context-aware parsing respects quote boundaries and escaping
- **Comment Preservation Logic**: 
  - Configurable patterns via .shush.toml (exact matches + wildcards)
  - Comment-only lines are deleted entirely unless preserved
  - Inline comments are stripped but lines preserved unless preserved
  - Spacing preserved unless comments are actually removed
- **Git Integration**: Surgical targeting of only changed lines
- **Language Support**: Extensible via `languageMap` in `languages.go`

### Comment Processing Rules
- **Line comments** (`//`, `#`, `--`): Remove entire line if comment-only, strip inline comments
- **Block comments** (`/* */`): Remove single-line or multi-line blocks  
- **Preservation patterns**: Check against .shush.toml preserve patterns before removal
- **String protection**: Comment markers inside strings are ignored
- **Git mode**: Only processes lines within detected change ranges
- **Flag exclusions**: `--inline` and `--block` are mutually exclusive; git flags are mutually exclusive; hook commands are mutually exclusive
- **Structure preservation**: Never remove intentional blank lines; preserve spacing unless comments removed

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

### Completed: Git-Aware Processing (v0.1.3)
The git-aware comment removal feature is now implemented:
- `--staged`: Remove comments only from staged changes ✅
- `--unstaged`: Remove comments only from unstaged changes ✅  
- `--changes-only`: Process all git changes (staged + unstaged + untracked) ✅

### Completed: Configuration System & Claude Code Hooks (v0.2.0)
Smart comment preservation and seamless Claude Code integration:
- TOML configuration with wildcard pattern support ✅
- `--install-hook`: Auto-configure Claude Code to run shush after file modifications ✅
- `--install-hook -s project`: Project-specific hook installation ✅
- Cross-scope conflict detection prevents duplicate execution ✅
- String-aware parsing preserves URLs and code in strings ✅
- Automatic comment cleanup with configurable preservation ✅

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