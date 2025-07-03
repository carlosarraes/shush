package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
)

type CLI struct {
	File    string `arg:"" name:"file" help:"Source code file to process"`
	Inline  bool   `help:"Remove only line comments"`
	Block   bool   `help:"Remove only block comments"`
	DryRun  bool   `help:"Show what would be removed without making changes"`
	Backup  bool   `help:"Create backup file before modification"`
	Verbose bool   `help:"Show detailed output"`
	Version bool   `help:"Show version information"`
}

type Language struct {
	LineComment  string
	BlockComment *BlockComment
}

type BlockComment struct {
	Start string
	End   string
}

var languageMap = map[string]Language{
	"lua":  {LineComment: "--"},
	"py":   {LineComment: "#"},
	"sh":   {LineComment: "#"},
	"js":   {LineComment: "//", BlockComment: &BlockComment{Start: "/*", End: "*/"}},
	"ts":   {LineComment: "//", BlockComment: &BlockComment{Start: "/*", End: "*/"}},
	"go":   {LineComment: "//", BlockComment: &BlockComment{Start: "/*", End: "*/"}},
	"c":    {LineComment: "//", BlockComment: &BlockComment{Start: "/*", End: "*/"}},
	"cpp":  {LineComment: "//", BlockComment: &BlockComment{Start: "/*", End: "*/"}},
	"java": {LineComment: "//", BlockComment: &BlockComment{Start: "/*", End: "*/"}},
	"rb":   {LineComment: "#"},
	"pl":   {LineComment: "#"},
	"yml":  {LineComment: "#"},
	"yaml": {LineComment: "#"},
}

func main() {
	var cli CLI
	kong.Parse(&cli, kong.Description("Remove comments from source code files"))

	if cli.Version {
		fmt.Println("shush version 0.0.1")
		return
	}

	if cli.Inline && cli.Block {
		fmt.Fprintf(os.Stderr, "Error: --inline and --block flags are mutually exclusive\n")
		os.Exit(1)
	}

	if err := processFile(cli); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func processFile(cli CLI) error {
	if _, err := os.Stat(cli.File); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", cli.File)
	}

	language, err := detectLanguage(cli.File)
	if err != nil {
		return err
	}

	if cli.Verbose {
		fmt.Printf("Processing %s...\n", cli.File)
		fmt.Printf("Detected language: %s\n", getLanguageName(cli.File))
		if language.BlockComment != nil {
			fmt.Printf("Comment types: line (%s), block (%s %s)\n", 
				language.LineComment, language.BlockComment.Start, language.BlockComment.End)
		} else {
			fmt.Printf("Comment types: line (%s)\n", language.LineComment)
		}
	}

	sedCmd := buildSedCommand(language, cli)
	
	if cli.DryRun {
		fmt.Printf("Would execute: %s\n", sedCmd)
		return nil
	}

	if cli.Backup {
		if err := createBackup(cli.File); err != nil {
			return fmt.Errorf("failed to create backup: %v", err)
		}
		if cli.Verbose {
			fmt.Printf("✓ Backup created: %s.bak\n", cli.File)
		}
	}

	if cli.Verbose {
		fmt.Printf("Executing: %s\n", sedCmd)
	}

	cmd := exec.Command("sed", "-i", sedCmd, cli.File)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sed command failed: %v", err)
	}

	if cli.Verbose {
		fmt.Printf("✓ Comments removed from %s\n", cli.File)
	}

	return nil
}

func detectLanguage(filename string) (Language, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return Language{}, fmt.Errorf("no file extension found")
	}

	ext = strings.TrimPrefix(ext, ".")
	
	if language, ok := languageMap[ext]; ok {
		return language, nil
	}

	return Language{}, fmt.Errorf("unsupported file extension: %s", ext)
}

func getLanguageName(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	ext = strings.TrimPrefix(ext, ".")
	
	names := map[string]string{
		"lua":  "Lua",
		"py":   "Python",
		"sh":   "Shell",
		"js":   "JavaScript",
		"ts":   "TypeScript",
		"go":   "Go",
		"c":    "C",
		"cpp":  "C++",
		"java": "Java",
		"rb":   "Ruby",
		"pl":   "Perl",
		"yml":  "YAML",
		"yaml": "YAML",
	}
	
	if name, ok := names[ext]; ok {
		return name
	}
	return ext
}

func buildSedCommand(language Language, cli CLI) string {
	var commands []string
	
	if !cli.Block && language.LineComment != "" {
		escaped := escapeForSed(language.LineComment)
		commands = append(commands, fmt.Sprintf("/^[[:space:]]*%s/d", escaped))
		commands = append(commands, fmt.Sprintf("s/%s.*//g", escaped))
	}
	
	if !cli.Inline && language.BlockComment != nil {
		startEscaped := escapeForSed(language.BlockComment.Start)
		endEscaped := escapeForSed(language.BlockComment.End)
		commands = append(commands, fmt.Sprintf("s/%s.*%s//g", startEscaped, endEscaped))
		commands = append(commands, fmt.Sprintf("/%s/,/%s/d", startEscaped, endEscaped))
	}
	
	commands = append(commands, "/^[[:space:]]*$/d")
	
	return strings.Join(commands, "; ")
}

func escapeForSed(pattern string) string {
	replacer := strings.NewReplacer(
		"/", "\\/",
		"*", "\\*",
		".", "\\.",
		"[", "\\[",
		"]", "\\]",
		"^", "\\^",
		"$", "\\$",
		"\\", "\\\\",
	)
	return replacer.Replace(pattern)
}

func createBackup(filename string) error {
	backupName := filename + ".bak"
	
	srcFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	
	dstFile, err := os.Create(backupName)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	
	buffer := make([]byte, 1024)
	for {
		n, err := srcFile.Read(buffer)
		if err != nil && err.Error() != "EOF" {
			return err
		}
		if n == 0 {
			break
		}
		
		if _, err := dstFile.Write(buffer[:n]); err != nil {
			return err
		}
	}
	
	return nil
}