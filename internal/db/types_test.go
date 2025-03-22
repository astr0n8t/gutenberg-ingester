package db

import (
	"encoding/json"
	"fmt"
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
