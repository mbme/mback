package utils

import (
	"encoding/json"
)

func Encode(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
