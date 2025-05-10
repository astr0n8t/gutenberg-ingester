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
		"en":  "english",
		"es":  "spanish",
		"fr":  "french",
		"de":  "german",
		"fi":  "finnish",
		"nl":  "dutch",
		"it":  "italian",
		"pt":  "portuguese",
		"eo":  "esperanto",
		"af":  "afrikaans",
		"ale": "aleut",
		"ang": "old english",
		"ar":  "arabic",
		"arp": "arapaho",
		"bg":  "bulgarian",
		"bgs": "tagabawa",
		"br":  "breton",
		"brx": "bodo",
		"ca":  "catalan",
		"ceb": "cebuano",
		"cs":  "czech",
		"csb": "kashubian",
		"cy":  "welsh",
		"da":  "danish",
		"el":  "greek",
		"enm": "middle english",
		"et":  "estonian",
		"fa":  "persian",
		"fur": "friulian",
		"fy":  "frisian",
		"ga":  "irish",
		"gl":  "galician",
		"gla": "scottish gaelic",
		"grc": "ancient greek",
		"hai": "haida",
		"he":  "hebrew",
		"hu":  "hungarian",
		"ia":  "interlingua",
		"ilo": "ilocano",
		"is":  "icelandic",
		"iu":  "inuktitut",
		"ja":  "japanese",
		"kha": "khasi",
		"kld": "gamilaraay",
		"ko":  "korean",
		"la":  "latin",
		"lt":  "lithuanian",
		"mi":  "māori",
		"myn": "mayan languages",
		"nah": "nahuatl",
		"nai": "north american indian",
		"nap": "neapolitan",
		"nav": "navajo",
		"no":  "norwegian",
		"oc":  "occitan",
		"oji": "ojibwe",
		"pl":  "polish",
		"rmq": "caló",
		"ro":  "romanian",
		"ru":  "russian",
		"sa":  "sanskrit",
		"sco": "scots",
		"sl":  "slovenian",
		"sr":  "serbian",
		"sv":  "swedish",
		"te":  "telugu",
		"tl":  "tagalog",
		"yi":  "yiddish",
		"zh":  "chinese",
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
