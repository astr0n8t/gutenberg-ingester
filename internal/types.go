package internal

import (
	"encoding/xml"
)

type Collection struct {
	records []BookRecord `xml:"collection"`
}

// Put any types in here
type BookRecord struct {
	XMLName xml.Name `xml:"record"`
}
