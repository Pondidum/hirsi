package renderer

import (
	"bytes"
	"fmt"
	"hirsi/message"
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
		"Tags":      r.buildTags(m.Tags),
	})

	fmt.Println(buf.String())

	return nil
}

func (r *CliRenderer) buildTags(t map[string]string) string {
	tags := make([]string, 0, len(t))

	for key, value := range t {
		if value != "" {
			tags = append(tags, fmt.Sprintf("#%s: %s", key, value))
		} else {
			tags = append(tags, fmt.Sprintf("#%s", key))
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
