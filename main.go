package main

import (
	"fmt"
	"hirsi/command"
	"hirsi/command/send"
	"hirsi/command/server"
	"hirsi/command/watch"
	"os"

	"github.com/mitchellh/cli"
)

func main() {

	commands := map[string]cli.CommandFactory{
		"send":   command.NewCommand(send.NewSendCommand()),
		"server": command.NewCommand(server.NewServerCommand()),
		"watch":  command.NewCommand(watch.NewWatchCommand()),
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
