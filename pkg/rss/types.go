package rss

import (
	"encoding/xml"
)

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
