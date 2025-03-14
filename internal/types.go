package internal

import (
	"encoding/xml"
)

type Book interface {
	Id() (int, error)
	Title() (string, error)
	URL() (string, error)
}

type Collection struct {
	XMLName xml.Name     `xml:"collection"`
	Records []BookRecord `xml:"record"`
}

type BookRecord struct {
	Text          string         `xml:",chardata"`
	Leader        string         `xml:"leader"`
	Controlfields []ControlField `xml:"controlfield"`
	Datafields    []DataField    `xml:"datafield"`
}

type ControlField struct {
	Text string `xml:",chardata"`
	Tag  string `xml:"tag,attr"`
}

type DataField struct {
	Text      string     `xml:",chardata"`
	Ind1      string     `xml:"ind1,attr"`
	Ind2      string     `xml:"ind2,attr"`
	Tag       string     `xml:"tag,attr"`
	Subfields []SubField `xml:"subfield"`
}

type SubField struct {
	Text string `xml:",chardata"`
	Code string `xml:"code,attr"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title         string     `xml:"title"`
	Description   string     `xml:"description"`
	Language      string     `xml:"language"`
	WebMaster     string     `xml:"webMaster"`
	PubDate       string     `xml:"pubDate"`
	LastBuildDate string     `xml:"lastBuildDate"`
	Items         []BookItem `xml:"item"`
}

type BookItem struct {
	Name        string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
}
