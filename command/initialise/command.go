package initialise

import (
	"context"
	"hirsi/config"
	"hirsi/storage"
	"hirsi/tracing"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"

	_ "github.com/mattn/go-sqlite3"
)

var tr = otel.Tracer("init")

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

func (c *InitCommand) Execute(ctx context.Context, args []string) error {
	ctx, span := tr.Start(ctx, "execute")
	defer span.End()

	cfg, err := config.CreateConfig(ctx)
	if err != nil {
		return tracing.Error(span, err)
	}

	if err := storage.MigrateDatabase(ctx, cfg.DbPath); err != nil {
		return err
	}

	return nil
}
