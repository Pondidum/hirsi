package main

import (
	"fmt"
	"hirsi/command"
	importcmd "hirsi/command/import"
	"hirsi/command/initialise"
	"hirsi/command/ls"
	"hirsi/command/write"
	"os"

	"github.com/mitchellh/cli"
)

func main() {

	commands := map[string]cli.CommandFactory{
		"write":  command.NewCommand(write.NewWriteCommand()),
		"init":   command.NewCommand(initialise.NewInitCommand()),
		"ls":     command.NewCommand(ls.NewLsCommand()),
		"import": command.NewCommand(importcmd.NewImportCommand()),
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
