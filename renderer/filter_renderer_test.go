package renderer

import (
	"hirsi/enhancement"
	"hirsi/message"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterRenderer(t *testing.T) {
	cases := []struct {
		name     string
		filter   func(m *message.Message) bool
		expected []string
	}{
		{
			name:     "only prefixed",
			filter:   HasTagWithPrefix("pwd", "/home/andy/dev/secret"),
			expected: []string{"with tag"},
		},
		{
			name:     "only non-prefix",
			filter:   HasTagWithoutPrefix("pwd", "/home/andy/dev/secret"),
			expected: []string{"without tag"},
		},
		{
			name:     "no filters",
			expected: []string{"with tag", "without tag"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			renderer := &MemoryRenderer{}
			f := NewFilteredRenderer(tc.filter, renderer)

			f.Render(&message.Message{
				Message: "with tag",
				Tags: map[string]string{
					enhancement.PwdTag: "/home/andy/dev/secret/project",
				},
			})

			f.Render(&message.Message{
				Message: "without tag",
				Tags: map[string]string{
					enhancement.PwdTag: "/home/andy/dev/projects/hirsi",
				},
			})

			rendered := []string{}
			for _, m := range renderer.messages {
				rendered = append(rendered, m.Message)
			}

			assert.Equal(t, tc.expected, rendered)
		})
	}

	// assert.Equal(t, "without tag", renderer.messages[0].Message)
	// assert.Len(t, renderer.messages, 1)
}
