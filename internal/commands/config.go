package commands

import (
	"fmt"

	"github.com/carlosarraes/shush/internal/config"
	"github.com/carlosarraes/shush/internal/types"
)

func HandleConfig(cli types.CLI) error {
	switch {
	case cli.Config:
		return ShowConfig()
	case cli.CreateConfig:
		return CreateConfig()
	}
	return nil
}

func ShowConfig() error {
	cfg, configPath, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("Shush Configuration:")
	if configPath != "" {
		fmt.Printf("Config file: %s\n", configPath)
	} else {
		fmt.Println("Config file: Using defaults (no config file found)")
	}

	fmt.Printf("\nPreserve patterns (%d):\n", len(cfg.Preserve))
	for i, pattern := range cfg.Preserve {
		fmt.Printf("  %2d. %s\n", i+1, pattern)
	}

	fmt.Println("\nConfig file search order:")
	fmt.Println("  1. .shush.toml (current directory)")
	fmt.Println("  2. .shush.toml (git repository root)")
	fmt.Println("  3. ~/.config/.shush.toml (global)")

	return nil
}

func CreateConfig() error {
	if err := config.CreateExampleConfig(); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	fmt.Println("✓ Created .shush.toml configuration file")
	fmt.Println("  Edit this file to customize which comments to preserve")
	fmt.Println("  Patterns support wildcards with * (e.g., '*IMPORTANT*')")

	return nil
}
