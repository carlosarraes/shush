package types

import "github.com/alecthomas/kong"

type CLI struct {
	Path          string           `arg:"" name:"path" help:"Source code file or directory to process" optional:""`
	Inline        bool             `help:"Remove only line comments"`
	Block         bool             `help:"Remove only block comments"`
	Recursive     bool             `short:"r" help:"Process directories recursively"`
	DryRun        bool             `help:"Show what would be removed without making changes"`
	Backup        bool             `help:"Create backup files before modification"`
	Verbose       bool             `help:"Show detailed output"`
	PreserveLines bool             `help:"Remove comments but preserve empty lines"`
	ContextLines  int              `short:"c" help:"Number of context lines to show in preview mode (default: from config)" default:"-1"`
	LLM           bool             `help:"Show LLM-friendly usage guide"`
	ChangesOnly   bool             `help:"Remove comments only from git changes (staged + unstaged + untracked)"`
	Staged        bool             `help:"Remove comments only from staged git changes"`
	Unstaged      bool             `help:"Remove comments only from unstaged git changes"`
	InstallHook   bool             `help:"Install Claude Code hooks for automatic comment cleanup"`
	UninstallHook bool             `help:"Uninstall Claude Code hooks"`
	ListHooks     bool             `help:"List current Claude Code hooks configuration"`
	HookStatus    bool             `help:"Check if shush hooks are installed"`
	HookScope     string           `short:"s" help:"Hook scope: 'project' for local, default for user-wide"`
	Config        bool             `help:"Show current configuration and location"`
	CreateConfig  bool             `help:"Create example .shush.toml configuration file"`
	Version       kong.VersionFlag `help:"Show version information"`
}

type Language struct {
	LineComment          string
	AlternateLineComment string
	BlockComment         *BlockComment
}

type BlockComment struct {
	Start string
	End   string
}
