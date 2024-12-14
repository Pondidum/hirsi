package write

import (
	"context"
	"fmt"
	"hirsi/config"
	"hirsi/message"
	"hirsi/storage"
	"os"
	"os/exec"
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
	return "write a message to the log"
}

func (c *WriteCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("basic", pflag.ContinueOnError)
	return flags
}

func (c *WriteCommand) Execute(ctx context.Context, cfg *config.Config, args []string) error {
	ctx, span := tr.Start(ctx, "execute")
	defer span.End()

	content := strings.Join(args, " ")

	if len(args) == 0 {
		note, err := c.launchEditor(ctx)
		if err != nil {
			return err
		}

		content = note
	}

	message := &message.Message{
		WrittenAt: time.Now(),
		Message:   content,
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

func (c *WriteCommand) launchEditor(ctx context.Context) (string, error) {
	tmp, err := os.CreateTemp("", "hirsi")
	if err != nil {
		return "", err
	}
	tmp.Close() // so that the editor writes it!

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	switch strings.ToLower(editor) {
	case "hx":
		cmd := exec.CommandContext(ctx, "hx", tmp.Name())
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", err
		}

	case "vim":
	default:
		cmd := exec.CommandContext(ctx, "vim", "+normal G$", "+startinsert", tmp.Name())
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", err
		}

	}

	content, err := os.ReadFile(tmp.Name())
	if err != nil {
		return "", err
	}

	if len(content) == 0 {
		return "", fmt.Errorf("file was empty, aborting")
	}

	return string(content), nil
}
