package records 

import (
	"fmt"
	"strconv"
)

// Testing method to make sure parse worked
func (b *BookRecord) dumpRecord() {
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
func (b *BookRecord) Id() (int, error) {
	id := -1
	var err error

	for _, c := range b.Controlfields {
		if c.Tag == "001" {
			id, err = strconv.Atoi(c.Text)
			if err != nil {
				return -1, nil
			}
			break
		}
	}

	return id, nil
}

// Returns the Title of the record
func (b *BookRecord) Title() (string, error) {
	title := ""

	for _, d := range b.Datafields {
		if d.Ind1 == "1" && d.Ind2 == "4" && d.Tag == "245" {
			if len(d.Subfields) == 0 {
				return "", fmt.Errorf("failed to process title for record")
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

// Returns the URL of the record
func (b *BookRecord) URL() (string, error) {
	url := ""

	// Go in reverse since its usually the last datafield
	for i := len(b.Datafields) - 1; i >= 0; i-- {
		d := b.Datafields[i]
		if d.Ind1 == "4" && d.Ind2 == "0" && d.Tag == "856" {
			if len(d.Subfields) != 1 {
				return "", fmt.Errorf("failed to process url for record")
			} else {
				url = d.Subfields[0].Text
				break
			}
		}
	}

	return url, nil
}
