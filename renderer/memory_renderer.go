package renderer

import "hirsi/message"

type MemoryRenderer struct {
	messages []*message.Message
}

func (r *MemoryRenderer) Render(message *message.Message) error {
	r.messages = append(r.messages, message)
	return nil
}
