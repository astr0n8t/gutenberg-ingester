package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/astr0n8t/gutenberg-ingester/pkg/rdf"
)

func (r *Runtime) downloadFromRDF(rdf rdf.RDF) error {
	formats, formatErr := rdf.Formats()
	if formatErr != nil {
		return fmt.Errorf("could not get format list for record with error: %v", formatErr)
	}

	downloads := make(map[string]string, 0)

	download_type := r.Config.GetStringSlice("download_type")
	download_type_precedence := r.Config.GetString("download_type_precedence")

	mirror_url := r.Config.GetString("gutenberg_mirror_url")
	mirror_url = strings.TrimSuffix(mirror_url, "/")

	if download_type_precedence == "sequential" {
		for _, requestedFormat := range download_type {
			if urlSlug, hasFormat := formats[requestedFormat]; hasFormat {
				url := mirror_url + urlSlug
				downloads[requestedFormat] = url
				break
			}
		}
	} else if download_type_precedence == "parallel" {
		for _, requestedFormat := range download_type {
			if urlSlug, hasFormat := formats[requestedFormat]; hasFormat {
				url := mirror_url + urlSlug
				downloads[requestedFormat] = url
			}
		}
	} else {
		log.Fatalf("unknown setting for download type precedence: %v", download_type_precedence)
	}

	id, idErr := rdf.Id()
	if idErr != nil {
		return fmt.Errorf("unable to get ID for record: %v", idErr)
	}

	idStr := strconv.Itoa(id)
	download_location := r.Config.GetString("download_location")
	download_delay := time.Duration(r.Config.GetInt("download_delay")) * time.Second

	for format, url := range downloads {
		if format == ".epub3.images" && r.Config.GetBool("epub_use_proper_extension") {
			format = ".epub"
		}
		filename := strings.TrimSuffix(download_location, "/") + "/" + idStr + format
		downloadErr := downloadFromURLToFile(url, filename)
		if downloadErr != nil {
			return fmt.Errorf("unable to download format: %v for record: %v with error: %v", format, id, downloadErr)
		}

		time.Sleep(download_delay)
		r.DB.SetDownloaded(id)
	}

	return nil
}

func downloadFromURLToFile(fileUrl string, filename string) error {
	// 1. Fetch the initial URL and handle redirection (302)
	resp, getErr := http.Get(fileUrl)
	if getErr != nil {
		return fmt.Errorf("error fetching URL: %w", getErr)
	}
	defer resp.Body.Close()

	// Check for redirection (302)
	if resp.StatusCode == http.StatusFound { // StatusFound is 302
		locationURL, locationHeaderErr := url.Parse(resp.Header.Get("Location"))
		if locationHeaderErr != nil {
			return fmt.Errorf("error parsing redirection URL: %w", locationHeaderErr)
		}
		// Handle relative URLs, making sure they're absolute.
		if !locationURL.IsAbs() {
			baseURL, _ := url.Parse(fileUrl) // Use the original URL for the base.
			locationURL = baseURL.ResolveReference(locationURL)
		}
		resp, getErr := http.Get(locationURL.String()) // Fetch the redirected URL
		if getErr != nil {
			return fmt.Errorf("error fetching redirected URL: %w", getErr)
		}
		defer resp.Body.Close()
	}

	// Check for successful status code
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 2. Create the output file
	out, fileWriteErr := os.Create(filename)
	if fileWriteErr != nil {
		return fmt.Errorf("error creating file: %w", fileWriteErr)
	}
	defer out.Close()

	// 3. Copy the response body to the output file
	_, copyErr := io.Copy(out, resp.Body)
	if copyErr != nil {
		return fmt.Errorf("error copying response body: %w", copyErr)
	}

	return nil
}

func downloadFromURLToBytes(fileUrl string) (*bytes.Buffer, error) {
	// 1. Fetch the initial URL and handle redirection (302)
	resp, getErr := http.Get(fileUrl)
	if getErr != nil {
		return nil, fmt.Errorf("error fetching URL: %w", getErr)
	}
	defer resp.Body.Close()

	// Check for redirection (302)
	if resp.StatusCode == http.StatusFound { // StatusFound is 302
		locationURL, locationHeaderErr := url.Parse(resp.Header.Get("Location"))
		if locationHeaderErr != nil {
			return nil, fmt.Errorf("error parsing redirection URL: %w", locationHeaderErr)
		}
		// Handle relative URLs, making sure they're absolute.
		if !locationURL.IsAbs() {
			baseURL, _ := url.Parse(fileUrl) // Use the original URL for the base.
			locationURL = baseURL.ResolveReference(locationURL)
		}
		resp, getErr := http.Get(locationURL.String()) // Fetch the redirected URL
		if getErr != nil {
			return nil, fmt.Errorf("error fetching redirected URL: %w", getErr)
		}
		defer resp.Body.Close()
	}

	// Check for successful status code
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	// 3. Copy the response body to the output file
	_, copyErr := io.Copy(buf, resp.Body)
	if copyErr != nil {
		return nil, fmt.Errorf("error copying response body: %w", copyErr)
	}

	return buf, nil
}
