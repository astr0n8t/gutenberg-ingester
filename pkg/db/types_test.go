package db

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestMarshalDB(t *testing.T) {
	db := NewDB()

	jsonData1, err1 := json.Marshal(db)
	if err1 != nil {
		t.Errorf("Unable to marshal empty DB to json %v", err1)
	}
	fmt.Printf("Empty DB in json: %v\n", string(jsonData1))
}

func TestUnMarshalDB(t *testing.T) {
	s1 := []byte(`{"version":1,"download_history":{"history":"H4sIAAAAAAAA/+zAgQAAAADCsPypEzjCNgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACA8AAP//ZgMpkYA4AQA="}}`)
	var db DB

	json.Unmarshal(s1, &db)
	if db.Version != 1 {
		t.Errorf("Unable to unmarshal empty DB from json")
	}
}

func TestWriteDB(t *testing.T) {
	db := NewDB()

	testFile := "/tmp/gutenberg_test1_db.json"
	err := db.WriteDBToFile(testFile)
	if err != nil {
		t.Errorf("Unable to write empty DB to json file %v", err)
	}

	os.Remove(testFile)
}

func TestOpenDB(t *testing.T) {
	testFile := "/tmp/gutenberg_test2_db.json"
	db1, err1 := OpenDBFromFile("/tmp/gutenberg_test_db.json")

	if err1 != nil {
		t.Errorf("unable to write and read empty DB %v", err1)
	}

	if db1.GetDownloaded(1) {
		t.Errorf("created and read in DB file is not empty")
	}

	db2, err2 := OpenDBFromFile("/tmp/gutenberg_test_db.json")

	if err2 != nil {
		t.Errorf("unable to read empty DB %v", err2)
	}

	if db2.GetDownloaded(1) {
		t.Errorf("read in DB file is not empty")
	}
	os.Remove(testFile)
}
