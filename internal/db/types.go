package db 

import (
	"encoding/xml"
)

type DB struct {
	version int `json: version`
	
}


func (b *BitField) UnmarshalJSON(data []byte) error {
    // ...
    return nil
}
