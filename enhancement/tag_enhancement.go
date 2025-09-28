package enhancement

import (
	"hirsi/message"
	"strings"
)

type TagAddConfig struct {
	Check     string
	Condition string
	Values    map[string][]string
}

func (ec *TagAddConfig) Build() Enhancement {
	return &TagAddEnhancement{ec.Check, ec.Condition, ec.Values}
}

type TagAddEnhancement struct {
	field        string
	condition    string
	replacements map[string][]string
}

func (e *TagAddEnhancement) Enhance(m *message.Message) error {

	for tag := range m.Tags.All() {
		if tag.Key == e.field {

			// note the matching logic can be pre built on construction, rather than per tag
			switch e.condition {
			case "equals":
				if replacement, found := e.replacements[tag.Value]; found {
					m.Tags.AddT(buildTag(replacement))
				}

			case "prefix":
				for search, replacement := range e.replacements {
					if strings.HasPrefix(tag.Value, search) {
						m.Tags.AddT(buildTag(replacement))
					}
				}

			}

		}
	}

	return nil
}

func buildTag(replacement []string) message.Tag {
	if len(replacement) == 1 {
		return message.Tag{Key: replacement[0]}
	}

	return message.Tag{Key: replacement[0], Value: replacement[1]}
}
