package obsidian

import (
	"hirsi/enhancement"
	"hirsi/message"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExtractFrontMatter(t *testing.T) {

	content, err := os.ReadFile("../../obsidian-dev/hirsi-dev/devportal.md")
	assert.NoError(t, err)

	aliases, err := extractAliases(content)
	assert.NoError(t, err)
	assert.Equal(t, []string{"devportal-src"}, aliases)
}

func TestFormatMessage(t *testing.T) {

	cases := []struct {
		message  string
		tags     map[string]string
		expected string
	}{
		{
			message: "This is a single line with one tag",
			tags:    map[string]string{enhancement.PwdTag: "/tmp/test/dir"},
			expected: `- 22:26 #pwd/tmp/test/dir
	This is a single line with one tag
`,
		},
		{
			message: "This is multiple\nlines of message.\noh,and one tag",
			tags:    map[string]string{enhancement.PwdTag: "/tmp/test/dir"},
			expected: `- 22:26 #pwd/tmp/test/dir
	This is multiple
	lines of message.
	oh,and one tag
`,
		},
	}

	teamcity, _ := NewTerm("teamcity", "")
	terms := []*Term{
		teamcity,
	}

	for _, tc := range cases {
		t.Run(tc.message, func(t *testing.T) {

			actual := formatMessage(terms, &message.Message{
				WrittenAt: time.Date(2024, 10, 03, 22, 26, 31, 0, time.UTC),
				Message:   tc.message,
				Tags:      tc.tags,
			})

			assert.Equal(t, tc.expected, string(actual))
		})
	}

}
