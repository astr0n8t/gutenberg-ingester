package collection

import (
	"fmt"
	"strconv"
)

// Testing method to make sure parse worked
func (b *CollectionRecord) dumpRecord() {
	for _, d := range b.Datafields {
		fmt.Printf("ind1: %v\n", d.Ind1)
		fmt.Printf("ind2: %v\n", d.Ind2)
		fmt.Printf("tag: %v\n", d.Tag)
		for _, s := range d.Subfields {
			fmt.Printf("subfield: %v\n", s.Text)

		}
	}
}

// Returns the ID of the record
func (b *CollectionRecord) Id() (int, error) {
	id := -1
	var err error

	for _, c := range b.Controlfields {
		if c.Tag == "001" {
			id, err = strconv.Atoi(c.Text)
			if err != nil {
				return -1, fmt.Errorf("failed to process id for record: %v", id)
			}
			break
		}
	}

	return id, nil
}

// Returns the Name of the record
func (b *CollectionRecord) Name() (string, error) {
	title := ""

	for _, d := range b.Datafields {
		if d.Ind1 == "1" && d.Ind2 == "4" && d.Tag == "245" {
			if len(d.Subfields) == 0 {
				id, _ := b.Id()
				return "", fmt.Errorf("failed to process title for record: %v", id)
			} else {
				for _, t := range d.Subfields {
					title += " " + t.Text
				}
				break
			}
		}
	}

	return title, nil
}

// Returns the language of the record
func (b *CollectionRecord) Language() (string, error) {
	lang := ""

	// Go in reverse since its usually the last datafield
	for i := len(b.Datafields) - 1; i >= 0; i-- {
		d := b.Datafields[i]
		if d.Ind1 == " " && d.Ind2 == "7" && d.Tag == "041" {
			for _, t := range d.Subfields {
				if t.Code == "a" {
					lang = t.Text
					break
				}
			}
			if lang != "" {
				break
			} else {
				id, _ := b.Id()
				return "", fmt.Errorf("failed to process language for record: %v", id)
			}
		}
	}

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

	normalizedLang, ok := languages[lang]

	if !ok {
		id, _ := b.Id()
		return "", fmt.Errorf("failed to process language for record: %v with lang: %v", id, lang)

	}

	return normalizedLang, nil
}

// Returns the URL of the record
func (b *CollectionRecord) URL() (string, error) {
	url := ""

	// Go in reverse since its usually the last datafield
	for i := len(b.Datafields) - 1; i >= 0; i-- {
		d := b.Datafields[i]
		if d.Ind1 == "4" && d.Ind2 == "0" && d.Tag == "856" {
			if len(d.Subfields) != 1 {
				id, _ := b.Id()
				return "", fmt.Errorf("failed to process url for record: %v", id)
			} else {
				url = d.Subfields[0].Text
				break
			}
		}
	}

	return url, nil
}
