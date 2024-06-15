package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObsidianLinkin(t *testing.T) {

	input := `- 10:45 TEAMCITY: teamcity build for proxy - [https://teamcity.example.com/project.html?projectId=someproject&tab=projectOverview](https://teamcity.example.com/project.html?projectId=someproject&tab=projectOverview)`

	r, err := NewObsidianRenderer(".")
	assert.NoError(t, err)

	actual := r.linkify(input)

	expected := "- 10:45 [[TEAMCITY]]: [[teamcity]] build for proxy - [https://teamcity.example.com/project.html?projectId=someproject&tab=projectOverview](https://teamcity.example.com/project.html?projectId=someproject&tab=projectOverview)"
	assert.Equal(t, expected, actual)
}
