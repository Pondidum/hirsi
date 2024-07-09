package render

import (
	"context"
	"hirsi/config"
	"hirsi/message"
	"hirsi/renderer"
	"hirsi/storage"
	"hirsi/tracing"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"

	_ "github.com/mattn/go-sqlite3"
)

var tr = otel.Tracer("command.render")

type RenderCommand struct {
	cliOnly bool
}

func NewRenderCommand() *RenderCommand {
	return &RenderCommand{}
}

func (c *RenderCommand) Synopsis() string {
	return "render all logs"
}

func (c *RenderCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("render", pflag.ContinueOnError)

	flags.BoolVar(&c.cliOnly, "stdout", false, "render to stdout only")
	return flags
}

func (c *RenderCommand) Execute(ctx context.Context, args []string) error {
	ctx, span := tr.Start(ctx, "execute")
	defer span.End()

	cfg, err := config.CreateConfig(ctx)
	if err != nil {
		return tracing.Error(span, err)
	}

	var sink renderer.Renderer

	if c.cliOnly {

		r, err := renderer.NewCliRenderer()
		if err != nil {
			return err
		}
		sink = r
	} else {
		cfg, err := config.CreateConfig(ctx)
		if err != nil {
			return err
		}
		sink = renderer.NewCompositeRendererM(cfg.Renderers)
	}

	err = storage.EachMessage(ctx, cfg.DbPath, 0, func(m *message.Message) error {
		return sink.Render(m)
	})
	if err != nil {
		return err
	}

	return nil
}
