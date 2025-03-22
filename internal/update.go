package internal

import (
	"encoding/xml"
	"github.com/astr0n8t/gutenberg-ingester/internal/rss"
	"io/ioutil"
	"net/http"
)

func pullRSS(url string) (rss.RSS, error) {
	var rss rss.RSS

	resp, err := http.Get(url)
	if err != nil {
		return rss, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = xml.Unmarshal(body, &rss)

	if err != nil {
		return rss, err
	}

	return rss, nil
}
