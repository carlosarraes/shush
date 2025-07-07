package hooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type HookScope int

const (
	ScopeUser HookScope = iota
	ScopeProject
)

type ClaudeSettings struct {
	Data map[string]interface{} `json:"-"`
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

func CreateBackup(settingsPath string) error {
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		return nil
	}

	timestamp := time.Now().Unix()
	backupPath := fmt.Sprintf("%s.backup.%d", settingsPath, timestamp)

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return fmt.Errorf("failed to read original settings for backup: %w", err)
	}

	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to create backup at %s: %w", backupPath, err)
	}

	fmt.Printf("✓ Backup created: %s\n", backupPath)
	return nil
}

func LoadSettings(path string) (*ClaudeSettings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &ClaudeSettings{
				Data: make(map[string]interface{}),
			}, nil
		}
		return nil, fmt.Errorf("failed to read settings file %s: %w", path, err)
	}

	var allSettings map[string]interface{}
	if err := json.Unmarshal(data, &allSettings); err != nil {
		return nil, fmt.Errorf("failed to parse settings file %s: %w", path, err)
	}

	if allSettings == nil {
		allSettings = make(map[string]interface{})
	}

	return &ClaudeSettings{
		Data: allSettings,
	}, nil
}

func SaveSettings(path string, settings *ClaudeSettings) error {
	if err := EnsureSettingsDirectory(path); err != nil {
		return err
	}

	data, err := json.MarshalIndent(settings.Data, "", "  ")
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
	hooksData, exists := settings.Data["hooks"]
	if !exists {
		return false
	}

	hooksMap, ok := hooksData.(map[string]interface{})
	if !ok {
		return false
	}

	postToolUseData, exists := hooksMap["PostToolUse"]
	if !exists {
		return false
	}

	postToolUseConfigs, ok := postToolUseData.([]interface{})
	if !ok {
		return false
	}

	for _, configData := range postToolUseConfigs {
		config, ok := configData.(map[string]interface{})
		if !ok {
			continue
		}

		hooksData, exists := config["hooks"]
		if !exists {
			continue
		}

		hooks, ok := hooksData.([]interface{})
		if !ok {
			continue
		}

		for _, hookData := range hooks {
			hook, ok := hookData.(map[string]interface{})
			if !ok {
				continue
			}

			command, exists := hook["command"]
			if !exists {
				continue
			}

			commandStr, ok := command.(string)
			if ok && strings.Contains(commandStr, "shush --changes-only") {
				return true
			}
		}
	}

	return false
}

func AddShushHook(settings *ClaudeSettings) error {
	if HasShushHook(settings) {
		return errors.New("shush hook already installed")
	}

	if _, exists := settings.Data["hooks"]; !exists {
		settings.Data["hooks"] = make(map[string]interface{})
	}

	hooksMap := settings.Data["hooks"].(map[string]interface{})

	shushHook := map[string]interface{}{
		"type":    "command",
		"command": "shush --changes-only",
		"timeout": 30,
	}

	if postToolUseData, exists := hooksMap["PostToolUse"]; exists {
		postToolUseConfigs, ok := postToolUseData.([]interface{})
		if !ok {
			return errors.New("invalid PostToolUse format in existing settings")
		}

		for _, configData := range postToolUseConfigs {
			config, ok := configData.(map[string]interface{})
			if !ok {
				continue
			}

			matcher, exists := config["matcher"]
			if !exists {
				continue
			}

			matcherStr, ok := matcher.(string)
			if !ok {
				continue
			}

			if matcherStr == "Write|Edit|MultiEdit" {
				hooksData, exists := config["hooks"]
				if !exists {
					config["hooks"] = []interface{}{shushHook}
					return nil
				}

				hooks, ok := hooksData.([]interface{})
				if !ok {
					return errors.New("invalid hooks format in existing settings")
				}

				config["hooks"] = append(hooks, shushHook)
				return nil
			}
		}

		newConfig := map[string]interface{}{
			"matcher": "Write|Edit|MultiEdit",
			"hooks":   []interface{}{shushHook},
		}

		hooksMap["PostToolUse"] = append(postToolUseConfigs, newConfig)
		return nil
	}

	newConfig := map[string]interface{}{
		"matcher": "Write|Edit|MultiEdit",
		"hooks":   []interface{}{shushHook},
	}

	hooksMap["PostToolUse"] = []interface{}{newConfig}
	return nil
}

func RemoveShushHook(settings *ClaudeSettings) error {
	if !HasShushHook(settings) {
		return errors.New("shush hook not found")
	}

	hooksData, exists := settings.Data["hooks"]
	if !exists {
		return errors.New("no hooks configuration found")
	}

	hooksMap, ok := hooksData.(map[string]interface{})
	if !ok {
		return errors.New("invalid hooks format")
	}

	postToolUseData, exists := hooksMap["PostToolUse"]
	if !exists {
		return errors.New("no PostToolUse hooks found")
	}

	postToolUseConfigs, ok := postToolUseData.([]interface{})
	if !ok {
		return errors.New("invalid PostToolUse format")
	}

	found := false
	var newConfigs []interface{}

	for _, configData := range postToolUseConfigs {
		config, ok := configData.(map[string]interface{})
		if !ok {
			newConfigs = append(newConfigs, configData)
			continue
		}

		hooksData, exists := config["hooks"]
		if !exists {
			newConfigs = append(newConfigs, configData)
			continue
		}

		hooks, ok := hooksData.([]interface{})
		if !ok {
			newConfigs = append(newConfigs, configData)
			continue
		}

		var newHooks []interface{}
		for _, hookData := range hooks {
			hook, ok := hookData.(map[string]interface{})
			if !ok {
				newHooks = append(newHooks, hookData)
				continue
			}

			command, exists := hook["command"]
			if !exists {
				newHooks = append(newHooks, hookData)
				continue
			}

			commandStr, ok := command.(string)
			if !ok {
				newHooks = append(newHooks, hookData)
				continue
			}

			if strings.Contains(commandStr, "shush --changes-only") {
				found = true
			} else {
				newHooks = append(newHooks, hookData)
			}
		}

		if len(newHooks) > 0 {
			newConfig := make(map[string]interface{})
			for k, v := range config {
				newConfig[k] = v
			}
			newConfig["hooks"] = newHooks
			newConfigs = append(newConfigs, newConfig)
		}
	}

	if !found {
		return errors.New("shush hook not found")
	}

	if len(newConfigs) == 0 {
		delete(hooksMap, "PostToolUse")
	} else {
		hooksMap["PostToolUse"] = newConfigs
	}

	if len(hooksMap) == 0 {
		delete(settings.Data, "hooks")
	}

	return nil
}
