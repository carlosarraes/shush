package git

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)



func DetectRepo() (*GitStatus, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return &GitStatus{IsRepo: false}, nil
	}

	rootDir := strings.TrimSpace(string(output))
	return &GitStatus{
		IsRepo:  true,
		RootDir: rootDir,
	}, nil
}


func GetStagedChanges() ([]FileChange, error) {

	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get staged files: %v", err)
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(files) == 1 && files[0] == "" {
return []FileChange{}, nil
	}

	var changes []FileChange
	for _, file := range files {
		if file == "" {
			continue
		}


		lineRanges, err := getLineRangesFromDiff(file, true)
		if err != nil {
			return nil, fmt.Errorf("failed to get line ranges for %s: %v", file, err)
		}

		changes = append(changes, FileChange{
			Path:       file,
			Status:     StatusStaged,
			LineRanges: lineRanges,
		})
	}

	return changes, nil
}


func GetUnstagedChanges() ([]FileChange, error) {
	var changes []FileChange


	cmd := exec.Command("git", "diff", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get unstaged modified files: %v", err)
	}

	modifiedFiles := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, file := range modifiedFiles {
		if file == "" {
			continue
		}


		lineRanges, err := getLineRangesFromDiff(file, false)
		if err != nil {
			return nil, fmt.Errorf("failed to get line ranges for %s: %v", file, err)
		}

		changes = append(changes, FileChange{
			Path:       file,
			Status:     StatusUnstaged,
			LineRanges: lineRanges,
		})
	}


	cmd = exec.Command("git", "ls-files", "--others", "--exclude-standard")
	output, err = cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get untracked files: %v", err)
	}

	untrackedFiles := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, file := range untrackedFiles {
		if file == "" {
			continue
		}


		changes = append(changes, FileChange{
			Path:       file,
			Status:     StatusUntracked,
LineRanges: []LineRange{},
		})
	}

	return changes, nil
}


func getLineRangesFromDiff(file string, staged bool) ([]LineRange, error) {
	var cmd *exec.Cmd
	if staged {
		cmd = exec.Command("git", "diff", "--cached", "--unified=0", file)
	} else {
		cmd = exec.Command("git", "diff", "--unified=0", file)
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return ParseDiffUnified(string(output))
}



func ParseDiffUnified(diff string) ([]LineRange, error) {
	var ranges []LineRange


	headerRegex := regexp.MustCompile(`@@\s+-(\d+)(?:,(\d+))?\s+\+(\d+)(?:,(\d+))?\s+@@`)

	scanner := bufio.NewScanner(strings.NewReader(diff))
	for scanner.Scan() {
		line := scanner.Text()

		matches := headerRegex.FindStringSubmatch(line)
		if len(matches) > 0 {

			newStart, err := strconv.Atoi(matches[3])
			if err != nil {
				continue
			}

			var newCount int
			if matches[4] == "" {
newCount = 1
			} else {
				newCount, err = strconv.Atoi(matches[4])
				if err != nil {
					continue
				}
			}


			if newCount > 0 {
				ranges = append(ranges, LineRange{
					Start: newStart,
					End:   newStart + newCount - 1,
				})
			}
		}
	}

	return ranges, scanner.Err()
}


func IsInLineRanges(lineNum int, ranges []LineRange) bool {
	for _, r := range ranges {
		if lineNum >= r.Start && lineNum <= r.End {
			return true
		}
	}
	return false
}


func GetChangesOnly() ([]FileChange, error) {
	var allChanges []FileChange


	staged, err := GetStagedChanges()
	if err != nil {
		return nil, fmt.Errorf("failed to get staged changes: %v", err)
	}
	allChanges = append(allChanges, staged...)


	unstaged, err := GetUnstagedChanges()
	if err != nil {
		return nil, fmt.Errorf("failed to get unstaged changes: %v", err)
	}
	allChanges = append(allChanges, unstaged...)

	return allChanges, nil
}
