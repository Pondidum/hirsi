package initialise

import (
	"context"
	"hirsi/config"
	"hirsi/storage"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"

	_ "github.com/mattn/go-sqlite3"
)

var tr = otel.Tracer("command.init")

type InitCommand struct {
}

func NewInitCommand() *InitCommand {
	return &InitCommand{}
}

func (c *InitCommand) Synopsis() string {
	return "init hirsi"
}

func (c *InitCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("init", pflag.ContinueOnError)
	return flags
}

func (c *InitCommand) Execute(ctx context.Context, cfg *config.Config, args []string) error {
	ctx, span := tr.Start(ctx, "execute")
	defer span.End()

	if err := storage.MigrateDatabase(ctx, cfg.DbPath); err != nil {
		return err
	}

	return nil
}
