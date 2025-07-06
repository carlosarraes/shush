package hooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)







type HookScope int

const (
	ScopeUser HookScope = iota
	ScopeProject
)


type ClaudeSettings struct {
	Hooks map[string][]EventConfig `json:"hooks,omitempty"`

}


type EventConfig struct {
	Matcher string      `json:"matcher"`
	Hooks   []HookEntry `json:"hooks"`
}


type HookEntry struct {
	Type    string `json:"type"`
	Command string `json:"command"`
	Timeout int    `json:"timeout,omitempty"`
}


func GetSettingsPath(scope HookScope) (string, error) {
	switch scope {
	case ScopeUser:
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		return filepath.Join(home, ".claude", "settings.json"), nil
	case ScopeProject:
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %w", err)
		}
		return filepath.Join(cwd, ".claude", "settings.json"), nil
	default:
		return "", errors.New("invalid hook scope")
	}
}


func EnsureSettingsDirectory(settingsPath string) error {
	dir := filepath.Dir(settingsPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create settings directory %s: %w", dir, err)
	}
	return nil
}


func LoadSettings(path string) (*ClaudeSettings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {

			return &ClaudeSettings{
				Hooks: make(map[string][]EventConfig),
			}, nil
		}
		return nil, fmt.Errorf("failed to read settings file %s: %w", path, err)
	}

	var settings ClaudeSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse settings file %s: %w", path, err)
	}


	if settings.Hooks == nil {
		settings.Hooks = make(map[string][]EventConfig)
	}

	return &settings, nil
}


func SaveSettings(path string, settings *ClaudeSettings) error {
	if err := EnsureSettingsDirectory(path); err != nil {
		return err
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings file %s: %w", path, err)
	}

	return nil
}


func CreateShushHookEntry() HookEntry {
	return HookEntry{
		Type:    "command",
		Command: "shush --changes-only",
		Timeout: 30,
	}
}


func CreateShushEventConfig() EventConfig {
	return EventConfig{
		Matcher: "Write|Edit|MultiEdit",
		Hooks:   []HookEntry{CreateShushHookEntry()},
	}
}


func HasShushHook(settings *ClaudeSettings) bool {
	if settings.Hooks == nil {
		return false
	}

	postToolUseConfigs, exists := settings.Hooks["PostToolUse"]
	if !exists {
		return false
	}

	for _, config := range postToolUseConfigs {
		for _, hook := range config.Hooks {
			if strings.Contains(hook.Command, "shush --changes-only") {
				return true
			}
		}
	}

	return false
}


func AddShushHook(settings *ClaudeSettings) error {
	if settings.Hooks == nil {
		settings.Hooks = make(map[string][]EventConfig)
	}

	if HasShushHook(settings) {
		return errors.New("shush hook already installed")
	}


	postToolUseConfigs, exists := settings.Hooks["PostToolUse"]
	if !exists {

		settings.Hooks["PostToolUse"] = []EventConfig{CreateShushEventConfig()}
		return nil
	}


	for i, config := range postToolUseConfigs {
		if config.Matcher == "Write|Edit|MultiEdit" || config.Matcher == "" {

			postToolUseConfigs[i].Hooks = append(config.Hooks, CreateShushHookEntry())
			return nil
		}
	}


	settings.Hooks["PostToolUse"] = append(postToolUseConfigs, CreateShushEventConfig())
	return nil
}


func RemoveShushHook(settings *ClaudeSettings) error {
	if settings.Hooks == nil {
		return errors.New("no hooks configuration found")
	}

	postToolUseConfigs, exists := settings.Hooks["PostToolUse"]
	if !exists {
		return errors.New("shush hook not found")
	}

	found := false
	for i, config := range postToolUseConfigs {
		newHooks := make([]HookEntry, 0, len(config.Hooks))
		for _, hook := range config.Hooks {
			if !strings.Contains(hook.Command, "shush --changes-only") {
				newHooks = append(newHooks, hook)
			} else {
				found = true
			}
		}
		postToolUseConfigs[i].Hooks = newHooks
	}

	if !found {
		return errors.New("shush hook not found")
	}


	filteredConfigs := make([]EventConfig, 0, len(postToolUseConfigs))
	for _, config := range postToolUseConfigs {
		if len(config.Hooks) > 0 {
			filteredConfigs = append(filteredConfigs, config)
		}
	}

	if len(filteredConfigs) == 0 {
		delete(settings.Hooks, "PostToolUse")
	} else {
		settings.Hooks["PostToolUse"] = filteredConfigs
	}

	return nil
}
