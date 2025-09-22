package log

import (
	"bytes"
	"fmt"
	"hirsi/message"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"
)

func NewLogRenderer(path string) (*LogRenderer, error) {
	tpl, err := template.New("").Parse(`- {{ .WrittenAt.Format "15:04" }} {{ .Tags }}: {{ .Message }}`)
	if err != nil {
		return nil, err
	}

	return &LogRenderer{path, tpl}, nil
}

type LogRenderer struct {
	Path     string
	template *template.Template
}

func (r *LogRenderer) Render(message *message.Message) error {
	if err := os.MkdirAll(r.Path, os.ModePerm); err != nil {
		return err
	}

	filename := message.WrittenAt.Format("2006-01-02")
	filepath := path.Join(r.Path, filename+".md")

	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := r.writeMessage(f, message); err != nil {
		return err
	}

	return nil
}

func (r *LogRenderer) writeMessage(w io.Writer, m *message.Message) error {
	content := r.formatMessage(m)

	if _, err := w.Write(append(content, '\n')); err != nil {
		return err
	}

	return nil
}

func (r *LogRenderer) formatMessage(m *message.Message) []byte {
	// nested tags for each entry, i.e. `path=/home/andy/dev` would be `#path/home/andy/dev`

	buf := bytes.Buffer{}
	r.template.Execute(&buf, map[string]any{
		"Message":   m.Message,
		"WrittenAt": m.WrittenAt,
		"Tags":      r.buildTags(m.Tags),
	})

	return buf.Bytes()
}

func (r *LogRenderer) buildTags(t []message.Tag) string {
	tags := make([]string, 0, len(t))

	for _, tag := range t {
		if tag.Value != "" {
			tags = append(tags, fmt.Sprintf("#%s: %s", tag.Key, tag.Value))
		} else {
			tags = append(tags, fmt.Sprintf("#%s", tag.Key))
		}
	}

	sort.Strings(tags)

	return strings.Join(tags, " ")
}
