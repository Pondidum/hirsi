package obsidian

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

type Term struct {
	Name  string
	Path  string
	regex *regexp.Regexp
}

func NewTerm(name string, path string) (*Term, error) {
	rx, err := regexp.Compile(`(?i)^(\W*?)(` + name + `)(\W*?)$`)
	if err != nil {
		return nil, err
	}

	return &Term{
		Name:  name,
		Path:  path,
		regex: rx,
	}, nil
}

func (t *Term) Linkify(word string) (string, bool) {
	prefix := ""
	if filename := path.Base(t.Path); !strings.EqualFold(filename, word) {
		prefix = t.Path + "|"
	}

	if word == t.Name {
		return "[[" + prefix + word + "]]", true
	}

	if t.regex.MatchString(word) {
		groups := t.regex.FindStringSubmatch(word)

		return fmt.Sprintf("%s[[%s%s]]%s", groups[1], prefix, groups[2], groups[3]), true
	}

	return word, false
}

func linkify(terms []*Term, message string) string {

	words := strings.Split(message, " ")
	for i, word := range words {
		for _, term := range terms {

			if replacement, ok := term.Linkify(word); ok {
				words[i] = replacement
				continue
			}
		}
	}

	return strings.Join(words, " ")
}
