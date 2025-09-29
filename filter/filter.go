package filter

import (
	"hirsi/message"
	"strings"
)

type Filter func(m *message.Message) bool

var _ Filter = HasTagWithPrefix("", "")
var _ Filter = HasTagWithoutPrefix("", "")

func HasTagWithPrefix(tag string, prefix string) func(m *message.Message) bool {
	return func(m *message.Message) bool {
		for tag := range m.Tags() {
			if tag.Key == "pwd" && strings.HasPrefix(tag.Value, prefix) {
				return true
			}
		}
		return false
	}
}

func HasTagWithoutPrefix(tag string, prefix string) func(m *message.Message) bool {
	return func(m *message.Message) bool {
		for tag := range m.Tags() {
			if tag.Key == "pwd" && !strings.HasPrefix(tag.Value, prefix) {
				return true
			}
		}
		return false
	}
}
