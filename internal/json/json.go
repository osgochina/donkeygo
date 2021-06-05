package json

import (
	"bytes"
	"encoding/json"
	"io"
)

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func UnmarshalUseNumber(data []byte, v interface{}) error {
	decoder := NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	return decoder.Decode(v)
}

func NewEncoder(writer io.Writer) *json.Encoder {
	return json.NewEncoder(writer)
}

func NewDecoder(reader io.Reader) *json.Decoder {
	return json.NewDecoder(reader)
}

func Valid(data []byte) bool {
	return json.Valid(data)
}
