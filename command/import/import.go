package importcmd

import (
	"context"
	"fmt"
	"hirsi/config"
	"hirsi/enhancement"
	"hirsi/message"
	"hirsi/storage"
	"hirsi/tracing"
	"io/fs"
	"os"
	"path"
	"path/filepath"
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

var tr = otel.Tracer("command.import")

type ImportCommand struct {
}

func NewImportCommand() *ImportCommand {
	return &ImportCommand{}
}

func (c *ImportCommand) Synopsis() string {
	return "import existing logs"
}

func (c *ImportCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("import", pflag.ContinueOnError)
	return flags
}

func (c *ImportCommand) Execute(ctx context.Context, cfg *config.Config, args []string) error {
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
		if err := c.importDir(ctx, cfg, args[0]); err != nil {
			return err
		}
	} else {
		if err := c.importFile(ctx, cfg, args[0]); err != nil {
			return err
		}
	}

	return nil
}

func (c *ImportCommand) importDir(ctx context.Context, cfg *config.Config, dirpath string) error {
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

		if err := c.parseFile(ctx, cfg, dirpath, filename, content); err != nil {
			return err
		}
	}

	return nil
}

func (c *ImportCommand) importFile(ctx context.Context, cfg *config.Config, filepath string) error {

	content, err := os.ReadFile(filepath)
	if err != nil {
		return tracing.ErrorCtx(ctx, err)
	}

	dirpath := path.Dir(filepath)
	filename := path.Base(filepath)

	if err := c.parseFile(ctx, cfg, dirpath, filename, content); err != nil {
		return err
	}

	return nil
}

var rx = regexp.MustCompile(`^- (\d\d):(\d\d) .*`)

func (c *ImportCommand) parseFile(ctx context.Context, cfg *config.Config, dirPath string, filename string, content []byte) error {
	ctx, span := tr.Start(ctx, "parse_file")
	defer span.End()

	dirPath, err := filepath.Abs(dirPath)
	if err != nil {
		return tracing.Error(span, err)
	}

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

		msg := &message.Message{
			WrittenAt: entryTime,
			Message:   entry[8:],
			Tags:      &message.Tags{},
		}
		msg.Tags.Add("imported", "")

		for _, e := range cfg.Enhancements {
			if err := e.Enhance(msg); err != nil {
				return err
			}
		}

		msg.Tags.Set(enhancement.PwdTag, dirPath)

		if err := storage.StoreMessage(ctx, cfg.DbPath, msg); err != nil {
			return err
		}

		for _, r := range cfg.Pipelines {
			if err := r.Handle(msg); err != nil {
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
