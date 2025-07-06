# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Shush is a CLI tool written in Go that removes comments from source code files using sed under the hood. It supports multiple programming languages and can process individual files or directories (with optional recursive traversal). The tool emphasizes preserving file structure while providing precise comment removal control.

**Key Features**:
- **Git-aware processing**: Surgical precision targeting only changed lines  
- **Claude Code integration**: Automatic comment cleanup via hooks (planned v0.2.0)
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
```

## Architecture

### Project Structure
- `cmd/shush/main.go` - Entry point, CLI parsing with kong, version management, LLM guide, git flag validation
- `internal/types/types.go` - Core type definitions (CLI struct with git flags, Language struct, BlockComment)
- `internal/processor/processor.go` - Main processing logic, routing between sed and git modes
- `internal/processor/git_processor.go` - Git-aware processing with line-based comment removal
- `internal/processor/languages.go` - Language detection and mapping (file extension → comment syntax)
- `internal/git/` - Git operations: repo detection, diff parsing, line range extraction
- `ai_docs/` - Design documents (GIT_FEATURE_PLAN.md, HOOKS_FEATURE_PLAN.md)

### Core Flow
1. **CLI Parsing**: Kong parses arguments into `types.CLI` struct with git flag validation
2. **Mode Detection**: Route to git-aware or traditional processing
3. **Git Mode** (`--staged`, `--unstaged`, `--changes-only`):
   - Git repository detection and file change analysis
   - Diff parsing to extract precise line ranges
   - Line-based comment removal with totals tracking
4. **Traditional Mode** (file/directory paths):
   - Language detection from file extension
   - Sed command generation and execution
   - Directory scanning (recursive if `-r` flag)
5. **Output**: Color-coded preview, totals summary, or file modification

### Key Design Principles
- **Dual Processing Architecture**: 
  - Traditional: sed-based for entire files/directories
  - Git-aware: in-memory line-based for surgical precision
- **Comment Removal Logic**: 
  - Comment-only lines are deleted entirely
  - Inline comments are stripped but lines preserved
  - Spacing preserved unless comments are actually removed
- **Git Integration**: Surgical targeting of only changed lines
- **Language Support**: Extensible via `languageMap` in `languages.go`

### Comment Processing Rules
- **Line comments** (`//`, `#`, `--`): Remove entire line if comment-only, strip inline comments
- **Block comments** (`/* */`): Remove single-line or multi-line blocks
- **Git mode**: Only processes lines within detected change ranges
- **Flag exclusions**: `--inline` and `--block` are mutually exclusive; git flags are mutually exclusive
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

### Planned: Claude Code Hooks Integration (v0.2.0)
Seamless integration with Claude Code via hooks (see `ai_docs/HOOKS_FEATURE_PLAN.md`):
- `--install-hooks`: Auto-configure Claude Code to run shush after file modifications
- `--install-hooks project`: Project-specific hook installation
- Automatic comment cleanup with zero manual intervention

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