package message

import (
	"iter"
	"time"
)

type Message struct {
	WrittenAt time.Time
	StoredAt  time.Time
	Message   string
	tags      []Tag
}

type Tag struct {
	Key   string
	Value string
}

func (m *Message) AddTags(tags ...Tag) {
	m.tags = append(m.tags, tags...)
}

func (m *Message) SetTag(key string, value string) {

	tags := make([]Tag, 0, len(m.tags))

	for _, tag := range m.tags {
		if tag.Key != key {
			tags = append(tags, tag)
		}
	}

	m.tags = append(tags, NewTag(key, value))
}

func (m *Message) AddTag(key string, value string) {
	m.AddTagT(NewTag(key, value))
}

func (m *Message) AddTagT(tag Tag) {
	m.tags = append(m.tags, tag)
}

func (m *Message) Tags() iter.Seq[Tag] {
	return func(yield func(Tag) bool) {
		for _, tag := range m.tags {
			if !yield(tag) {
				return
			}
		}
	}
}

func NewTag(k, v string) Tag {
	return Tag{k, v}
}
