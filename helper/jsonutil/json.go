package jsonutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

func EncodeJSON(in interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(in); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeJSON(data []byte, out interface{}) error {
	// Decoding requires a pointer type to be supplied
	value := reflect.ValueOf(out)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("decoding the value into an invalid type: %v", reflect.TypeOf(out))
	}

	return DecodeJSONFromReader(bytes.NewReader(data), out)
}

func DecodeJSONFromReader(r io.Reader, out interface{}) error {
	dec := json.NewDecoder(r)

	// While decoding JSON values, intepret the integer values as numbers instead of floats.
	dec.UseNumber()

	// Since 'out' is an interface representing a pointer, pass it to the decoder without an '&'
	return dec.Decode(out)
}
