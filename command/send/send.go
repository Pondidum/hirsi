package send

import (
	"context"
	"fmt"
	"hirsi/tracing"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"

	_ "github.com/mattn/go-sqlite3"
)

var tr = otel.Tracer("send")

type SendCommand struct {
	addr string
}

func NewSendCommand() *SendCommand {
	return &SendCommand{
		addr: "http://localhost:5757",
	}
}

func (c *SendCommand) Synopsis() string {
	return "sends a message"
}

func (c *SendCommand) Flags() *pflag.FlagSet {

	flags := pflag.NewFlagSet("send", pflag.ContinueOnError)
	return flags
}

func (c *SendCommand) Execute(ctx context.Context, args []string) error {
	ctx, span := tr.Start(ctx, "execute")
	defer span.End()

	values := url.Values{}
	values.Add("message", strings.Join(args, " "))

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.addr+"/api/messages",
		strings.NewReader(values.Encode()))
	if err != nil {
		return tracing.Error(span, err)
	}

	client := http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
				return otelhttptrace.NewClientTrace(ctx)
			}),
		),
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 && res.StatusCode > 299 {
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("error from server:\n " + string(body))
	}

	return nil
}
