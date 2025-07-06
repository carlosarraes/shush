package processor

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/carlosarraes/shush/internal/config"
	"github.com/carlosarraes/shush/internal/git"
	"github.com/carlosarraes/shush/internal/types"
	"github.com/fatih/color"
)

type GitTotals struct {
	FilesProcessed int
	TotalChanged   int
	TotalKept      int
	TotalPreserved int
}

func (p *Processor) processGitChanges() error {

	gitStatus, err := git.DetectRepo()
	if err != nil {
		return fmt.Errorf("failed to detect git repository: %v", err)
	}

	if !gitStatus.IsRepo {
		return fmt.Errorf("not in a git repository. Use without git flags to process files normally")
	}

	if p.cli.Verbose {
		fmt.Printf("Git repository detected: %s\n", gitStatus.RootDir)
	}

	var changes []git.FileChange
	switch {
	case p.cli.ChangesOnly:
		changes, err = git.GetChangesOnly()
		if err != nil {
			return fmt.Errorf("failed to get git changes: %v", err)
		}
	case p.cli.Staged:
		changes, err = git.GetStagedChanges()
		if err != nil {
			return fmt.Errorf("failed to get staged changes: %v", err)
		}
	case p.cli.Unstaged:
		changes, err = git.GetUnstagedChanges()
		if err != nil {
			return fmt.Errorf("failed to get unstaged changes: %v", err)
		}
	}

	if len(changes) == 0 {
		fmt.Println("No changes found to process")
		return nil
	}


	supportedChanges := make([]git.FileChange, 0, len(changes))
	for _, change := range changes {
		if IsSupportedFile(change.Path) {
			supportedChanges = append(supportedChanges, change)
		} else if p.cli.Verbose {
			fmt.Printf("Skipping unsupported file: %s\n", change.Path)
		}
	}

	if len(supportedChanges) == 0 {
		fmt.Println("No supported files found to process")
		return nil
	}

	if p.cli.Verbose {
		fmt.Printf("Found %d supported files with changes to process\n", len(supportedChanges))
	}

	totals := &GitTotals{}

	for _, change := range supportedChanges {

		if p.cli.Verbose {
			fmt.Printf("Processing: %s\n", change.Path)
		}

		if p.cli.DryRun {
			if err := p.showGitPreviewWithTotals(change.Path, change.LineRanges, totals); err != nil {
				fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", change.Path, err)
				continue
			}
		} else {
			if err := p.processFileWithLineRanges(change.Path, change.LineRanges); err != nil {
				fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", change.Path, err)
				continue
			}
		}
		totals.FilesProcessed++
	}

	if p.cli.DryRun && totals.FilesProcessed > 0 {
		p.showGitTotals(totals)
	}

	return nil
}

func (p *Processor) processFileWithLineRanges(filename string, lineRanges []git.LineRange) error {
	language, err := DetectLanguage(filename)
	if err != nil {
		return err
	}


	cfg, _, err := config.Load()
	if err != nil && p.cli.Verbose {
		fmt.Printf("Warning: failed to load config, using defaults: %v\n", err)
		cfg = config.Default()
	}

	if p.cli.Verbose {
		fmt.Printf("Processing %s (language: %s)\n", filename, GetLanguageName(filename))
		if len(lineRanges) == 0 {
			fmt.Printf("Processing entire file (untracked)\n")
		} else {
			fmt.Printf("Processing %d line ranges\n", len(lineRanges))
		}
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

	processEntireFile := len(lineRanges) == 0
	modified := false

	for i, line := range lines {
		lineNum := i + 1
		shouldProcess := processEntireFile || git.IsInLineRanges(lineNum, lineRanges)

		if shouldProcess {
			newLine := p.removeCommentsFromLine(line, language, cfg)
			if newLine != line {
				lines[i] = newLine
				modified = true
			}
		}
	}

	if modified {
		outFile, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer outFile.Close()

		for _, line := range lines {
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

func (p *Processor) removeCommentsFromLine(line string, language types.Language, cfg *config.Config) string {
	result := line
	hasChanges := false

	if !p.cli.Block && language.LineComment != "" {

		if idx := strings.Index(result, language.LineComment); idx != -1 {

			comment := strings.TrimSpace(result[idx:])
			

			if cfg.ShouldPreserveComment(comment) {
return line
			}
			
			result = result[:idx]
			hasChanges = true
		}
	}

	if !p.cli.Inline && language.BlockComment != nil {

		startComment := language.BlockComment.Start
		endComment := language.BlockComment.End

		for {
			startIdx := strings.Index(result, startComment)
			if startIdx == -1 {
				break
			}

			endIdx := strings.Index(result[startIdx:], endComment)
			if endIdx == -1 {

				comment := strings.TrimSpace(result[startIdx:])
				if cfg.ShouldPreserveComment(comment) {
return line
				}
				
				result = result[:startIdx]
				hasChanges = true
				break
			}


			blockComment := strings.TrimSpace(result[startIdx:startIdx+endIdx+len(endComment)])
			if cfg.ShouldPreserveComment(blockComment) {
return line
			}

			endIdx += startIdx + len(endComment)
			result = result[:startIdx] + result[endIdx:]
			hasChanges = true
		}
	}

	if hasChanges {
		result = strings.TrimSpace(result)
	}

	return result
}

func (p *Processor) showGitPreviewWithTotals(filename string, lineRanges []git.LineRange, totals *GitTotals) error {
	language, err := DetectLanguage(filename)
	if err != nil {
		return err
	}


	cfg, _, err := config.Load()
	if err != nil {
		cfg = config.Default()
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	red := color.New(color.FgRed, color.CrossedOut)
	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)
	gray := color.New(color.FgHiBlack)
	yellow := color.New(color.FgYellow)

	fmt.Printf("\n%s %s\n", yellow.Sprint("Git Preview:"), filename)
	if len(lineRanges) == 0 {
		fmt.Printf("%s\n", blue.Sprint("Processing entire file (untracked)"))
	} else {
		fmt.Printf("%s %d line ranges\n", blue.Sprint("Processing"), len(lineRanges))
	}
	fmt.Println()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	keptCount := 0
	changedCount := 0
	preservedCount := 0

	processEntireFile := len(lineRanges) == 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		shouldProcess := processEntireFile || git.IsInLineRanges(lineNum, lineRanges)

		lineNumStr := gray.Sprintf("%4d", lineNum)

		if shouldProcess {
			newLine := p.removeCommentsFromLine(line, language, cfg)
			if newLine != line {
				changedCount++
				fmt.Printf("%s %s %s\n", lineNumStr, red.Sprint("~"), red.Sprint(line))
				if strings.TrimSpace(newLine) != "" {
					fmt.Printf("%s %s %s\n", lineNumStr, green.Sprint("+"), green.Sprint(newLine))
				}
			} else {
				// Check if this line has comments that were preserved
				hasComment := p.lineHasComment(line, language)
				if hasComment {
					preservedCount++
					cyan := color.New(color.FgCyan)
					fmt.Printf("%s %s %s\n", lineNumStr, cyan.Sprint("P"), line)
				} else {
					keptCount++
					fmt.Printf("%s %s %s\n", lineNumStr, green.Sprint(" "), line)
				}
			}
		} else {
			keptCount++
			fmt.Printf("%s %s %s\n", lineNumStr, gray.Sprint(" "), gray.Sprint(line))
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Printf("\n%s\n", strings.Repeat("-", 50))
	fmt.Printf("%s %d lines would be changed\n", yellow.Sprint("~"), changedCount)
	fmt.Printf("%s %d lines would be kept\n", green.Sprint("✓"), keptCount)
	if preservedCount > 0 {
		cyan := color.New(color.FgCyan)
		fmt.Printf("%s %d comments would be preserved\n", cyan.Sprint("P"), preservedCount)
	}
	fmt.Println()

	totals.TotalChanged += changedCount
	totals.TotalKept += keptCount
	totals.TotalPreserved += preservedCount

	return nil
}

// lineHasComment checks if a line contains comments
func (p *Processor) lineHasComment(line string, language types.Language) bool {
	// Check for line comments
	if language.LineComment != "" && strings.Contains(line, language.LineComment) {
		return true
	}
	
	// Check for block comments
	if language.BlockComment != nil {
		if strings.Contains(line, language.BlockComment.Start) || strings.Contains(line, language.BlockComment.End) {
			return true
		}
	}
	
	return false
}

func (p *Processor) showGitTotals(totals *GitTotals) {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)
	cyan := color.New(color.FgCyan)

	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Printf("%s\n", blue.Sprint("GIT PROCESSING TOTALS"))
	fmt.Printf("%s\n", strings.Repeat("=", 60))
	fmt.Printf("%s %d files processed\n", blue.Sprint("📁"), totals.FilesProcessed)
	fmt.Printf("%s %d lines would be changed\n", yellow.Sprint("~"), totals.TotalChanged)
	fmt.Printf("%s %d lines would be kept\n", green.Sprint("✓"), totals.TotalKept)
	if totals.TotalPreserved > 0 {
		fmt.Printf("%s %d comments would be preserved\n", cyan.Sprint("P"), totals.TotalPreserved)
	}
	fmt.Printf("%s\n", strings.Repeat("=", 60))
}
