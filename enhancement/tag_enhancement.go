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
	field     string
	condition string
	values    map[string][]string
}

func (e *TagAddEnhancement) Enhance(m *message.Message) error {
	val, found := m.Tags[e.field]
	if !found {
		return nil
	}

	switch e.condition {
	case "equals":
		if kvp, found := e.values[val]; found {
			m.Tags[kvp[0]] = kvp[1]
		}

	case "prefix":
		for search, kvp := range e.values {
			if strings.HasPrefix(val, search) {
				m.Tags[kvp[0]] = kvp[1]
			}
		}

	}

	return nil
}
