package obsidian

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkifyingPhrases(t *testing.T) {

	input := `- 10:45 TEAMCITY: teamcity build for proxy - [https://teamcity.example.com/project.html?projectId=someproject&tab=projectOverview](https://teamcity.example.com/project.html?projectId=someproject&tab=projectOverview)`

	teamcity, _ := NewTerm("teamcity", "")
	terms := []*Term{
		teamcity,
	}

	actual := linkify(terms, input)

	expected := "- 10:45 [[TEAMCITY]]: [[teamcity]] build for proxy - [https://teamcity.example.com/project.html?projectId=someproject&tab=projectOverview](https://teamcity.example.com/project.html?projectId=someproject&tab=projectOverview)"
	assert.Equal(t, expected, actual)
}
