package basic

import (
	"context"
	"hirsi/config"
	"hirsi/message"
	"hirsi/storage"
	"hirsi/tracing"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"

	_ "github.com/mattn/go-sqlite3"
)

var tr = otel.Tracer("basic")

type BasicCommand struct {
}

func NewBasicCommand() *BasicCommand {
	return &BasicCommand{}
}

func (c *BasicCommand) Synopsis() string {
	return "storage a message"
}

func (c *BasicCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("basic", pflag.ContinueOnError)
	return flags
}

func (c *BasicCommand) Execute(ctx context.Context, args []string) error {
	ctx, span := tr.Start(ctx, "execute")
	defer span.End()

	cfg, err := config.CreateConfig(ctx)
	if err != nil {
		return tracing.Error(span, err)
	}

	message := &message.Message{
		WrittenAt: time.Now(),
		Message:   strings.Join(args, " "),
		Tags:      map[string]string{},
	}

	for _, e := range cfg.Enhancements {
		if err := e.Enhance(message); err != nil {
			return err
		}
	}

	if err := storage.StoreMessage(ctx, cfg.DbPath, message); err != nil {
		return err
	}

	for _, r := range cfg.Renderers {
		if err := r.Render(message); err != nil {
			return err
		}
	}

	return nil
}
