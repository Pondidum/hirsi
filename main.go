package main

import (
	"fmt"
	"hirsi/command"
	"hirsi/command/basic"
	importcmd "hirsi/command/import"
	"hirsi/command/initialise"
	"hirsi/command/ls"
	"hirsi/config"
	"os"

	"github.com/mitchellh/cli"
)

func main() {

	cfg, err := config.CreateConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		os.Exit(1)
	}

	commands := map[string]cli.CommandFactory{
		"write":  command.NewCommand(basic.NewBasicCommand(cfg)),
		"init":   command.NewCommand(initialise.NewInitCommand(cfg)),
		"ls":     command.NewCommand(ls.NewLsCommand(cfg)),
		"import": command.NewCommand(importcmd.NewImportCommand(cfg)),
	}

	cli := &cli.CLI{
		Name:                       "hirsi",
		Args:                       os.Args[1:],
		Commands:                   commands,
		Autocomplete:               true,
		AutocompleteNoDefaultFlags: false,
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
	}

	os.Exit(exitCode)
}
