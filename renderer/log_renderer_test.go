package renderer

import (
	"bytes"
	"hirsi/message"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogRendererMessage(t *testing.T) {

	r, err := NewLogRenderer(".")
	assert.NoError(t, err)

	m := &message.Message{
		WrittenAt: time.Date(2024, 06, 26, 14, 3, 37, 0, time.UTC),
		Message:   "this is a test",
		Tags: map[string]string{
			"pwd":  "/home/andy",
			"type": "test",
		},
	}

	buf := &bytes.Buffer{}
	assert.NoError(t, r.writeMessage(buf, m))

	assert.Equal(t, "- 14:03 #pwd: /home/andy #type: test: this is a test\n", buf.String())
}
