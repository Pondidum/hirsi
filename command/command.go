package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"
)

type CommandDefinition interface {
	Synopsis() string
	Flags() *pflag.FlagSet
	Execute(args []string) error
}

func NewCommand(definition CommandDefinition) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return &command{definition}, nil
	}
}

type command struct {
	CommandDefinition
}

func (c *command) Help() string {
	sb := strings.Builder{}

	sb.WriteString(c.Synopsis())
	sb.WriteString("\n\n")

	sb.WriteString("Flags:\n\n")

	sb.WriteString(c.Flags().FlagUsagesWrapped(80))

	return sb.String()
}

func (c *command) Run(args []string) int {

	flags := c.Flags()

	if err := flags.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	if err := c.Execute(flags.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	return 0
}
