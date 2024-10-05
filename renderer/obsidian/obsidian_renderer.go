package obsidian

import (
	"bytes"
	"context"
	"fmt"
	"hirsi/message"
	"hirsi/tracing"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tr = otel.Tracer("renderer.obsidian")

type ObsidianRenderer struct {
	path string //./obsidian-dev/hirsi-dev

	// terms map[string]*regexp.Regexp

	terms []*Term
}

func NewObsidianRenderer(dirpath string) (*ObsidianRenderer, error) {
	renderer := &ObsidianRenderer{
		terms: []*Term{},
		path:  dirpath,
	}
	return renderer, nil
}

func (r *ObsidianRenderer) AddTitles(titles []string) error {

	for _, title := range titles {
		term, err := NewTerm(title, "")
		if err != nil {
			return err
		}

		r.terms = append(r.terms, term)
	}

	return nil
}

func (r *ObsidianRenderer) PopulateAutoLinker(ctx context.Context) error {
	ctx, span := tr.Start(ctx, "populate_autolinker")
	defer span.End()

	logPath := path.Join(r.path, "log")

	return filepath.WalkDir(r.path, func(p string, d fs.DirEntry, e error) error {
		ctx, span := tr.Start(ctx, "walk")
		defer span.End()

		span.SetAttributes(
			attribute.String("path", p),
			attribute.Bool("is_dir", d.IsDir()),
		)

		if strings.HasPrefix(d.Name(), ".") {
			span.SetAttributes(attribute.Bool("skip_dir", true))
			return fs.SkipDir
		}

		if strings.HasPrefix(p, logPath) {
			span.SetAttributes(attribute.Bool("skip_dir", true))
			return fs.SkipDir
		}

		if d.IsDir() {
			return nil
		}

		name := strings.TrimSuffix(d.Name(), path.Ext(d.Name()))

		span.SetAttributes(attribute.String("term", name))

		term, err := NewTerm(name, p)
		if err != nil {
			return tracing.ErrorCtx(ctx, err)
		}
		r.terms = append(r.terms, term)

		content, err := os.ReadFile(p)
		if err != nil {
			return tracing.ErrorCtx(ctx, err)
		}

		aliases, err := extractAliases(content)
		if err != nil {
			return tracing.ErrorCtx(ctx, err)
		}

		span.SetAttributes(attribute.StringSlice("aliases", aliases))

		for _, alias := range aliases {

			term, err := NewTerm(alias, p)
			if err != nil {
				return tracing.ErrorCtx(ctx, err)
			}
			r.terms = append(r.terms, term)
		}

		return nil
	})
}

func (r *ObsidianRenderer) Render(message *message.Message) error {
	logPath := path.Join(r.path, "log")

	if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
		return err
	}

	filename := message.WrittenAt.Format("2006-01-02")
	filepath := path.Join(logPath, filename+".md")

	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	content := formatMessage(r.terms, message)
	if _, err := f.Write(content); err != nil {
		return err
	}

	return nil
}

func formatMessage(terms []*Term, m *message.Message) []byte {

	sb := strings.Builder{}
	sb.WriteString("- ")
	sb.WriteString(m.WrittenAt.Format("15:04"))
	sb.WriteString(" ")
	sb.WriteString(buildTags(m.Tags))
	sb.WriteString("\n")

	lines := strings.Split(linkify(terms, m.Message), "\n")
	for _, line := range lines {
		sb.WriteString("\t")
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	return []byte(sb.String())
}

func buildTags(tags map[string]string) string {
	sb := strings.Builder{}

	for key, value := range tags {
		sb.WriteString(fmt.Sprintf("#%s", key))

		if value != "" {
			sb.WriteString(fmt.Sprintf("/%s", strings.TrimPrefix(value, "/")))
		}

		sb.WriteString(" ")
	}

	return strings.TrimSpace(sb.String())
}

func extractAliases(content []byte) ([]string, error) {

	m := &matter{}
	if _, err := frontmatter.Parse(bytes.NewReader(content), m); err != nil {
		return nil, err
	}

	return m.Aliases, nil
}

type matter struct {
	Aliases []string
}
