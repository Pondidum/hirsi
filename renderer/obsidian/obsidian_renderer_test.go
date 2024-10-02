package obsidian

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObsidianLinkin(t *testing.T) {

	input := `- 10:45 TEAMCITY: teamcity build for proxy - [https://teamcity.example.com/project.html?projectId=someproject&tab=projectOverview](https://teamcity.example.com/project.html?projectId=someproject&tab=projectOverview)`

	r, err := NewObsidianRenderer(".")
	assert.NoError(t, err)

	r.AddTitles([]string{"teamcity"})

	actual := linkify(r.terms, input)

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
