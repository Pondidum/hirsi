package message

import "iter"

type Tags struct {
	tags []Tag
}

func NewTagsFrom(tags []Tag) *Tags {
	return &Tags{tags: tags}
}

func (t *Tags) Set(key string, value string) {

	tags := make([]Tag, 0, len(t.tags))

	for _, tag := range t.tags {
		if tag.Key != key {
			tags = append(tags, tag)
		}
	}

	t.tags = append(tags, NewTag(key, value))
}

func (t *Tags) Add(key string, value string) {
	t.AddT(NewTag(key, value))
}

func (t *Tags) AddT(tag Tag) {
	t.tags = append(t.tags, tag)
}

func (t *Tags) All() iter.Seq[Tag] {
	return func(yield func(Tag) bool) {
		for _, tag := range t.tags {
			if !yield(tag) {
				return
			}
		}
	}
}

type Tag struct {
	Key   string
	Value string
}

func NewTag(k, v string) Tag {
	return Tag{k, v}
}
