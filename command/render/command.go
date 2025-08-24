package render

import (
	"context"
	"hirsi/config"
	"hirsi/message"
	"hirsi/pipeline"
	"hirsi/renderer"
	"hirsi/storage"
	"hirsi/tracing"
	"slices"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"

	_ "github.com/mattn/go-sqlite3"
)

var tr = otel.Tracer("command.render")

type RenderCommand struct {
	cliOnly        bool
	namedPipelines []string
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
	flags.StringSliceVar(&c.namedPipelines, "pipelines", []string{}, "only render with the given pipelines")

	return flags
}

func (c *RenderCommand) Execute(ctx context.Context, cfg *config.Config, args []string) error {
	ctx, span := tr.Start(ctx, "execute")
	defer span.End()

	if c.cliOnly {

		r, err := renderer.NewCliRenderer()
		if err != nil {
			return tracing.Error(span, err)
		}
		return storage.EachMessage(ctx, cfg.DbPath, 0, func(m *message.Message) error {
			return r.Render(m)
		})

	}

	cfg, err := config.CreateConfig()
	if err != nil {
		return err
	}

	pipelines := cfg.Pipelines

	if len(c.namedPipelines) > 0 {

		pipelines = make([]*pipeline.Pipeline, 0, len(c.namedPipelines))
		for _, pipeline := range cfg.Pipelines {
			if slices.Contains(c.namedPipelines, pipeline.Name) {
				pipelines = append(pipelines, pipeline)
			}
		}
	}

	return storage.EachMessage(ctx, cfg.DbPath, 0, func(m *message.Message) error {
		for _, pipe := range pipelines {
			if err := pipe.Handle(m); err != nil {
				return err
			}
		}
		return nil
	})
}
