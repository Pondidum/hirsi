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
