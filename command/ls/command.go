package ls

import (
	"context"
	"fmt"
	"hirsi/config"
	"hirsi/storage"
	"hirsi/tracing"
	"strings"

	"github.com/ryanuber/columnize"
	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"

	_ "github.com/mattn/go-sqlite3"
)

var tr = otel.Tracer("command.ls")

type LsCommand struct {
}

func NewLsCommand() *LsCommand {
	return &LsCommand{}
}

func (c *LsCommand) Synopsis() string {
	return "view the logs"
}

func (c *LsCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("ls", pflag.ContinueOnError)
	return flags
}

func (c *LsCommand) Execute(ctx context.Context, args []string) error {
	ctx, span := tr.Start(ctx, "execute")
	defer span.End()

	cfg, err := config.CreateConfig(ctx)
	if err != nil {
		return tracing.Error(span, err)
	}

	messages, err := storage.ListMessages(ctx, cfg.DbPath, 10)
	if err != nil {
		return err
	}

	output := make([]string, len(messages)+1)
	output[0] = "stored_at | written_at | message | tags"

	for i, m := range messages {
		output[i+1] = fmt.Sprintf("%s | %s | %s | %s", m.StoredAt, m.WrittenAt, m.Message, tagsCsv(m.Tags))
	}

	fmt.Println(tableOutput(output))

	return nil
}
func tagsCsv(tags map[string]string) string {

	sb := strings.Builder{}
	for k, v := range tags {
		sb.WriteString(fmt.Sprintf("%s=%s,", k, v))
	}

	return strings.TrimSuffix(sb.String(), ",")
}
func tableOutput(list []string) string {
	if len(list) == 0 {
		return ""
	}

	delim := "|"
	underline := ""
	headers := strings.Split(list[0], delim)
	for i, h := range headers {
		h = strings.TrimSpace(h)
		u := strings.Repeat("-", len(h))

		underline = underline + u
		if i != len(headers)-1 {
			underline = underline + delim
		}
	}

	list = append(list, "")
	copy(list[2:], list[1:])
	list[1] = underline

	return columnize.Format(list, &columnize.Config{
		Glue: "    ",
	})
}
