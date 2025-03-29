package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/astr0n8t/gutenberg-ingester/pkg/db"
)

func fullSync(config ConfigStore, db *db.Db) err {
}

func pullCatalog(config ConfigStore) string {
	filePath := config.GetString("temporary_directory") + "/gutenberg-catalog.xml.zip"
	// Create the file
	out, fileCreateErr := os.Create(filepath)
	if fileCreateErr != nil {
		return fmt.Errorf("issue creating temporary file:  %v with error: %v", filepath, fileCreateErr)
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	out.Close()

	return nil
}
