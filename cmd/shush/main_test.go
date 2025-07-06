package main

import (
	"testing"

	"github.com/carlosarraes/shush/internal/types"
)

func TestGitFlagValidation(t *testing.T) {
	tests := []struct {
		name          string
		cli           types.CLI
		shouldPass    bool
		expectedError string
	}{
		{
			name:       "single git flag should pass",
			cli:        types.CLI{Staged: true},
			shouldPass: true,
		},
		{
			name:       "changes-only flag should pass",
			cli:        types.CLI{ChangesOnly: true},
			shouldPass: true,
		},
		{
			name:       "unstaged flag should pass",
			cli:        types.CLI{Unstaged: true},
			shouldPass: true,
		},
		{
			name:          "multiple git flags should fail",
			cli:           types.CLI{Staged: true, Unstaged: true},
			shouldPass:    false,
			expectedError: "git flags (--changes-only, --staged, --unstaged) are mutually exclusive",
		},
		{
			name:          "all git flags should fail",
			cli:           types.CLI{ChangesOnly: true, Staged: true, Unstaged: true},
			shouldPass:    false,
			expectedError: "git flags (--changes-only, --staged, --unstaged) are mutually exclusive",
		},
		{
			name:          "git flag with path should fail",
			cli:           types.CLI{Staged: true, Path: "src/"},
			shouldPass:    false,
			expectedError: "cannot use git flags with explicit path argument",
		},
		{
			name:          "git flag with recursive should fail",
			cli:           types.CLI{Unstaged: true, Recursive: true},
			shouldPass:    false,
			expectedError: "cannot use git flags with --recursive (git handles repository scope)",
		},
		{
			name:       "git flag with compatible flags should pass",
			cli:        types.CLI{Staged: true, DryRun: true, Verbose: true, Backup: true},
			shouldPass: true,
		},
		{
			name:       "git flag with inline should pass",
			cli:        types.CLI{ChangesOnly: true, Inline: true},
			shouldPass: true,
		},
		{
			name:       "git flag with block should pass",
			cli:        types.CLI{Unstaged: true, Block: true},
			shouldPass: true,
		},
		{
			name:       "non-git mode with path should pass",
			cli:        types.CLI{Path: "src/", Recursive: true},
			shouldPass: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGitFlags(tt.cli)

			if tt.shouldPass {
				if err != nil {
					t.Errorf("validateGitFlags() should pass but got error: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("validateGitFlags() should fail but passed")
				} else if err.Error() != tt.expectedError {
					t.Errorf("validateGitFlags() error = %v, want %v", err.Error(), tt.expectedError)
				}
			}
		})
	}
}

func validateGitFlags(cli types.CLI) error {

	gitFlags := []bool{cli.ChangesOnly, cli.Staged, cli.Unstaged}
	gitFlagCount := 0
	for _, flag := range gitFlags {
		if flag {
			gitFlagCount++
		}
	}

	if gitFlagCount > 1 {
		return &GitFlagError{"git flags (--changes-only, --staged, --unstaged) are mutually exclusive"}
	}

	if gitFlagCount > 0 && cli.Path != "" {
		return &GitFlagError{"cannot use git flags with explicit path argument"}
	}

	if gitFlagCount > 0 && cli.Recursive {
		return &GitFlagError{"cannot use git flags with --recursive (git handles repository scope)"}
	}

	return nil
}

type GitFlagError struct {
	message string
}

func (e *GitFlagError) Error() string {
	return e.message
}
