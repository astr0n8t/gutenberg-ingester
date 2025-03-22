package history 

import (
	"fmt"
	"testing"
	"encoding/json"
)

func TestSetAndGetIndex(t *testing.T) {
	var h History

	fmt.Printf("Value at 20: %v\n", h.getHistory(20))

	for i := 0; i < 9; i++ {
		fmt.Printf("testing at index %v\n", i)
		h.setHistory(i)
	}


	h.setHistory(2137)
	h.setHistory(21)
	h.setHistory(4096)

	fmt.Printf("Value at 2137: %v\n", h.getHistory(2137))
	fmt.Printf("Value at 21: %v\n", h.getHistory(21))
	fmt.Printf("Value at 4096: %v\n", h.getHistory(4096))
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
		h2.setHistory(i)
	}
	for i := 300; i < 423; i++ {
		h2.setHistory(i)
	}

	jsonData2, err2 := json.Marshal(h2) 
	if err2 != nil {
		t.Errorf("Unable to marshal partially set history to json %v", err2)
	}

	fmt.Printf("Partially set history in json: %v\n", string(jsonData2))


	var h3 History

	for i := 0; i < 20000; i++ {
		h3.setHistory(i)
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
	if h1.getHistory(0) == true {
		t.Errorf("Unable to unmarshal empty history to json")
	}
	fmt.Printf("Length of unmarshalled history for empty history: %v\n", len(h1.bitmap))

	s2 := []byte(`{"history":"H4sIAAAAAAAA//qPCtgZsIEPqIrqAQEAAP//Ftid7jUAAAA="}`)
	var h2 History

	json.Unmarshal(s2, &h2)
	if h2.getHistory(122) == true && h2.getHistory(300) == true {
		t.Errorf("Unable to unmarshal partially set history to json")
	}
	fmt.Printf("Length of unmarshalled history for partially set history: %v\n", len(h2.bitmap))

	s3 := []byte(`{"history":"H4sIAAAAAAAA/+zAAQ0AAADCoP6pNcc3GAAAOQ8AAP//MwLNoMQJAAA="}`)
	var h3 History

	json.Unmarshal(s3, &h3)
	if h1.getHistory(1999) == true {
		t.Errorf("Unable to unmarshal partially set history to json")
	}
	fmt.Printf("Length of unmarshalled history for fully set history: %v\n", len(h3.bitmap))
}
