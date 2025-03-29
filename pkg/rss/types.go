package rss

import (
	"encoding/xml"
)

// RSS was generated 2025-03-29 10:02:12 by https://xml-to-go.github.io/ in Ukraine.
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Text          string     `xml:",chardata"`
		Title         string     `xml:"title"`
		Link          string     `xml:"link"`
		Description   string     `xml:"description"`
		Language      string     `xml:"language"`
		WebMaster     string     `xml:"webMaster"`
		PubDate       string     `xml:"pubDate"`
		LastBuildDate string     `xml:"lastBuildDate"`
		Items         []BookItem `xml:"item"`
	} `xml:"channel"`
}

type BookItem struct {
	Text        string `xml:",chardata"`
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
}
