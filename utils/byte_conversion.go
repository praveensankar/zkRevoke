package utils

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// GetBytes converts an arbitrary interface to a byte array.
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func BytesToString(input []byte) string {
	res := fmt.Sprintf("%s", input)
	return res
}
