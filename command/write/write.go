package write

import (
	"context"
	"hirsi/config"
	"hirsi/message"
	"hirsi/storage"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"

	_ "github.com/mattn/go-sqlite3"
)

var tr = otel.Tracer("command.writer")

type WriteCommand struct {
}

func NewWriteCommand() *WriteCommand {
	return &WriteCommand{}
}

func (c *WriteCommand) Synopsis() string {
	return "storage a message"
}

func (c *WriteCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("basic", pflag.ContinueOnError)
	return flags
}

func (c *WriteCommand) Execute(ctx context.Context, cfg *config.Config, args []string) error {
	ctx, span := tr.Start(ctx, "execute")
	defer span.End()

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
