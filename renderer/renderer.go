package renderer

import "hirsi/message"

type Renderer interface {
	Render(message *message.Message) error
}
