package message

import "time"

type Message struct {
	WrittenAt time.Time
	StoredAt  time.Time
	Message   string
	Tags      []Tag
}

type Tag struct {
	Key   string
	Value string
}

func NewTag(k, v string) Tag {
	return Tag{k, v}
}
