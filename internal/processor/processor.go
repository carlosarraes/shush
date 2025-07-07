package processor

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/carlosarraes/shush/internal/config"
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

	return p.processFileInMemory(filename, language)
}

func (p *Processor) processFileInMemory(filename string, language types.Language) error {

	cfg, _, err := config.Load()
	if err != nil && p.cli.Verbose {
		fmt.Printf("Warning: failed to load config, using defaults: %v\n", err)
		cfg = config.Default()
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if p.cli.Backup {
		if err := p.createBackup(filename); err != nil {
			return fmt.Errorf("failed to create backup: %v", err)
		}
		if p.cli.Verbose {
			fmt.Printf("✓ Backup created: %s.bak\n", filename)
		}
	}

	modified := false
	var processedLines []string

	for _, line := range lines {
		newLine := p.removeCommentsFromLine(line, language, cfg)
		if newLine != line {
			modified = true

			if newLine != "" {
				processedLines = append(processedLines, newLine)
			}

		} else {
			processedLines = append(processedLines, line)
		}
	}

	if modified {
		outFile, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer outFile.Close()

		for _, line := range processedLines {
			if _, err := fmt.Fprintln(outFile, line); err != nil {
				return err
			}
		}

		if p.cli.Verbose {
			fmt.Printf("✓ Comments removed from %s\n", filename)
		}
	} else if p.cli.Verbose {
		fmt.Printf("✓ No changes made to %s\n", filename)
	}

	return nil
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

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func (p *Processor) showPreview(filename string, language types.Language) error {

	cfg, _, err := config.Load()
	if err != nil {
		cfg = config.Default()
	}

	contextLines := cfg.ContextLines
	if p.cli.ContextLines >= 0 {
		contextLines = p.cli.ContextLines
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	red := color.New(color.FgRed, color.CrossedOut)
	green := color.New(color.FgGreen)
	gray := color.New(color.FgHiBlack)
	yellow := color.New(color.FgYellow)
	cyan := color.New(color.FgCyan)
	dimGray := color.New(color.FgWhite, color.Faint)

	fmt.Printf("\n%s %s\n\n", yellow.Sprint("Preview:"), filename)

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	changedCount := 0
	keptCount := 0
	preservedCount := 0

	var changes []changeInfo

	for i, line := range lines {
		lineNum := i + 1
		newLine := p.removeCommentsFromLine(line, language, cfg)

		if newLine != line {
			changedCount++
			if newLine == "" {
				changes = append(changes, changeInfo{lineNum, line, newLine, "removed"})
			} else {
				changes = append(changes, changeInfo{lineNum, line, newLine, "modified"})
			}
		} else {
			hasComment := p.lineHasComment(line, language)
			if hasComment {
				preservedCount++
				changes = append(changes, changeInfo{lineNum, line, line, "preserved"})
			} else {
				keptCount++
			}
		}
	}

	if len(changes) == 0 {
		fmt.Printf("%s No comments found to remove\n", gray.Sprint("→"))
		fmt.Println()
		return nil
	}

	if contextLines > 0 {
		p.displayChangesWithContext(lines, changes, contextLines, red, green, nil, gray, dimGray)
	} else {
		for _, change := range changes {
			lineNumStr := gray.Sprintf("%4d", change.lineNum)
			switch change.changeType {
			case "removed":
				fmt.Printf("%s %s %s\n", lineNumStr, red.Sprint("-"), red.Sprint(change.oldLine))
			case "modified":
				fmt.Printf("%s %s %s\n", lineNumStr, red.Sprint("~"), red.Sprint(change.oldLine))
				fmt.Printf("%s %s %s\n", lineNumStr, green.Sprint("+"), green.Sprint(change.newLine))
			case "preserved":
				fmt.Printf("%s %s %s\n", lineNumStr, cyan.Sprint("P"), change.oldLine)
			}
		}
	}

	fmt.Printf("\n%s\n", strings.Repeat("-", 50))
	fmt.Printf("%s %d lines would be changed\n", yellow.Sprint("~"), changedCount)
	fmt.Printf("%s %d lines would be kept\n", green.Sprint("✓"), keptCount)
	if preservedCount > 0 {
		fmt.Printf("%s %d comments would be preserved\n", cyan.Sprint("P"), preservedCount)
	}
	fmt.Println()

	return nil
}
