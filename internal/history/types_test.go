package history

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"testing"
)

func TestSetAndGetIndex(t *testing.T) {
	var h History

	fmt.Printf("Value at 20: %v\n", h.GetHistory(20))

	for i := 0; i < 9; i++ {
		h.SetHistory(i)
	}

	for i := 0; i < 9; i++ {
		if !h.GetHistory(i) {
			t.Errorf("value not set or get correctly")
		}
	}

	h.SetHistory(2137)
	h.SetHistory(21)
	h.SetHistory(4096)

	fmt.Printf("Value at 2137: %v\n", h.GetHistory(2137))
	fmt.Printf("Value at 21: %v\n", h.GetHistory(21))
	fmt.Printf("Value at 4096: %v\n", h.GetHistory(4096))

	if !h.GetHistory(2137) || !h.GetHistory(21) || !h.GetHistory(4096) {
		t.Errorf("value not set or get correctly")
	}

	h.UnsetHistory(2137)
	fmt.Printf("Value at 2137: %v\n", h.GetHistory(2137))
	if h.GetHistory(2137) == true {
		t.Errorf("value not unset correctly")
	}
}

func TestMarshalHistory(t *testing.T) {
	var h1 History

	jsonData1, err1 := json.Marshal(h1)
	if err1 != nil {
		t.Errorf("Unable to marshal empty history to json %v", err1)
	}

	fmt.Printf("Empty history in json: %v\n", string(jsonData1))

	var h2 History

	for i := 0; i < 123; i++ {
		h2.SetHistory(i)
	}
	for i := 300; i < 423; i++ {
		h2.SetHistory(i)
	}

	h2.UnsetHistory(50)

	jsonData2, err2 := json.Marshal(h2)
	if err2 != nil {
		t.Errorf("Unable to marshal partially set history to json %v", err2)
	}

	fmt.Printf("Partially set history in json: %v\n", string(jsonData2))

	var h3 History

	for i := 0; i < 20000; i++ {
		h3.SetHistory(i)
	}

	jsonData3, err3 := json.Marshal(h3)
	if err3 != nil {
		t.Errorf("Unable to marshal fully set history to json %v", err3)
	}

	fmt.Printf("Fully set history in json: %v\n", string(jsonData3))
}

func TestUnMarshalHistory(t *testing.T) {
	s1 := []byte(`{"history":"H4sIAAAAAAAA/wEAAP//AAAAAAAAAAA="}`)
	var h1 History

	json.Unmarshal(s1, &h1)
	if h1.GetHistory(0) == true {
		t.Errorf("Unable to unmarshal empty history to json")
	}

	hash1 := sha256.New()
	hash1.Write(h1.bitmap)
	hashSum1 := hash1.Sum(nil)
	hexHash1 := fmt.Sprintf("%x", hashSum1)

	if hexHash1 != "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
		fmt.Printf("SHA256 Hash of Empty History: %s\n", hexHash1)
		t.Errorf("Hash mismatch on unmarshalled array for empty history")
	}

	fmt.Printf("Length of unmarshalled history for empty history: %v\n", len(h1.bitmap))

	s2 := []byte(`{"history":"H4sIAAAAAAAA//oPBr//QwE7Azbw4T8KqAcEAAD//8PzCJA1AAAA"}`)
	var h2 History

	json.Unmarshal(s2, &h2)
	if h2.GetHistory(50) == true || h2.GetHistory(122) == false || h2.GetHistory(300) == false || h2.GetHistory(200) == true {
		t.Errorf("Unable to unmarshal partially set history to json")
	}

	hash2 := sha256.New()
	hash2.Write(h2.bitmap)
	hashSum2 := hash2.Sum(nil)
	hexHash2 := fmt.Sprintf("%x", hashSum2)
	if hexHash2 != "565d0cf4b14aadb6ec6ac63d8b645b81c3975491f7aa846b312b1b130ee74b10" {
		fmt.Printf("SHA256 Hash of Partially Set History: %s\n", hexHash2)
		t.Errorf("Hash mismatch on unmarshalled array for partially set history")
	}

	fmt.Printf("Length of unmarshalled history for partially set history: %v\n", len(h2.bitmap))

	s3 := []byte(`{"history":"H4sIAAAAAAAA/+zAAQ0AAADCoP6pNcc3GAAAOQ8AAP//MwLNoMQJAAA="}`)
	var h3 History

	json.Unmarshal(s3, &h3)
	if h3.GetHistory(19999) == false {
		t.Errorf("Unable to unmarshal partially set history to json")
	}

	hash3 := sha256.New()
	hash3.Write(h3.bitmap)
	hashSum3 := hash3.Sum(nil)
	hexHash3 := fmt.Sprintf("%x", hashSum3)
	if hexHash3 != "7bf918ed22ac03fcccb88561850fa0d8196d89a5b42334364d972d9606f404f4" {
		fmt.Printf("SHA256 Hash of Fully Set History: %s\n", hexHash3)
		t.Errorf("Hash mismatch on unmarshalled array for fully set history")
	}

	fmt.Printf("Length of unmarshalled history for fully set history: %v\n", len(h3.bitmap))
}
