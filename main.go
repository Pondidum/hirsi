package main

import (
	"fmt"
	"hirsi/command"
	"hirsi/command/send"
	"os"

	"github.com/mitchellh/cli"
)

func main() {

	commands := map[string]cli.CommandFactory{
		"send": command.NewCommand(&send.SendCommand{}),
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
