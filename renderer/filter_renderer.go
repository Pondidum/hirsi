package renderer

import (
	"hirsi/message"
	"strings"
)

func HasTagWithPrefix(tag string, prefix string) func(m *message.Message) bool {
	return func(m *message.Message) bool {
		for tag, val := range m.Tags {
			if tag == "pwd" && strings.HasPrefix(val, "/home/andy/dev/secret") {
				return true
			}
		}
		return false
	}
}

func HasTagWithoutPrefix(tag string, prefix string) func(m *message.Message) bool {
	return func(m *message.Message) bool {
		for tag, val := range m.Tags {
			if tag == "pwd" && !strings.HasPrefix(val, "/home/andy/dev/secret") {
				return true
			}
		}
		return false
	}
}

func NewFilteredRenderer(filter func(m *message.Message) bool, other Renderer) *FilterRenderer {
	return &FilterRenderer{
		other:  other,
		filter: filter,
	}
}

type FilterRenderer struct {
	other  Renderer
	filter func(m *message.Message) bool
}

func (r *FilterRenderer) Render(message *message.Message) error {

	if r.filter == nil || r.filter(message) {
		return r.other.Render(message)
	}

	return nil
}

// func (r *FilterRenderer) isExcluded(message *message.Message) bool {
// 	if len(r.filter.excludeTagFilters) == 0 {
// 		return false
// 	}

// 	for tag, filters := range r.filter.excludeTagFilters {
// 		val, found :=message.Tags[tag]
// 		if !found {
// 			continue
// 		}

// 	}

// 	return false
// }

// func hasTagMatch(message *message.Message, tag string, filters []*regexp.Regexp) bool {
// 	val, found := message.Tags[tag]
// 	if !found {
// 		return false
// 	}

// 	for _, filter := range filters {
// 		if filter.MatchString(val) {
// 			return true
// 		}
// 	}

// 	return false
// }

// type FilterOption func(cfg *filterConfig) error

// type filterConfig struct {
// 	includeTagFilters map[string][]*regexp.Regexp
// 	excludeTagFilters map[string][]*regexp.Regexp
// }

// func WithTagPrefix(tag string, prefix string) FilterOption {
// 	return func(cfg *filterConfig) error {
// 		rx, err := regexp.Compile(`^` + prefix + `.*`)
// 		if err != nil {
// 			return err
// 		}
// 		cfg.includeTagFilters[tag] = append(cfg.includeTagFilters[tag], rx)
// 		return nil
// 	}
// }
// func WithoutTagPrefix(tag string, prefix string) FilterOption {
// 	return func(cfg *filterConfig) error {
// 		rx, err := regexp.Compile(`^` + prefix + `.*`)
// 		if err != nil {
// 			return err
// 		}
// 		cfg.excludeTagFilters[tag] = append(cfg.excludeTagFilters[tag], rx)
// 		return nil
// 	}
// }
