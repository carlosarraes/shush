package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/carlosarraes/shush/internal/cli"
	"github.com/carlosarraes/shush/internal/commands"
	"github.com/carlosarraes/shush/internal/guide"
	"github.com/carlosarraes/shush/internal/processor"
	"github.com/carlosarraes/shush/internal/types"
)

var version = "0.3.3"

func main() {
	var cliArgs types.CLI
	kong.Parse(&cliArgs,
		kong.Description("SHUSH: Sloppily Hushing Unwanted Source-code Heavy"),
		kong.Vars{"version": version})

	if cliArgs.LLM {
		guide.Show()
		return
	}

	if err := cli.ValidateFlags(cliArgs); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	hookCommands := []bool{
		cliArgs.InstallHook,
		cliArgs.UninstallHook,
		cliArgs.ListHooks,
		cliArgs.HookStatus,
	}
	for _, flag := range hookCommands {
		if flag {
			if err := commands.HandleHooks(cliArgs); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		}
	}

	configCommands := []bool{
		cliArgs.Config,
		cliArgs.CreateConfig,
	}
	for _, flag := range configCommands {
		if flag {
			if err := commands.HandleConfig(cliArgs); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		}
	}

	proc := processor.New(cliArgs)
	if err := proc.Process(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
