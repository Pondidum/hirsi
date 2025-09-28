package message

type Tags struct {
	tags []Tag
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
	t.tags = append(t.tags, NewTag(key, value))
}

type Tag struct {
	Key   string
	Value string
}

func NewTag(k, v string) Tag {
	return Tag{k, v}
}
