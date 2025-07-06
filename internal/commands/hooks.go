package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/carlosarraes/shush/internal/hooks"
	"github.com/carlosarraes/shush/internal/types"
)

func HandleHooks(cli types.CLI) error {
	switch {
	case cli.InstallHook:
		return InstallHooks(cli.HookScope)
	case cli.UninstallHook:
		return UninstallHooks(cli.HookScope)
	case cli.ListHooks:
		return ListHooks()
	case cli.HookStatus:
		return ShowHooksStatus()
	}
	return nil
}

func InstallHooks(scope string) error {
	hookScope := hooks.ScopeUser
	if scope == "project" {
		hookScope = hooks.ScopeProject
	}

	userPath, _ := hooks.GetSettingsPath(hooks.ScopeUser)
	projectPath, _ := hooks.GetSettingsPath(hooks.ScopeProject)

	userSettings, userErr := hooks.LoadSettings(userPath)
	projectSettings, projectErr := hooks.LoadSettings(projectPath)

	userHasShush := userErr == nil && hooks.HasShushHook(userSettings)
	projectHasShush := projectErr == nil && hooks.HasShushHook(projectSettings)

	if hookScope == hooks.ScopeUser {
		if userHasShush {
			return fmt.Errorf("shush hook already installed for user-wide scope at %s", userPath)
		}
		if projectHasShush {
			fmt.Printf("ℹ  Project-specific hooks found at %s\n", projectPath)
			fmt.Println("   User-wide hooks will take precedence and override project hooks")
		}
	} else {
		if userHasShush {
			return fmt.Errorf("user-wide hooks already installed at %s\n"+
				"   User-wide hooks cover all projects including this one.\n"+
				"   Use --uninstall-hook first if you want project-specific hooks instead", userPath)
		}
		if projectHasShush {
			return fmt.Errorf("shush hook already installed for project scope at %s", projectPath)
		}
	}

	path, err := hooks.GetSettingsPath(hookScope)
	if err != nil {
		return fmt.Errorf("failed to get settings path: %w", err)
	}

	settings, err := hooks.LoadSettings(path)
	if err != nil {
		return fmt.Errorf("failed to load settings: %w", err)
	}

	if err := hooks.AddShushHook(settings); err != nil {
		return fmt.Errorf("failed to add shush hook: %w", err)
	}

	if err := hooks.SaveSettings(path, settings); err != nil {
		return fmt.Errorf("failed to save settings: %w", err)
	}

	scopeName := "user-wide"
	if hookScope == hooks.ScopeProject {
		scopeName = "project"
	}
	fmt.Printf("✓ Hooks installed for %s scope at %s\n", scopeName, path)
	fmt.Println("  Auto-cleanup will run after Claude modifies files")
	return nil
}

func UninstallHooks(scope string) error {
	hookScope := hooks.ScopeUser
	if scope == "project" {
		hookScope = hooks.ScopeProject
	}

	path, err := hooks.GetSettingsPath(hookScope)
	if err != nil {
		return fmt.Errorf("failed to get settings path: %w", err)
	}

	settings, err := hooks.LoadSettings(path)
	if err != nil {
		return fmt.Errorf("failed to load settings: %w", err)
	}

	scopeName := "user-wide"
	if hookScope == hooks.ScopeProject {
		scopeName = "project"
	}

	if !hooks.HasShushHook(settings) {
		return fmt.Errorf("shush hook not found for %s scope at %s", scopeName, path)
	}

	if _, err := exec.LookPath("jq"); err == nil {
		if err := uninstallWithJQ(path, scopeName); err == nil {
			return nil
		}
		fmt.Printf("⚠️  jq removal failed, using Go implementation\n")
	} else {
		fmt.Printf("⚠️  jq not available. For surgical removal, install jq and run again.\n")
		fmt.Printf("   Manual removal: edit %s and remove shush hook entries\n", path)
	}

	if err := hooks.RemoveShushHook(settings); err != nil {
		return fmt.Errorf("failed to remove shush hook: %w", err)
	}

	if err := hooks.SaveSettings(path, settings); err != nil {
		return fmt.Errorf("failed to save settings: %w", err)
	}

	fmt.Printf("✓ Hooks uninstalled for %s scope at %s\n", scopeName, path)
	return nil
}

func uninstallWithJQ(path, scopeName string) error {
	jqFilter := `
		.hooks.PostToolUse |= (
			map(
				.hooks |= map(select(.command != "shush --changes-only"))
			) |
			map(select(.hooks | length > 0))
		) |
		if .hooks.PostToolUse | length == 0 then
			.hooks |= del(.PostToolUse)
		else
			.
		end
	`

	cmd := exec.Command("jq", jqFilter, path)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("jq command failed: %w", err)
	}

	if err := os.WriteFile(path, output, 0644); err != nil {
		return fmt.Errorf("failed to write updated settings: %w", err)
	}

	fmt.Printf("✓ Hooks uninstalled for %s scope at %s (using jq)\n", scopeName, path)
	return nil
}

func ShowHooksStatus() error {
	userPath, _ := hooks.GetSettingsPath(hooks.ScopeUser)
	projectPath, _ := hooks.GetSettingsPath(hooks.ScopeProject)

	userSettings, userErr := hooks.LoadSettings(userPath)
	projectSettings, projectErr := hooks.LoadSettings(projectPath)

	userHasShush := userErr == nil && hooks.HasShushHook(userSettings)
	projectHasShush := projectErr == nil && hooks.HasShushHook(projectSettings)

	fmt.Println("Shush Hooks Status:")
	fmt.Printf("User-wide (%s): %s\n", userPath, getHookStatus(userSettings, userErr))
	fmt.Printf("Project (%s): %s\n", projectPath, getHookStatus(projectSettings, projectErr))

	if userHasShush && projectHasShush {
		fmt.Println("\n⚠️  Warning: Both user-wide and project hooks are installed")
		fmt.Println("   This will cause shush to run twice on every file modification")
		fmt.Println("   Consider removing project hooks: shush --uninstall-hook -s project")
	} else if userHasShush {
		fmt.Println("\n✓ User-wide hooks will handle all projects including this one")
	} else if projectHasShush {
		fmt.Println("\n✓ Project-specific hooks active for this project only")
	}

	return nil
}

func getHookStatus(settings *hooks.ClaudeSettings, err error) string {
	if err != nil {
		return "Not configured"
	}
	if hooks.HasShushHook(settings) {
		return "✓ Installed"
	}
	return "Not installed"
}

func ListHooks() error {
	userPath, _ := hooks.GetSettingsPath(hooks.ScopeUser)
	projectPath, _ := hooks.GetSettingsPath(hooks.ScopeProject)

	fmt.Println("Claude Code Hooks Configuration:")

	fmt.Println("User-wide:")
	if err := listHooksForPath(userPath); err != nil {
		fmt.Println("  No hooks configured")
	}

	fmt.Println("Project:")
	if err := listHooksForPath(projectPath); err != nil {
		fmt.Println("  No hooks configured")
	}

	return nil
}

func listHooksForPath(path string) error {
	settings, err := hooks.LoadSettings(path)
	if err != nil {
		return err
	}

	if len(settings.Hooks) == 0 {
		return fmt.Errorf("no hooks")
	}

	hasShush := hooks.HasShushHook(settings)

	for event, configs := range settings.Hooks {
		for _, config := range configs {
			for _, hook := range config.Hooks {
				if hook.Command == "shush --changes-only" {
					fmt.Printf("  ✓ shush --changes-only (%s: %s)\n", event, config.Matcher)
				} else {
					fmt.Printf("  ✓ %s (%s: %s)\n", hook.Command, event, config.Matcher)
				}
			}
		}
	}

	if !hasShush {
		fmt.Println("  (No shush hooks found)")
	}

	return nil
}
