package collection

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestParseXML(t *testing.T) {
	xmlFile, err := os.Open("../../tests/example.xml")
	if err != nil {
		t.Errorf("Unable to open tests/example.xml %v", err)
	}

	bytes, _ := ioutil.ReadAll(xmlFile)

	var collection Collection

	err = xml.Unmarshal(bytes, &collection)

	if err != nil {
		t.Errorf("Issue unmarshling XML: %v", err)
	}

	for _, r := range collection.Records {
		id, idErr := r.Id()
		if idErr != nil {
			t.Errorf("Issue unmarshling XML: %v", idErr)
		}
		fmt.Printf("ID is %v\n", id)
		title, titleErr := r.Name()
		if titleErr != nil {
			t.Errorf("Issue unmarshling XML: %v", titleErr)
		}
		fmt.Printf("Title is %v\n", title)
		url, urlErr := r.URL()
		if urlErr != nil {
			t.Errorf("Issue unmarshling XML: %v", urlErr)
		}
		fmt.Printf("URL is %v\n", url)
		lang, langErr := r.Language()
		if langErr != nil {
			t.Errorf("Issue unmarshling XML: %v", langErr)
		}
		fmt.Printf("Language is %v\n", lang)
	}
}
