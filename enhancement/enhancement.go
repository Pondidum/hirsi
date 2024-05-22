package enhancement

import "hirsi/message"

type Enhancement interface {
	Enhance(message *message.Message) error
}
