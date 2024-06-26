package renderer

import (
	"hirsi/message"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompositeRenderer(t *testing.T) {

	renderers := make([]Renderer, 10)
	for i := 0; i < len(renderers); i++ {
		renderers[i] = &MemoryRenderer{}
	}

	composite := NewCompositeRenderer(renderers...)

	err := composite.Render(&message.Message{
		Message: "the message",
	})
	assert.NoError(t, err)

	for _, r := range renderers {
		mr := r.(*MemoryRenderer)
		assert.Len(t, mr.messages, 1)
		assert.Equal(t, "the message", mr.messages[0].Message)
	}
}
