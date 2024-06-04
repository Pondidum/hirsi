package importcmd

import (
	"context"
	"fmt"
	"hirsi/config"
	"hirsi/message"
	"hirsi/storage"
	"hirsi/tracing"
	"io/fs"
	"os"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/md"
	"github.com/gomarkdown/markdown/parser"

	_ "github.com/mattn/go-sqlite3"
)

var tr = otel.Tracer("import")

type ImportCommand struct {
	config *config.Config
}

func NewImportCommand(config *config.Config) *ImportCommand {
	return &ImportCommand{config}
}

func (c *ImportCommand) Synopsis() string {
	return "import existing logs"
}

func (c *ImportCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("import", pflag.ContinueOnError)
	return flags
}

func (c *ImportCommand) Execute(ctx context.Context, args []string) error {
	ctx, span := tr.Start(ctx, "execute")
	defer span.End()

	if len(args) != 1 {
		return fmt.Errorf("this command takes exactly 1 argument: dir_path or file_path")
	}

	stats, err := os.Stat(args[0])
	if err != nil {
		return err
	}

	if stats.IsDir() {
		if err := c.importDir(ctx, args[0]); err != nil {
			return err
		}
	} else {
		if err := c.importFile(ctx, args[0]); err != nil {
			return err
		}
	}

	return nil
}

func (c *ImportCommand) importDir(ctx context.Context, dirpath string) error {
	source := os.DirFS(dirpath)

	entries, err := fs.ReadDir(source, ".")
	if err != nil {
		return tracing.ErrorCtx(ctx, err)
	}

	files := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.Type().IsRegular() {
			files = append(files, e.Name())
		}
	}

	slices.Sort(files)

	fmt.Println("==> Importing", len(files), "log files")

	for _, filename := range files {

		content, err := fs.ReadFile(source, filename)
		if err != nil {
			return tracing.ErrorCtx(ctx, err)
		}

		if err := c.parseFile(ctx, filename, content); err != nil {
			return err
		}
	}

	return nil
}

func (c *ImportCommand) importFile(ctx context.Context, filepath string) error {

	content, err := os.ReadFile(filepath)
	if err != nil {
		return tracing.ErrorCtx(ctx, err)
	}

	filename := path.Base(filepath)

	if err := c.parseFile(ctx, filename, content); err != nil {
		return err
	}

	return nil
}

var rx = regexp.MustCompile(`^- (\d\d):(\d\d) .*`)

func (c *ImportCommand) parseFile(ctx context.Context, filename string, content []byte) error {

	fmt.Println("==>", filename)
	date, err := time.Parse("2006-01-02", strings.TrimSuffix(filename, path.Ext(filename)))
	if err != nil {
		return tracing.ErrorCtx(ctx, err)
	}

	seconds := 0
	prevHours := 0
	prevMins := 0

	for _, entry := range extractEntries(ctx, content) {
		groups := rx.FindStringSubmatch(entry)

		if len(groups) < 3 {
			continue
		}

		hours, err := strconv.Atoi(groups[1])
		if err != nil {
			continue
		}
		minutes, err := strconv.Atoi(groups[2])
		if err != nil {
			continue
		}

		entryTime := date.Add(time.Duration(hours) * time.Hour).Add(time.Duration(minutes) * time.Minute)

		if hours == prevHours && minutes == prevMins {
			seconds++
			entryTime = entryTime.Add(time.Duration(seconds) * time.Second)
		}

		message := &message.Message{
			WrittenAt: entryTime,
			Message:   entry[8:],
			Tags: map[string]string{
				"import.source": filename,
			},
		}

		for _, e := range c.config.Enhancements {
			if err := e.Enhance(message); err != nil {
				return err
			}
		}

		if err := storage.StoreMessage(ctx, c.config.DbPath, message); err != nil {
			return err
		}

		for _, r := range c.config.Renderers {
			if err := r.Render(message); err != nil {
				return err
			}
		}

		prevHours = hours
		prevMins = minutes
	}

	return nil
}

func extractEntries(ctx context.Context, content []byte) []string {
	ctx, span := tr.Start(ctx, "extract_entries")
	defer span.End()

	p := parser.New()
	doc := p.Parse(content)

	span.SetAttributes(attribute.Bool("content.parsed", true))

	list := doc.GetChildren()[0]
	items := list.GetChildren()

	entries := make([]string, len(items))

	for i, item := range items {
		r := md.NewRenderer()

		output := markdown.Render(item, r)
		entries[i] = string(output)
	}

	return entries
}
