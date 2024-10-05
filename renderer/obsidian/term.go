package obsidian

import (
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

	if word == t.Name {
		return "[[" + word + "]]", true
	}

	if t.regex.MatchString(word) {
		return t.regex.ReplaceAllString(word, "$1[[$2]]$3"), true
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
