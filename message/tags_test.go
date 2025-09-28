package message

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTagAdd(t *testing.T) {
	tags := Tags{}

	tags.Add("key", "one")

	require.Equal(t, []Tag{
		Tag{"key", "one"},
	}, tags.tags)

	tags.Add("key", "two")

	require.Equal(t, []Tag{
		Tag{"key", "one"},
		Tag{"key", "two"},
	}, tags.tags)
}

func TestTagSet(t *testing.T) {
	tags := Tags{}

	tags.Set("key", "one")

	require.Equal(t, []Tag{
		Tag{"key", "one"},
	}, tags.tags)

	tags.Set("key", "two")

	require.Equal(t, []Tag{
		Tag{"key", "two"},
	}, tags.tags)

	tags.Add("key", "three")
	tags.Add("key", "four")

	tags.Set("key", "five")

	require.Equal(t, []Tag{
		Tag{"key", "five"},
	}, tags.tags)

}

func TestTagIteration(t *testing.T) {
	tags := Tags{}
	tags.Add("1", "one")
	tags.Add("2", "two")
	tags.Add("3", "three")

	seen := map[string]string{}
	for tag := range tags.All() {
		seen[tag.Key] = tag.Value
	}

	require.Equal(t, map[string]string{
		"1": "one",
		"2": "two",
		"3": "three",
	}, seen)
}
