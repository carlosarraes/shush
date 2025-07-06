package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Preserve []string `toml:"preserve"`
}

func Default() *Config {
	return &Config{
		Preserve: []string{
			"TODO:",
			"FIXME:",
			"HACK:",
			"XXX:",
			"@ts-ignore",
			"@ts-expect-error",
			"eslint-",
			"prettier-ignore",
			"pylint:",
			"mypy:",
			"type: ignore",
		},
	}
}

func Load() (*Config, string, error) {

	if configPath, found := findProjectConfig(); found {
		config, err := loadFromFile(configPath)
		if err != nil {
			return Default(), configPath, err
		}
		return config, configPath, nil
	}

	if configPath, found := findGlobalConfig(); found {
		config, err := loadFromFile(configPath)
		if err != nil {
			return Default(), configPath, err
		}
		return config, configPath, nil
	}

	return Default(), "", nil
}

func findProjectConfig() (string, bool) {

	cwd, err := os.Getwd()
	if err != nil {
		return "", false
	}

	localConfig := filepath.Join(cwd, ".shush.toml")
	if fileExists(localConfig) {
		return localConfig, true
	}

	if gitRoot := findGitRoot(cwd); gitRoot != "" {
		gitConfig := filepath.Join(gitRoot, ".shush.toml")
		if fileExists(gitConfig) && gitConfig != localConfig {
			return gitConfig, true
		}
	}

	return "", false
}

func findGlobalConfig() (string, bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", false
	}

	configPath := filepath.Join(home, ".config", ".shush.toml")
	if fileExists(configPath) {
		return configPath, true
	}

	return "", false
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

func loadFromFile(path string) (*Config, error) {
	config := &Config{}
	if _, err := toml.DecodeFile(path, config); err != nil {
		return nil, err
	}

	if len(config.Preserve) == 0 {
		config.Preserve = Default().Preserve
	}

	return config, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (c *Config) ShouldPreserveComment(comment string) bool {
	comment = strings.TrimSpace(comment)

	for _, pattern := range c.Preserve {
		if matchesPattern(comment, pattern) {
			return true
		}
	}

	return false
}

func matchesPattern(text, pattern string) bool {

	if strings.Contains(pattern, "*") {
		return matchesWildcard(text, pattern)
	}

	return strings.Contains(text, pattern)
}

func matchesWildcard(text, pattern string) bool {

	parts := strings.Split(pattern, "*")

	if !strings.HasPrefix(pattern, "*") {
		if !strings.HasPrefix(text, parts[0]) {
			return false
		}
		text = text[len(parts[0]):]
		parts = parts[1:]
	} else {
		parts = parts[1:]
	}

	if !strings.HasSuffix(pattern, "*") && len(parts) > 0 {
		lastPart := parts[len(parts)-1]
		if !strings.HasSuffix(text, lastPart) {
			return false
		}
		text = text[:len(text)-len(lastPart)]
		parts = parts[:len(parts)-1]
	}

	for _, part := range parts {
		if part == "" {
			continue
		}

		index := strings.Index(text, part)
		if index == -1 {
			return false
		}

		text = text[index+len(part):]
	}

	return true
}

func CreateExampleConfig() error {
	content := `# Shush Configuration
# Patterns to preserve in comments (supports wildcards with *)
preserve = [
    "TODO:",
    "FIXME:", 
    "HACK:",
    "XXX:",
    "@ts-ignore",
    "@ts-expect-error", 
    "eslint-",
    "prettier-ignore",
    "pylint:",
    "mypy:",
    "type: ignore",
    "*IMPORTANT*",   # Example wildcard: preserves any comment containing IMPORTANT
    "*DEBUG*",       # Example wildcard: preserves any comment containing DEBUG
]
`

	return os.WriteFile(".shush.toml", []byte(content), 0644)
}
