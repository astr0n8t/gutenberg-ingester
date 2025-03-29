package collection

import (
	"encoding/xml"
)

// Collection was generated 2025-03-29 10:21:01 by https://xml-to-go.github.io/ in Ukraine.
type Collection struct {
	XMLName xml.Name `xml:"collection"`
	Text    string   `xml:",chardata"`
	Xmlns   string   `xml:"xmlns,attr"`
	Records  []CollectionRecord  `xml:"record"`
}

type CollectionRecord struct {
	Text         string `xml:",chardata"`
	Leader       string `xml:"leader"`
	Controlfields []struct {
		Text string `xml:",chardata"`
		Tag  string `xml:"tag,attr"`
	} `xml:"controlfield"`
	Datafields []struct {
		Text     string `xml:",chardata"`
		Ind1     string `xml:"ind1,attr"`
		Ind2     string `xml:"ind2,attr"`
		Tag      string `xml:"tag,attr"`
		Subfields []struct {
			Text string `xml:",chardata"`
			Code string `xml:"code,attr"`
		} `xml:"subfield"`
	} `xml:"datafield"`
}

