package message

import "time"

type Message struct {
	WrittenAt time.Time
	Message   string
	Tags      map[string]string
}
