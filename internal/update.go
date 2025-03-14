package internal

import (
	"net/http"
	"encoding/xml"
	"io/ioutil"
)

func pullRSS() (RSS, error) {
	var rss RSS

	resp, err := http.Get("https://www.gutenberg.org/cache/epub/feeds/today.rss")
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
