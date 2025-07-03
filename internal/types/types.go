package types

import "github.com/alecthomas/kong"

type CLI struct {
	Path      string `arg:"" name:"path" help:"Source code file or directory to process" optional:""`
	Inline    bool   `help:"Remove only line comments"`
	Block     bool   `help:"Remove only block comments"`
	Recursive bool   `short:"r" help:"Process directories recursively"`
	DryRun    bool   `help:"Show what would be removed without making changes"`
	Backup    bool   `help:"Create backup files before modification"`
	Verbose   bool   `help:"Show detailed output"`
	LLM       bool   `help:"Show LLM-friendly usage guide"`
	Version   kong.VersionFlag `help:"Show version information"`
}

type Language struct {
	LineComment  string
	BlockComment *BlockComment
}

type BlockComment struct {
	Start string
	End   string
}