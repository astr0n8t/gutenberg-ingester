package history 

import (
    "log"
    "fmt"
    "bytes"
    "io/ioutil"
    "encoding/json"
    "encoding/base64"
    "encoding/binary"
    "compress/gzip"
)

// Function to set the index
// i: index to set 
func (h *History) setHistory(i int) {
    offset := i / 8

    // If the array isn't big enough, we need to make it big enough to hold the index
    if len(h.bitmap) < (offset + 1) {
        numToAdd := (offset - len(h.bitmap)) + 1
        for i := 0; i < numToAdd; i++ {
            h.bitmap = append(h.bitmap, 0)
        }
    }

    h.bitmap[offset] = setIntAtIndex(i % 8, h.bitmap[offset])
}

// Function to check whether the index is set
// i: index to check
// returns whether the index is set or not
func (h *History) getHistory(i int) bool {
    offset := i / 8

    // If the array isn't big enough, it was never set
    if len(h.bitmap) < offset || len(h.bitmap) == 0 {
        return false
    }

    return getIntAtIndex(i % 8, h.bitmap[offset])
}

// Function to easily get a specific bit of the int
// i: index of int to access
// d: int to access
// returns whether the bit is set or not
func getIntAtIndex(i int, d uint8) bool {
    if i > 7 || i < 0 {
		log.Printf("getting wrong index for uint8 in DB, should never see this: %v", i)
        return false
    }
    if ((d & ( 1 << i )) >> i) == 1 {
        return true
    }
    return false
}

// Function to easily set a specific bit of the int
// i: index of int to access
// d: int to access
// returns the updated uint8
func setIntAtIndex(i int, d uint8) uint8 {
    if i > 7 || i < 0 {
		log.Printf("getting wrong index for uint8 in DB, should never see this: %v", i)
        return d 
    }
    return d + (1 << i)
}

func (h *History) UnmarshalJSON(data []byte) error {
    var jsonData map[string]interface{}
    if unmarhsalErr := json.Unmarshal(data, &jsonData); unmarhsalErr != nil {
		return fmt.Errorf("failed to unmarshal DB download history json: %v", unmarhsalErr)
    }
    historyB64, dataErr := jsonData["history"].(string)
    if dataErr != true {
		return fmt.Errorf("failed to retrieve DB download history json: %v", dataErr)
    }

    historyCompressedData, b64err := base64.StdEncoding.DecodeString(historyB64)
    if b64err != nil {
		return fmt.Errorf("failed to decode DB download history base64: %v", b64err)
    }

    bytesReader := bytes.NewReader(historyCompressedData)
    gzipReader, gzipErr := gzip.NewReader(bytesReader)
    if gzipErr != nil {
		return fmt.Errorf("failed to decompress DB download history: %v", gzipErr)
    }
    defer gzipReader.Close()

    historyData, bytesReadErr := ioutil.ReadAll(gzipReader)
    if bytesReadErr != nil {
		return fmt.Errorf("failed to read decompressed DB download history: %v", bytesReadErr)
    }

    h.bitmap = historyData

    return nil
}

func (h History) MarshalJSON() ([]byte, error) {
   // historyAsBytes := make([]byte, len(h.bitmap))

   // for i := 0; i < len(h.bitmap); i++ {
   //     binary.PutUvarint(historyAsBytes[i+0:i+4], h.bitmap[i]) 
   // }

    // Convert to byte slice in big-endian order
    var historyBuffer bytes.Buffer
    binary.Write(&historyBuffer, binary.BigEndian, h.bitmap) 

    // get a gzip writer ready
    var buf bytes.Buffer
    writer := gzip.NewWriter(&buf)

    // compress the data
    _, compressErr := writer.Write(historyBuffer.Bytes())
    if compressErr != nil {
		return nil, fmt.Errorf("failed to compress DB download history: %v", compressErr)
    }

    writer.Close()

    // base64 encode and serialize it
    serializedHistory := base64.StdEncoding.EncodeToString(buf.Bytes())
    return json.Marshal(map[string]interface{}{
        "history": serializedHistory,
    })
}

