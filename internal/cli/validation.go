package cli

import (
	"fmt"

	"github.com/carlosarraes/shush/internal/types"
)

func ValidateFlags(cli types.CLI) error {
	hookCommands := []bool{
		cli.InstallHook,
		cli.UninstallHook,
		cli.ListHooks,
		cli.HookStatus,
		cli.Config,
		cli.CreateConfig,
	}
	hookFlagCount := 0
	for _, flag := range hookCommands {
		if flag {
			hookFlagCount++
		}
	}

	if hookFlagCount > 1 {
		return fmt.Errorf("utility commands (--install-hook, --uninstall-hook, --list-hooks, --hook-status, --config, --create-config) are mutually exclusive")
	}

	if hookFlagCount > 0 {
		if cli.Recursive || cli.Inline || cli.Block || cli.DryRun || cli.Backup || cli.Verbose {
			return fmt.Errorf("hook commands cannot be combined with processing flags")
		}

		gitFlags := []bool{cli.ChangesOnly, cli.Staged, cli.Unstaged}
		for _, flag := range gitFlags {
			if flag {
				return fmt.Errorf("hook commands cannot be combined with git flags")
			}
		}
		return nil
	}

	gitFlags := []bool{cli.ChangesOnly, cli.Staged, cli.Unstaged}
	gitFlagCount := 0
	for _, flag := range gitFlags {
		if flag {
			gitFlagCount++
		}
	}

	if gitFlagCount > 1 {
		return fmt.Errorf("git flags (--changes-only, --staged, --unstaged) are mutually exclusive")
	}

	if gitFlagCount > 0 && cli.Path != "" {
		return fmt.Errorf("cannot use git flags with explicit path argument")
	}

	if gitFlagCount > 0 && cli.Recursive {
		return fmt.Errorf("cannot use git flags with --recursive (git handles repository scope)")
	}

	if gitFlagCount == 0 && cli.Path == "" {
		return fmt.Errorf("path argument is required")
	}

	if cli.Inline && cli.Block {
		return fmt.Errorf("--inline and --block flags are mutually exclusive")
	}

	return nil
}
