package rdf

import (
	"encoding/xml"
)

// RDF was generated 2025-03-22 17:50:53 by https://xml-to-go.github.io/ in Ukraine.
type RDF struct {
	XMLName xml.Name `xml:"RDF"`
	Text    string   `xml:",chardata"`
	Base    string   `xml:"base,attr"`
	Dcterms string   `xml:"dcterms,attr"`
	Pgterms string   `xml:"pgterms,attr"`
	Cc      string   `xml:"cc,attr"`
	Rdfs    string   `xml:"rdfs,attr"`
	Rdf     string   `xml:"rdf,attr"`
	Dcam    string   `xml:"dcam,attr"`
	Ebook   struct {
		Text      string `xml:",chardata"`
		About     string `xml:"about,attr"`
		Publisher string `xml:"publisher"`
		License   struct {
			Text     string `xml:",chardata"`
			Resource string `xml:"resource,attr"`
		} `xml:"license"`
		Issued struct {
			Text     string `xml:",chardata"`
			Datatype string `xml:"datatype,attr"`
		} `xml:"issued"`
		Rights    string `xml:"rights"`
		Downloads struct {
			Text     string `xml:",chardata"`
			Datatype string `xml:"datatype,attr"`
		} `xml:"downloads"`
		Creator []struct {
			Text  string `xml:",chardata"`
			Agent struct {
				Text    string `xml:",chardata"`
				About   string `xml:"about,attr"`
				Name    string `xml:"name"`
				Webpage struct {
					Text     string `xml:",chardata"`
					Resource string `xml:"resource,attr"`
				} `xml:"webpage"`
			} `xml:"agent"`
		} `xml:"creator"`
		Title       string `xml:"title"`
		Description string `xml:"description"`
		Language    struct {
			Text        string `xml:",chardata"`
			Description struct {
				Text   string `xml:",chardata"`
				NodeID string `xml:"nodeID,attr"`
				Value  struct {
					Text     string `xml:",chardata"`
					Datatype string `xml:"datatype,attr"`
				} `xml:"value"`
			} `xml:"Description"`
		} `xml:"language"`
		Subject []struct {
			Text        string `xml:",chardata"`
			Description struct {
				Text     string `xml:",chardata"`
				NodeID   string `xml:"nodeID,attr"`
				MemberOf struct {
					Text     string `xml:",chardata"`
					Resource string `xml:"resource,attr"`
				} `xml:"memberOf"`
				Value string `xml:"value"`
			} `xml:"Description"`
		} `xml:"subject"`
		Type struct {
			Text        string `xml:",chardata"`
			Description struct {
				Text     string `xml:",chardata"`
				NodeID   string `xml:"nodeID,attr"`
				MemberOf struct {
					Text     string `xml:",chardata"`
					Resource string `xml:"resource,attr"`
				} `xml:"memberOf"`
				Value string `xml:"value"`
			} `xml:"Description"`
		} `xml:"type"`
		Bookshelf []struct {
			Text        string `xml:",chardata"`
			Description struct {
				Text     string `xml:",chardata"`
				NodeID   string `xml:"nodeID,attr"`
				MemberOf struct {
					Text     string `xml:",chardata"`
					Resource string `xml:"resource,attr"`
				} `xml:"memberOf"`
				Value string `xml:"value"`
			} `xml:"Description"`
		} `xml:"bookshelf"`
		HasFormat []struct {
			Text string `xml:",chardata"`
			File struct {
				Text       string `xml:",chardata"`
				About      string `xml:"about,attr"`
				IsFormatOf struct {
					Text     string `xml:",chardata"`
					Resource string `xml:"resource,attr"`
				} `xml:"isFormatOf"`
				Extent struct {
					Text     string `xml:",chardata"`
					Datatype string `xml:"datatype,attr"`
				} `xml:"extent"`
				Modified struct {
					Text     string `xml:",chardata"`
					Datatype string `xml:"datatype,attr"`
				} `xml:"modified"`
				Format []struct {
					Text        string `xml:",chardata"`
					Description struct {
						Text     string `xml:",chardata"`
						NodeID   string `xml:"nodeID,attr"`
						MemberOf struct {
							Text     string `xml:",chardata"`
							Resource string `xml:"resource,attr"`
						} `xml:"memberOf"`
						Value struct {
							Text     string `xml:",chardata"`
							Datatype string `xml:"datatype,attr"`
						} `xml:"value"`
					} `xml:"Description"`
				} `xml:"format"`
			} `xml:"file"`
		} `xml:"hasFormat"`
	} `xml:"ebook"`
	Work struct {
		Text    string `xml:",chardata"`
		About   string `xml:"about,attr"`
		License struct {
			Text     string `xml:",chardata"`
			Resource string `xml:"resource,attr"`
		} `xml:"license"`
		Comment string `xml:"comment"`
	} `xml:"Work"`
	Description struct {
		Text        string `xml:",chardata"`
		About       string `xml:"about,attr"`
		Description string `xml:"description"`
	} `xml:"Description"`
}
