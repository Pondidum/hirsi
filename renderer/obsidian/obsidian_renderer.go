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
	"regexp"
	"strings"
	"text/template"

	"github.com/adrg/frontmatter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tr = otel.Tracer("renderer.obsidian")

type ObsidianRenderer struct {
	path     string //./obsidian-dev/hirsi-dev
	template *template.Template

	terms map[string]*regexp.Regexp
}

func NewObsidianRenderer(dirpath string) (*ObsidianRenderer, error) {
	tpl, err := template.New("").Parse(`- {{ .WrittenAt.Format "15:04" }} {{ .Tags }}
	{{ .Message }}
`)
	if err != nil {
		return nil, err
	}

	renderer := &ObsidianRenderer{
		terms:    map[string]*regexp.Regexp{},
		path:     dirpath,
		template: tpl,
	}
	return renderer, nil
}

func (r *ObsidianRenderer) AddTitles(titles []string) error {

	for _, title := range titles {
		if err := r.addTerm(title); err != nil {
			return err
		}
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

		term := strings.TrimSuffix(d.Name(), path.Ext(d.Name()))

		span.SetAttributes(attribute.String("term", term))

		if err := r.addTerm(term); err != nil {
			return tracing.ErrorCtx(ctx, err)
		}

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
			if err := r.addTerm(alias); err != nil {
				return tracing.ErrorCtx(ctx, err)
			}
		}

		return nil
	})
}

func (r *ObsidianRenderer) addTerm(term string) error {
	if _, found := r.terms[term]; !found {

		rx, err := regexp.Compile(`(?i)^(\W*?)(` + term + `)(\W*?)$`)
		if err != nil {
			return err
		}
		r.terms[term] = rx

	}
	return nil
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

	content := r.formatMessage(message)
	if _, err := f.Write(content); err != nil {
		return err
	}

	return nil
}

func (r *ObsidianRenderer) formatMessage(m *message.Message) []byte {
	// nested tags for each entry, i.e. `path=/home/andy/dev` would be `#path/home/andy/dev`

	buf := bytes.Buffer{}
	r.template.Execute(&buf, map[string]any{
		"Message":   linkify(r.terms, m.Message),
		"WrittenAt": m.WrittenAt,
		"Tags":      buildTags(m.Tags),
	})

	return buf.Bytes()
}

func linkify(terms map[string]*regexp.Regexp, message string) string {

	words := strings.Split(message, " ")
	for i, word := range words {
		for term, rx := range terms {

			if word == term {
				words[i] = "[[" + word + "]]"
				continue
			} else if rx.MatchString(word) {
				words[i] = rx.ReplaceAllString(word, "$1[[$2]]$3")
				continue
			}
		}
	}

	return strings.Join(words, " ")
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
