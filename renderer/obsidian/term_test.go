package obsidian

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkifyingPhrases(t *testing.T) {

	input := `- 10:45 TEAMCITY: teamcity build for proxy - [https://teamcity.example.com/project.html?projectId=someproject&tab=projectOverview](https://teamcity.example.com/project.html?projectId=someproject&tab=projectOverview)`

	teamcity, _ := NewTerm("teamcity", "ci/teamcity")
	terms := []*Term{
		teamcity,
	}

	actual := linkify(terms, input)

	expected := "- 10:45 [[ci/teamcity|TEAMCITY]]: [[teamcity]] build for proxy - [https://teamcity.example.com/project.html?projectId=someproject&tab=projectOverview](https://teamcity.example.com/project.html?projectId=someproject&tab=projectOverview)"
	assert.Equal(t, expected, actual)
}

func TestLinking(t *testing.T) {

	t.Run("plain", func(t *testing.T) {

		term, err := NewTerm("part-one", "projects/big-project/part-one")
		assert.NoError(t, err)

		replacement, changed := term.Linkify(term.Name)
		assert.True(t, changed)
		assert.Equal(t, "[[part-one]]", replacement)
	})
	t.Run("alias", func(t *testing.T) {

		term, err := NewTerm("other-name", "projects/other-project/duplicate")
		assert.NoError(t, err)

		replacement, changed := term.Linkify(term.Name)
		assert.True(t, changed)
		assert.Equal(t, "[[projects/other-project/duplicate|other-name]]", replacement)
	})

	t.Run("alias special", func(t *testing.T) {

		term, err := NewTerm("other-name", "projects/other-project/duplicate")
		assert.NoError(t, err)

		replacement, changed := term.Linkify("other-name:")
		assert.True(t, changed)
		assert.Equal(t, "[[projects/other-project/duplicate|other-name]]:", replacement)
	})
}
