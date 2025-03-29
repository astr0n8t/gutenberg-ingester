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

	"github.com/astr0n8t/gutenberg-ingester/pkg/db"
	"github.com/astr0n8t/gutenberg-ingester/pkg/rdf"
)

func fullSync(config ConfigStore, db *db.DB) error {
	return nil
}

func pullCatalog(config ConfigStore) (string, error) {
	filePath := config.GetString("temporary_directory") + "/gutenberg-catalog.tar.zip"
	// Create the file
	outFile, fileCreateErr := os.Create(filePath)
	if fileCreateErr != nil {
		return "", fmt.Errorf("issue creating temporary file:  %v with error: %v", filePath, fileCreateErr)
	}

	url := config.GetString("gutenberg_feed_url") + "feeds/rdf-files.tar.zip"
	log.Printf("Pulling Gutenberg Catalog from url: %v", url)

	// Get the data
	resp, httpErr := http.Get(url)
	if httpErr != nil {
		return "", fmt.Errorf("issue pulling gutenberg catalog from url: %v with issue: %v", url, httpErr)
	}

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, writeErr := io.Copy(outFile, resp.Body)
	if writeErr != nil {
		return "", fmt.Errorf("issue saving gutenberg catalog to working dir with issue: %v", writeErr)
	}

	resp.Body.Close()
	outFile.Close()

	log.Printf("Successfully pulled Gutenberg Catalog")

	return filePath, nil
}

func getNumberOfRDFRecords(catalogFilePath string) (int, error) {
	recordCount := 0

	zipReader, zipError := zip.OpenReader(catalogFilePath)
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
		header, tarReadErr := tarReader.Next()
		if tarReadErr == io.EOF {
			break // End of archive reached
		} else if tarReadErr != nil {
			return -1, fmt.Errorf("error reading catalog tar archive: %v", tarReadErr)
		}
		log.Printf("%v", header.Name)
		recordCount++
	}
	return recordCount, nil
}

func getRDFByID(id int, catalogFilePath string) (*rdf.RDF, error) {
	if id < 1 {
		return nil, fmt.Errorf("id out of range: %v", id)
	}

	idStr := strconv.Itoa(id)

	targetFileName := "cache/epub/" + idStr + "/pg" + idStr + ".rdf"

	zipReader, zipError := zip.OpenReader(catalogFilePath)
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
