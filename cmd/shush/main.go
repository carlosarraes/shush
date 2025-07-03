package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/carlosarraes/shush/internal/processor"
	"github.com/carlosarraes/shush/internal/types"
)

var version = "0.1.0"

func main() {
	var cli types.CLI
	kong.Parse(&cli, 
		kong.Description("Remove comments from source code files"),
		kong.Vars{"version": version})

	if cli.Path == "" {
		fmt.Fprintf(os.Stderr, "Error: path argument is required\n")
		os.Exit(1)
	}

	if cli.Inline && cli.Block {
		fmt.Fprintf(os.Stderr, "Error: --inline and --block flags are mutually exclusive\n")
		os.Exit(1)
	}

	proc := processor.New(cli)
	if err := proc.Process(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}