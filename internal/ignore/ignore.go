package ignore

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Ignore struct {
	patterns []pattern
}

type pattern struct {
	original    string
	normalized  string
	isNegation  bool
	isDirectory bool
	source      string
}

func New() *Ignore {
	return &Ignore{
		patterns: make([]pattern, 0),
	}
}

func Load() (*Ignore, error) {
	ignore := New()

	if gitignorePath := getGitIgnorePath(); gitignorePath != "" {
		if err := ignore.loadFromFileWithSource(gitignorePath, "gitignore"); err != nil && !os.IsNotExist(err) {
		}
	}

	if globalPath := getGlobalIgnorePath(); globalPath != "" {
		if err := ignore.loadFromFileWithSource(globalPath, "global"); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}

	if projectPath := getProjectIgnorePath(); projectPath != "" {
		if err := ignore.loadFromFileWithSource(projectPath, "project"); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}

	return ignore, nil
}

func (i *Ignore) loadFromFile(path string) error {
	return i.loadFromFileWithSource(path, "")
}

func (i *Ignore) loadFromFileWithSource(path string, source string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		i.addPatternWithSource(line, source)
	}

	return scanner.Err()
}

func (i *Ignore) addPattern(patternStr string) {
	i.addPatternWithSource(patternStr, "")
}

func (i *Ignore) addPatternWithSource(patternStr string, source string) {
	p := pattern{
		original: patternStr,
		source:   source,
	}

	if strings.HasPrefix(patternStr, "!") {
		p.isNegation = true
		patternStr = patternStr[1:]
	}

	if strings.HasSuffix(patternStr, "/") {
		p.isDirectory = true
		patternStr = patternStr[:len(patternStr)-1]
	}

	p.normalized = patternStr

	i.patterns = append(i.patterns, p)
}

func (i *Ignore) IsIgnored(filePath string) bool {
	relPath := filepath.Clean(filePath)
	
	if strings.HasPrefix(relPath, "./") {
		relPath = relPath[2:]
	}

	matched := false

	for _, p := range i.patterns {
		if i.matchesPattern(relPath, p) {
			matched = !p.isNegation
		}
	}

	return matched
}

func (i *Ignore) matchesPattern(filePath string, p pattern) bool {
	if p.isDirectory {
		if info, err := os.Stat(filePath); err != nil || !info.IsDir() {
			return false
		}
	}

	return i.globMatch(p.normalized, filePath)
}

func (i *Ignore) globMatch(pattern, path string) bool {
	if pattern == path {
		return true
	}

	if strings.Contains(pattern, "*") {
		return i.wildcardMatch(pattern, path)
	}

	if strings.HasSuffix(pattern, "/") {
		dirPattern := pattern[:len(pattern)-1]
		return strings.HasPrefix(path, dirPattern+"/") || path == dirPattern
	}

	if strings.HasPrefix(path, pattern+"/") {
		return true
	}

	if filepath.Base(path) == pattern {
		return true
	}

	return false
}

func (i *Ignore) wildcardMatch(pattern, path string) bool {
	
	parts := strings.Split(pattern, "*")
	
	if len(parts) == 1 {
		return pattern == path
	}

	if len(parts[0]) > 0 && !strings.HasPrefix(path, parts[0]) {
		return false
	}

	if len(parts[len(parts)-1]) > 0 && !strings.HasSuffix(path, parts[len(parts)-1]) {
		return false
	}

	currentPath := path
	for i := 0; i < len(parts); i++ {
		part := parts[i]
		if part == "" {
			continue
		}

		if i == 0 {
			currentPath = currentPath[len(part):]
			continue
		}

		if i == len(parts)-1 {
			break
		}

		index := strings.Index(currentPath, part)
		if index == -1 {
			return false
		}
		currentPath = currentPath[index+len(part):]
	}

	return true
}

func getGitIgnorePath() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	localGitIgnore := filepath.Join(cwd, ".gitignore")
	if fileExists(localGitIgnore) {
		return localGitIgnore
	}

	if gitRoot := findGitRoot(cwd); gitRoot != "" {
		gitRootIgnore := filepath.Join(gitRoot, ".gitignore")
		if fileExists(gitRootIgnore) && gitRootIgnore != localGitIgnore {
			return gitRootIgnore
		}
	}

	return ""
}

func getGlobalIgnorePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", ".shushignore")
}

func getProjectIgnorePath() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	localIgnore := filepath.Join(cwd, ".shushignore")
	if fileExists(localIgnore) {
		return localIgnore
	}

	if gitRoot := findGitRoot(cwd); gitRoot != "" {
		gitIgnore := filepath.Join(gitRoot, ".shushignore")
		if fileExists(gitIgnore) && gitIgnore != localIgnore {
			return gitIgnore
		}
	}

	return ""
}

func findGitRoot(startDir string) string {
	dir := startDir
	for {
		gitDir := filepath.Join(dir, ".git")
		if fileExists(gitDir) {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
