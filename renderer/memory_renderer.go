package renderer

import "hirsi/message"

type MemoryRenderer struct {
	Messages []*message.Message
}

func (r *MemoryRenderer) Render(message *message.Message) error {
	r.Messages = append(r.Messages, message)
	return nil
}
