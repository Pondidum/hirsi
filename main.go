package main

import (
	"fmt"
	"hirsi/command"
	"hirsi/command/basic"
	"hirsi/command/initialise"
	"hirsi/command/ls"
	"hirsi/config"
	"os"

	"github.com/mitchellh/cli"
)

func main() {

	commands := map[string]cli.CommandFactory{
		"write": command.NewCommand(basic.NewBasicCommand(config.AppConfig)),
		"init":  command.NewCommand(initialise.NewInitCommand(config.AppConfig)),
		"ls":    command.NewCommand(ls.NewLsCommand(config.AppConfig)),
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
