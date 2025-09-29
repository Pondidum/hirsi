package message

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTagAdd(t *testing.T) {
	message := Message{}

	message.AddTag("key", "one")

	require.Equal(t, []Tag{
		Tag{"key", "one"},
	}, message.tags)

	message.AddTag("key", "two")

	require.Equal(t, []Tag{
		Tag{"key", "one"},
		Tag{"key", "two"},
	}, message.tags)
}

func TestTagSet(t *testing.T) {
	message := Message{}

	message.SetTag("key", "one")

	require.Equal(t, []Tag{
		Tag{"key", "one"},
	}, message.tags)

	message.SetTag("key", "two")

	require.Equal(t, []Tag{
		Tag{"key", "two"},
	}, message.tags)

	message.AddTag("key", "three")
	message.AddTag("key", "four")

	message.SetTag("key", "five")

	require.Equal(t, []Tag{
		Tag{"key", "five"},
	}, message.tags)

}

func TestTagIteration(t *testing.T) {
	message := Message{}
	message.AddTag("1", "one")
	message.AddTag("2", "two")
	message.AddTag("3", "three")

	seen := map[string]string{}
	for tag := range message.Tags() {
		seen[tag.Key] = tag.Value
	}

	require.Equal(t, map[string]string{
		"1": "one",
		"2": "two",
		"3": "three",
	}, seen)
}
