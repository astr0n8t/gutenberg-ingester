package rdf

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Returns the ID of the record
func (r *RDF) Id() (int, error) {
	id := -1
	var err error

	urlParts := strings.Split(r.Ebook.About, `/`)
	idStr := urlParts[len(urlParts)-1]
	id, err = strconv.Atoi(idStr)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// Returns the Name of the record
func (r *RDF) Name() (string, error) {
	return r.Ebook.Title, nil
}

// Returns the language of the record
func (r *RDF) Language() (string, error) {
	languages := map[string]string{
		"en": "english",
		"es": "spanish",
		"fr": "french",
		"de": "german",
		"fi": "finnish",
		"nl": "dutch",
		"it": "italian",
		"pt": "portuguese",
		"eo": "esperanto",
	}

	lang := r.Ebook.Language.Description.Value.Text

	normalizedLang, ok := languages[lang]

	if !ok {
		id, _ := r.Id()
		return "", fmt.Errorf("failed to process language for rdf item: %v with lang: %v", id, lang)

	}

	return normalizedLang, nil
}

// Returns the URL of the record
func (r *RDF) URL() (string, error) {
	rdfURL := ""

	for _, format := range r.Ebook.HasFormat {
		if len(format.File.Format) == 1 && format.File.Format[0].Description.Value.Text == "application/rdf+xml" {
			rdfURL = format.File.About
			break
		}
	}

	if len(rdfURL) == 0 {
		id, _ := r.Id()
		return "", fmt.Errorf("failed to get URL for rdf item: %v", id)

	}

	urlParts := strings.Split(rdfURL, `/`)
	urlParts = urlParts[:len(urlParts)-1]
	id, _ := r.Id()
	url := strings.Join(urlParts[:], "/") + "/" + strconv.Itoa(id)

	return url, nil
}

// Returns a map of download types and their urls
func (r *RDF) Formats() (map[string]string, error) {
	downloadMap := make(map[string]string)

	for _, format := range r.Ebook.HasFormat {

		url := format.File.About

		// ignore readme files
		if strings.Contains(url, "readme") {
			continue
		}

		urlParts := strings.Split(url, `/`)
		if len(urlParts) < 3 {
			log.Printf("could not parse url for format with unknown syntax: %v", url)
			continue
		}

		fileFormat := urlParts[len(urlParts)-1]
		fileFormatParts := strings.Split(fileFormat, `.`)
		fileFormat = ""
		for i := 1; i < len(fileFormatParts); i++ {
			fileFormat += "." + fileFormatParts[i]
		}

		urlSlug := ""
		for i := 3; i < len(urlParts); i++ {
			urlSlug += "/" + urlParts[i]
		}

		downloadMap[fileFormat] = urlSlug
	}

	return downloadMap, nil
}
