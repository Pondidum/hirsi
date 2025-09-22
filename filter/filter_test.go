package filter

import (
	"hirsi/enhancement"
	"hirsi/message"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasTagWithPrefix(t *testing.T) {

	matching := &message.Message{
		Message: "with tag",
		Tags: []message.Tag{
			message.Tag{Key: enhancement.PwdTag, Value: "/home/andy/dev/secret/project"},
		},
	}
	nonMatching := &message.Message{
		Message: "without tag",
		Tags: []message.Tag{
			message.Tag{Key: enhancement.PwdTag, Value: "/home/andy/dev/projects/hirsi"},
		},
	}

	noTag := &message.Message{
		Message: "no tags",
		Tags:    []message.Tag{},
	}

	assert.True(t, HasTagWithPrefix("pwd", "/home/andy/dev/secret")(matching))
	assert.False(t, HasTagWithPrefix("pwd", "/home/andy/dev/secret")(nonMatching))
	assert.False(t, HasTagWithPrefix("pwd", "/home/andy/dev/secret")(noTag))

}

func TestHasTagWithoutPrefix(t *testing.T) {

	matching := &message.Message{
		Message: "with tag",
		Tags: []message.Tag{
			message.Tag{Key: enhancement.PwdTag, Value: "/home/andy/dev/secret/project"},
		},
	}
	nonMatching := &message.Message{
		Message: "without tag",
		Tags: []message.Tag{
			message.Tag{Key: enhancement.PwdTag, Value: "/home/andy/dev/projects/hirsi"},
		},
	}

	noTag := &message.Message{
		Message: "no tags",
		Tags:    []message.Tag{},
	}

	assert.False(t, HasTagWithoutPrefix("pwd", "/home/andy/dev/secret")(matching))
	assert.True(t, HasTagWithoutPrefix("pwd", "/home/andy/dev/secret")(nonMatching))
	assert.False(t, HasTagWithoutPrefix("pwd", "/home/andy/dev/secret")(noTag))

}
