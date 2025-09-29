package filter

import (
	"hirsi/enhancement"
	"hirsi/message"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newMessage(m string, tags ...message.Tag) *message.Message {
	msg := &message.Message{
		Message: m,
	}
	msg.AddTags(tags...)
	return msg
}

func TestHasTagWithPrefix(t *testing.T) {

	matching := newMessage(
		"with tag",
		message.Tag{Key: enhancement.PwdTag, Value: "/home/andy/dev/secret/project"},
	)
	nonMatching := newMessage(
		"without tag",
		message.Tag{Key: enhancement.PwdTag, Value: "/home/andy/dev/projects/hirsi"},
	)

	noTag := newMessage("no tags")

	assert.True(t, HasTagWithPrefix("pwd", "/home/andy/dev/secret")(matching))
	assert.False(t, HasTagWithPrefix("pwd", "/home/andy/dev/secret")(nonMatching))
	assert.False(t, HasTagWithPrefix("pwd", "/home/andy/dev/secret")(noTag))

}

func TestHasTagWithoutPrefix(t *testing.T) {

	matching := newMessage(
		"with tag",
		message.Tag{Key: enhancement.PwdTag, Value: "/home/andy/dev/secret/project"},
	)

	nonMatching := newMessage(
		"without tag",
		message.Tag{Key: enhancement.PwdTag, Value: "/home/andy/dev/projects/hirsi"},
	)

	noTag := newMessage("no tags")

	assert.False(t, HasTagWithoutPrefix("pwd", "/home/andy/dev/secret")(matching))
	assert.True(t, HasTagWithoutPrefix("pwd", "/home/andy/dev/secret")(nonMatching))
	assert.False(t, HasTagWithoutPrefix("pwd", "/home/andy/dev/secret")(noTag))

}
