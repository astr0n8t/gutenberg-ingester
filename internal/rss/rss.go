package rss 

import (
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

// Returns the Title of the record
func (b *BookItem) Title() (string, error) {
	return b.Name, nil
}

// Returns the URL of the record
func (b *BookItem) URL() (string, error) {
	return b.Link, nil
}
