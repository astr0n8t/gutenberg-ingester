package internal

import (
	"archive/tar"
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/astr0n8t/gutenberg-ingester/pkg/rdf"
	"github.com/astr0n8t/gutenberg-ingester/pkg/rss"
)

func (r *Runtime) startSyncSchedule() {
	currentTimeUTC := time.Now().UTC()
	// duration := currentTimeUTC.Sub(r.DB.GetLastFullSync())
	// if duration.Hours()/24 >= float64(r.Config.GetInt("full_sync_frequency")) {
	// 	log.Printf("Full sync due according to schedule")
	// 	syncErr := r.fullSync()
	// 	if syncErr != nil {
	// 		log.Printf("error full syncing: %v", syncErr)
	// 	}
	// }

	duration := currentTimeUTC.Sub(r.DB.GetLastPartialSync())
	if duration.Hours() >= float64(r.Config.GetInt("partial_sync_frequency")) {
		log.Printf("Partial RSS sync due according to schedule")
		syncErr := r.partialSync()
		if syncErr != nil {
			log.Printf("error partial syncing: %v", syncErr)
		}
	}

	time.Sleep(time.Hour)
}

func (r *Runtime) fullSync() error {
	log.Printf("Starting full sync with Project Gutenberg...")
	errs := false

	catalogPullErr := r.pullCatalog()
	if catalogPullErr != nil {
		return fmt.Errorf("error syncing: %v", catalogPullErr)
	}

	zipReader, zipError := zip.OpenReader(r.Catalog)
	if zipError != nil {
		return fmt.Errorf("error creating zip reader for catalog: %v", zipError)
	}

	if len(zipReader.File) != 1 {
		return fmt.Errorf("unknown file format when opening catalog zip")
	}
	zippedTarFile, readZippedTarError := zipReader.File[0].Open()
	if readZippedTarError != nil {
		return fmt.Errorf("error opening zipped tar reader for catalog: %v", readZippedTarError)
	}

	filterByLang := true
	lang_filter := r.Config.GetStringSlice("download_languages")
	lang_map := make(map[string]bool, len(lang_filter))
	for _, v := range lang_filter {
		lang_map[v] = true
	}
	if _, ok := lang_map["all"]; ok {
		filterByLang = false
	}

	// Create a tar reader on top of the gzip reader
	tarReader := tar.NewReader(zippedTarFile)

	var rdfFile []byte

	for {
		header, tarReadErr := tarReader.Next()
		if tarReadErr == io.EOF {
			break // End of archive reached
		} else if tarReadErr != nil {
			return fmt.Errorf("error reading catalog tar archive: %v", tarReadErr)
		}

		headerParts := strings.Split(header.Name, `/`)
		if len(headerParts) != 4 {
			log.Printf("unknown file path in catalog: %v", header)
			errs = true
			continue
		}
		idStr := headerParts[2]
		// don't use the test file
		if idStr == "test" {
			continue
		}
		id, idErr := strconv.Atoi(idStr)
		if idErr != nil {
			log.Printf("issue casting id: %v to int: %v", id, idErr)
			errs = true
			continue
		}

		// Check if the file has already been downloaded
		if !r.DB.GetDownloaded(id) {
			// Open the file within the archive for reading
			var tarReadTargetFileErr error
			rdfFile, tarReadTargetFileErr = io.ReadAll(tarReader)
			if tarReadTargetFileErr != nil {
				return fmt.Errorf("error reading file from catalog: %v", tarReadTargetFileErr)
			}

			var rdf rdf.RDF
			xmlParseErr := xml.Unmarshal(rdfFile, &rdf)
			if xmlParseErr != nil {
				log.Printf("error parsing rdf of id: %v into xml into struct: %v", id, xmlParseErr)
				errs = true
				continue
			}

			if filterByLang {
				itemLang, langParseErr := rdf.Language()
				if langParseErr != nil {
					log.Printf("error getting lang from rdf of id: %v with: %v", id, langParseErr)
					errs = true
					continue
				}

				// if the lang is not in the map, skip it
				if _, ok := lang_map[itemLang]; !ok {
					continue
				}
			}

			downloadErr := r.downloadFromRDF(rdf)
			if downloadErr != nil {
				log.Printf("issue downloading record: %v with error: %v", id, downloadErr)
				errs = true
				continue
			}

		}
	}
	zipReader.Close()

	// Make sure we annotate when we finished the sync
	utcTime := time.Now().UTC()
	r.DB.SetLastFullSync(utcTime)
	r.DB.SetLastPartialSync(utcTime)

	if errs {
		log.Printf("Finished full sync with Project Gutenberg with one or more errors")
	} else {
		log.Printf("Finished full sync with Project Gutenberg with no errors")
	}

	return nil
}

func (r *Runtime) partialSync() error {
	log.Printf("Starting partial sync with Project Gutenberg...")
	errs := false

	rss, rssPullErr := r.pullRSS()
	if rssPullErr != nil {
		return fmt.Errorf("error pulling rss: %v", rssPullErr)
	}

	filterByLang := true
	lang_filter := r.Config.GetStringSlice("download_languages")
	lang_map := make(map[string]bool, len(lang_filter))
	for _, v := range lang_filter {
		lang_map[v] = true
	}
	if _, ok := lang_map["all"]; ok {
		filterByLang = false
	}

	for _, item := range rss.Channel.Items {
		id, idErr := item.Id()
		if idErr != nil {
			log.Printf("error getting id for rss item: %v", idErr)
			errs = true
			continue
		}

		// Skip item if we've already downloaded and we're not updating
		if r.DB.GetDownloaded(id) && !r.Config.GetBool("update_previously_downloaded") {
			continue
		}

		rdf, rdfPullErr := r.pullRDF(id)
		if rdfPullErr != nil {
			log.Printf("error pulling rdf for id: %v with: %v", id, rdfPullErr)
			errs = true
			continue
		}

		if filterByLang {
			itemLang, langParseErr := rdf.Language()
			if langParseErr != nil {
				log.Printf("error getting lang from rdf of id: %v with: %v", id, langParseErr)
				errs = true
				continue
			}

			// if the lang is not in the map, skip it
			if _, ok := lang_map[itemLang]; !ok {
				continue
			}
		}

		downloadErr := r.downloadFromRDF(*rdf)
		if downloadErr != nil {
			log.Printf("error downloading items for id: %v", id)
			errs = true
			continue
		}
	}

	if errs {
		log.Printf("Finished partial sync with Project Gutenberg with one or more errors")
	} else {
		log.Printf("Finished partial sync with Project Gutenberg with no errors")
	}

	return nil
}

func (r *Runtime) pullRSS() (*rss.RSS, error) {
	rssURL := strings.TrimSuffix(r.Config.GetString("gutenberg_feed_url"), "/")
	rssURL += "/cache/epub/feeds/today.rss"

	rssBytes, rssPullErr := downloadFromURLToBytes(rssURL)
	if rssPullErr != nil {
		return nil, fmt.Errorf("error pulling RSS feed: %v", rssPullErr)
	}

	var rss rss.RSS
	rssUnmarshalErr := xml.Unmarshal(rssBytes.Bytes(), &rss)
	if rssUnmarshalErr != nil {
		return nil, fmt.Errorf("error unmarshalling RSS feed: %v", rssUnmarshalErr)
	}

	return &rss, nil
}

func (r *Runtime) pullRDF(id int) (*rdf.RDF, error) {
	idStr := strconv.Itoa(id)
	rdfURL := strings.TrimSuffix(r.Config.GetString("gutenberg_mirror_url"), "/")
	rdfURL += "/ebooks/" + idStr + ".rdf"

	rdfBytes, rdfPullErr := downloadFromURLToBytes(rdfURL)
	if rdfPullErr != nil {
		return nil, fmt.Errorf("error pulling RDF for item: %v with: %v", id, rdfPullErr)
	}

	var rdf rdf.RDF
	rdfUnmarshalErr := xml.Unmarshal(rdfBytes.Bytes(), &rdf)
	if rdfUnmarshalErr != nil {
		return nil, fmt.Errorf("error unmarshalling RDF feed: %v", rdfUnmarshalErr)
	}

	return &rdf, nil
}

func (r *Runtime) pullCatalog() error {
	filePath := ""
	if r.Catalog == "" {
		filePath = r.Config.GetString("temporary_directory") + "/gutenberg-catalog.tar.zip"
		r.Catalog = filePath
	} else {
		filePath = r.Catalog
	}

	log.Printf("catalog is %v", filePath)

	// Try to housekeep and keep a local copy if we can't pull the new one
	oldFile := false
	_, fileNotFound := os.Stat(filePath)
	if fileNotFound == nil {
		// Don't constantly re-download in dev mode
		if r.Config.GetString("mode") == "development" {
			log.Printf("DEV: not re-pulling pg catalog, please remove file to re-download")
			return nil
		}

		renameErr := os.Rename(filePath, filePath+".bak")
		if renameErr == nil {
			oldFile = true
		}
	}

	// Create the file
	outFile, fileCreateErr := os.Create(filePath)
	if fileCreateErr != nil {
		if oldFile && os.Rename(filePath+".bak", filePath) != nil {
			log.Printf("error while trying to recover, could not restore old catalog")
		}
		return fmt.Errorf("issue creating temporary file:  %v with error: %v", filePath, fileCreateErr)
	}

	url := r.Config.GetString("gutenberg_feed_url") + "cache/epub/feeds/rdf-files.tar.zip"
	log.Printf("Pulling Gutenberg Catalog from url: %v", url)

	// Get the data
	resp, httpErr := http.Get(url)
	if httpErr != nil {
		if oldFile && os.Rename(filePath+".bak", filePath) != nil {
			log.Printf("error while trying to recover, could not restore old catalog")
		}
		return fmt.Errorf("issue pulling gutenberg catalog from url: %v with issue: %v", url, httpErr)
	}

	// Check server response
	if resp.StatusCode != http.StatusOK {
		if oldFile && os.Rename(filePath+".bak", filePath) != nil {
			log.Printf("error while trying to recover, could not restore old catalog")
		}
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, writeErr := io.Copy(outFile, resp.Body)
	if writeErr != nil {
		if oldFile && os.Rename(filePath+".bak", filePath) != nil {
			log.Printf("error while trying to recover, could not restore old catalog")
		}
		return fmt.Errorf("issue saving gutenberg catalog to working dir with issue: %v", writeErr)
	}

	resp.Body.Close()
	outFile.Close()

	log.Printf("Successfully pulled Gutenberg Catalog")

	if oldFile && os.Remove(filePath+".bak") != nil {
		log.Printf("Error: could not remove old catalog file: %v", filePath+".bak")
	}

	return nil
}

func (r *Runtime) getNumberOfRDFRecords() (int, error) {
	recordCount := 0

	zipReader, zipError := zip.OpenReader(r.Catalog)
	if zipError != nil {
		return -1, fmt.Errorf("error creating zip reader for catalog: %v", zipError)
	}

	if len(zipReader.File) != 1 {
		return -1, fmt.Errorf("unknown file format when opening catalog zip")
	}
	zippedTarFile, readZippedTarError := zipReader.File[0].Open()
	if readZippedTarError != nil {
		return -1, fmt.Errorf("error opening zipped tar reader for catalog: %v", readZippedTarError)
	}

	// Create a tar reader on top of the gzip reader
	tarReader := tar.NewReader(zippedTarFile)

	for {
		_, tarReadErr := tarReader.Next()
		if tarReadErr == io.EOF {
			break // End of archive reached
		} else if tarReadErr != nil {
			return -1, fmt.Errorf("error reading catalog tar archive: %v", tarReadErr)
		}
		recordCount++
	}
	return recordCount, nil
}

func (r *Runtime) getRDFByID(id int) (*rdf.RDF, error) {
	if id < 1 {
		return nil, fmt.Errorf("id out of range: %v", id)
	}

	idStr := strconv.Itoa(id)

	targetFileName := "cache/epub/" + idStr + "/pg" + idStr + ".rdf"

	zipReader, zipError := zip.OpenReader(r.Catalog)
	if zipError != nil {
		return nil, fmt.Errorf("error creating zip reader for catalog: %v", zipError)
	}

	if len(zipReader.File) != 1 {
		return nil, fmt.Errorf("unknown file format when opening catalog zip")
	}
	zippedTarFile, readZippedTarError := zipReader.File[0].Open()
	if readZippedTarError != nil {
		return nil, fmt.Errorf("error opening zipped tar reader for catalog: %v", readZippedTarError)
	}

	// Create a tar reader on top of the gzip reader
	tarReader := tar.NewReader(zippedTarFile)

	var rdfFile []byte

	for {
		header, tarReadErr := tarReader.Next()
		if tarReadErr == io.EOF {
			break // End of archive reached
		} else if tarReadErr != nil {
			return nil, fmt.Errorf("error reading catalog tar archive: %v", tarReadErr)
		}

		// Check if you've found the desired file
		if header.Name == targetFileName {
			// Open the file within the archive for reading
			var tarReadTargetFileErr error
			rdfFile, tarReadTargetFileErr = io.ReadAll(tarReader)
			if tarReadTargetFileErr != nil {
				return nil, fmt.Errorf("error reading target file from catalog: %v", tarReadTargetFileErr)
			}

			break // Stop after reading the desired file
		}
	}
	zipReader.Close()

	var rdf rdf.RDF
	xmlParseErr := xml.Unmarshal(rdfFile, &rdf)

	if xmlParseErr != nil {
		return nil, fmt.Errorf("error parsing rdf xml into struct: %v", xmlParseErr)
	}

	return &rdf, nil
}
