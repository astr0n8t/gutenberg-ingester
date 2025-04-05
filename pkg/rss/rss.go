package rss

import (
	"fmt"
	"strconv"
	"strings"
)

// Returns the ID of the record
func (b *BookItem) Id() (int, error) {
	id := -1
	var err error

	urlParts := strings.Split(b.Link, `/`)
	idStr := urlParts[len(urlParts)-1]
	id, err = strconv.Atoi(idStr)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// Returns the Name of the record
func (b *BookItem) Name() (string, error) {
	return b.Title, nil
}

// Returns the URL of the record
func (b *BookItem) URL() (string, error) {
	return b.Link, nil
}

// Returns the language of the record
func (b *BookItem) Language() (string, error) {
	descParts := strings.Split(b.Description, `Language: `)
	lang := descParts[len(descParts)-1]

	languages := map[string]string{
		"English":    "english",
		"French":     "french",
		"Spanish":    "spanish",
		"German":     "german",
		"Finnish":    "finnish",
		"Dutch":      "dutch",
		"Italian":    "italian",
		"Portuguese": "portuguese",
		"Esperanto":  "esperanto",
	}

	normalizedLang, ok := languages[lang]

	if !ok {
		id, _ := b.Id()
		return "", fmt.Errorf("failed to process language for rss item: %v with lang: %v", id, lang)

	}

	return normalizedLang, nil
}
