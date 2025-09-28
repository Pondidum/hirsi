package renderer

import (
	"bytes"
	"fmt"
	"hirsi/message"
	"iter"
	"strings"
	"text/template"
)

const format = `- {{ .WrittenAt.Format "15:04" }} {{ .Tags }}
{{ .Message | indent 2 }}`

func NewCliRenderer() (*CliRenderer, error) {

	funcs := template.FuncMap{
		"indent": indent,
	}
	tpl, err := template.New("").Funcs(funcs).Parse(format)
	if err != nil {
		return nil, err
	}

	return &CliRenderer{tpl}, nil
}

type CliRenderer struct {
	template *template.Template
}

func (r *CliRenderer) Render(m *message.Message) error {

	buf := bytes.Buffer{}
	r.template.Execute(&buf, map[string]any{
		"Message":   m.Message,
		"WrittenAt": m.WrittenAt,
		"Tags":      r.buildTags(m.Tags.All()),
	})

	fmt.Println(buf.String())

	return nil
}

func (r *CliRenderer) buildTags(t iter.Seq[message.Tag]) string {
	tags := []string{}

	for tag := range t {
		if tag.Value != "" {
			tags = append(tags, fmt.Sprintf("#%s: %s", tag.Key, tag.Value))
		} else {
			tags = append(tags, fmt.Sprintf("#%s", tag.Key))
		}
	}

	return strings.Join(tags, " ")
}

func indent(count int, s string) string {
	prefix := strings.Repeat(" ", count)
	lines := strings.Split(s, "\n")

	for i, line := range lines {
		lines[i] = prefix + line
	}

	return strings.Join(lines, "\n")
}
