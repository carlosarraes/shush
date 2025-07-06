package processor

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/carlosarraes/shush/internal/types"
	"github.com/fatih/color"
)

type Processor struct {
	cli types.CLI
}

func New(cli types.CLI) *Processor {
	return &Processor{cli: cli}
}

func (p *Processor) Process() error {

	if p.cli.ChangesOnly || p.cli.Staged || p.cli.Unstaged {
		return p.processGitChanges()
	}


	info, err := os.Stat(p.cli.Path)
	if os.IsNotExist(err) {
		return fmt.Errorf("path not found: %s", p.cli.Path)
	}
	if err != nil {
		return err
	}

	if info.IsDir() {
		return p.processDirectory(p.cli.Path)
	}

	return p.processFile(p.cli.Path)
}

func (p *Processor) processDirectory(dirPath string) error {
	var files []string

	if p.cli.Recursive {
		err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && IsSupportedFile(path) {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return err
		}
	} else {
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				fullPath := filepath.Join(dirPath, entry.Name())
				if IsSupportedFile(fullPath) {
					files = append(files, fullPath)
				}
			}
		}
	}

	if len(files) == 0 {
		return fmt.Errorf("no supported files found in directory: %s", dirPath)
	}

	if p.cli.Verbose {
		fmt.Printf("Found %d supported files to process\n", len(files))
	}

	for _, file := range files {
		if p.cli.Verbose {
			fmt.Printf("Processing: %s\n", file)
		}

		if err := p.processFile(file); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", file, err)
			continue
		}
	}

	return nil
}

func (p *Processor) processFile(filename string) error {
	language, err := DetectLanguage(filename)
	if err != nil {
		return err
	}

	if p.cli.Verbose {
		fmt.Printf("Processing %s...\n", filename)
		fmt.Printf("Detected language: %s\n", GetLanguageName(filename))
		if language.BlockComment != nil {
			fmt.Printf("Comment types: line (%s), block (%s %s)\n",
				language.LineComment, language.BlockComment.Start, language.BlockComment.End)
		} else {
			fmt.Printf("Comment types: line (%s)\n", language.LineComment)
		}
	}

	if p.cli.DryRun {
		return p.showPreview(filename, language)
	}

	sedCmd := p.buildSedCommand(language)

	if p.cli.Backup {
		if err := p.createBackup(filename); err != nil {
			return fmt.Errorf("failed to create backup: %v", err)
		}
		if p.cli.Verbose {
			fmt.Printf("✓ Backup created: %s.bak\n", filename)
		}
	}

	if p.cli.Verbose {
		fmt.Printf("Executing: %s\n", sedCmd)
	}

	cmd := exec.Command("sed", "-i", sedCmd, filename)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sed command failed: %v", err)
	}

	if p.cli.Verbose {
		fmt.Printf("✓ Comments removed from %s\n", filename)
	}

	return nil
}

func (p *Processor) buildSedCommand(language types.Language) string {
	var commands []string

	if !p.cli.Block && language.LineComment != "" {
		escaped := escapeForSed(language.LineComment)
		// Delete lines that are only comments (optionally with whitespace)
		commands = append(commands, fmt.Sprintf("/^[[:space:]]*%s/d", escaped))
		// Remove inline comments but keep the line
		commands = append(commands, fmt.Sprintf("s/%s.*//g", escaped))
	}

	if !p.cli.Inline && language.BlockComment != nil {
		startEscaped := escapeForSed(language.BlockComment.Start)
		endEscaped := escapeForSed(language.BlockComment.End)
		// Remove block comments on same line
		commands = append(commands, fmt.Sprintf("s/%s.*%s//g", startEscaped, endEscaped))
		// Remove multi-line block comments
		commands = append(commands, fmt.Sprintf("/%s/,/%s/d", startEscaped, endEscaped))
	}

	// Don't remove empty lines - preserve original file structure

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

func (p *Processor) createBackup(filename string) error {
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

func (p *Processor) showPreview(filename string, language types.Language) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	red := color.New(color.FgRed, color.CrossedOut)
	green := color.New(color.FgGreen)
	gray := color.New(color.FgHiBlack)
	yellow := color.New(color.FgYellow)

	scanner := bufio.NewScanner(file)
	lineNum := 0
	deletedCount := 0
	keptCount := 0

	// Compile regex patterns once
	var lineRegex, blockStartRegex, blockEndRegex *regexp.Regexp
	if language.LineComment != "" && !p.cli.Block {
		escaped := regexp.QuoteMeta(language.LineComment)
		lineRegex = regexp.MustCompile(fmt.Sprintf(`^\s*%s|%s.*$`, escaped, escaped))
	}
	if language.BlockComment != nil && !p.cli.Inline {
		blockStartRegex = regexp.MustCompile(regexp.QuoteMeta(language.BlockComment.Start))
		blockEndRegex = regexp.MustCompile(regexp.QuoteMeta(language.BlockComment.End))
	}

	fmt.Printf("\n%s %s\n\n", yellow.Sprint("Preview:"), filename)

	inBlockComment := false
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		shouldDelete := false

		// Check if line should be deleted
		if inBlockComment && blockEndRegex != nil {
			shouldDelete = true
			if blockEndRegex.MatchString(line) {
				inBlockComment = false
			}
		} else if blockStartRegex != nil && blockStartRegex.MatchString(line) {
			shouldDelete = true
			if !blockEndRegex.MatchString(line) {
				inBlockComment = true
			}
		} else if lineRegex != nil && lineRegex.MatchString(line) {
			shouldDelete = true
		}
		// Don't delete empty lines - preserve original file structure

		// Print the line with appropriate formatting
		lineNumStr := gray.Sprintf("%4d", lineNum)
		if shouldDelete {
			deletedCount++
			fmt.Printf("%s %s %s\n", lineNumStr, red.Sprint("-"), red.Sprint(line))
		} else {
			keptCount++
			fmt.Printf("%s %s %s\n", lineNumStr, green.Sprint(" "), line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Printf("\n%s\n", strings.Repeat("-", 50))
	fmt.Printf("%s %d lines would be removed\n", red.Sprint("✗"), deletedCount)
	fmt.Printf("%s %d lines would be kept\n\n", green.Sprint("✓"), keptCount)

	return nil
}
