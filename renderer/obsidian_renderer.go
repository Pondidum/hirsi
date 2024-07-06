package renderer

import (
	"bytes"
	"fmt"
	"hirsi/message"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

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
		if _, found := r.terms[title]; found {
			continue
		}

		var err error
		fmt.Println("add title", title)
		if r.terms[title], err = regexp.Compile(`(?i)^(\W*?)(` + title + `)(\W*?)$`); err != nil {
			return err
		}
	}

	return nil
}

func (r *ObsidianRenderer) PopulateAutoLinker() error {
	logPath := path.Join(r.path, "log")

	return filepath.WalkDir(r.path, func(p string, d fs.DirEntry, e error) error {
		if d.IsDir() {
			return nil
		}

		if strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		if strings.HasPrefix(p, logPath) {
			return nil
		}

		term := strings.TrimSuffix(d.Name(), path.Ext(d.Name()))
		if _, found := r.terms[term]; !found {

			fmt.Println("add term", term)
			rx, err := regexp.Compile(`(?i)^(\W*?)(` + term + `)(\W*?)$`)
			if err != nil {
				return err
			}
			r.terms[term] = rx

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
		"Message":   r.linkify(m.Message),
		"WrittenAt": m.WrittenAt,
		"Tags":      r.buildTags(m.Tags),
	})

	return buf.Bytes()
}

func (r *ObsidianRenderer) linkify(message string) string {

	words := strings.Split(message, " ")
	for i, word := range words {
		for term, rx := range r.terms {

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

func (r *ObsidianRenderer) buildTags(tags map[string]string) string {
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
