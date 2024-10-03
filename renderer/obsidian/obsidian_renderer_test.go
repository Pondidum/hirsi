package obsidian

import (
	"hirsi/enhancement"
	"hirsi/message"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestObsidianLinkin(t *testing.T) {

	input := `- 10:45 TEAMCITY: teamcity build for proxy - [https://teamcity.example.com/project.html?projectId=someproject&tab=projectOverview](https://teamcity.example.com/project.html?projectId=someproject&tab=projectOverview)`

	terms := map[string]*regexp.Regexp{}
	addTerm(terms, "teamcity")

	actual := linkify(terms, input)

	expected := "- 10:45 [[TEAMCITY]]: [[teamcity]] build for proxy - [https://teamcity.example.com/project.html?projectId=someproject&tab=projectOverview](https://teamcity.example.com/project.html?projectId=someproject&tab=projectOverview)"
	assert.Equal(t, expected, actual)
}

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

	terms := map[string]*regexp.Regexp{}
	addTerm(terms, "teamcity")

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
