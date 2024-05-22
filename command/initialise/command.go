package initialise

import (
	"context"
	"hirsi/config"
	"hirsi/storage"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"

	_ "github.com/mattn/go-sqlite3"
)

var tr = otel.Tracer("init")

type InitCommand struct {
	config *config.Config
}

func NewInitCommand(config *config.Config) *InitCommand {
	return &InitCommand{config}
}

func (c *InitCommand) Synopsis() string {
	return "init hirsi"
}

func (c *InitCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("init", pflag.ContinueOnError)
	return flags
}

func (c *InitCommand) Execute(ctx context.Context, args []string) error {
	ctx, span := tr.Start(ctx, "execute")
	defer span.End()

	if err := storage.MigrateDatabase(ctx, c.config.DbPath); err != nil {
		return err
	}

	return nil
}
