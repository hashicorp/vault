package jsonutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/hashicorp/vault/helper/compressutil"
)

// Encodes/Marshals the given object into JSON
func EncodeJSON(in interface{}) ([]byte, error) {
	if in == nil {
		return nil, fmt.Errorf("input for encoding is nil")
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(in); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decodes/Unmarshals the given JSON into a desired object
func DecodeJSON(data []byte, out interface{}) error {
	if data == nil {
		return fmt.Errorf("'data' being decoded is nil")
	}
	if out == nil {
		return fmt.Errorf("output parameter 'out' is nil")
	}

	return DecodeJSONFromReader(bytes.NewReader(data), out)
}

// Decodes/Unmarshals the given io.Reader pointing to a JSON, into a desired object
func DecodeJSONFromReader(r io.Reader, out interface{}) error {
	if r == nil {
		return fmt.Errorf("'io.Reader' being decoded is nil")
	}
	if out == nil {
		return fmt.Errorf("output parameter 'out' is nil")
	}

	dec := json.NewDecoder(r)

	// While decoding JSON values, intepret the integer values as `json.Number`s instead of `float64`.
	dec.UseNumber()

	// Since 'out' is an interface representing a pointer, pass it to the decoder without an '&'
	return dec.Decode(out)
}

// DecompressAndDecodeJSON checks if the first byte in the input matches the
// canary byte. If it does, the input will be decompressed (lzw) before being
// JSON decoded. If the does not, the input will be JSON decoded without
// attempting to decompress it.
func DecompressAndDecodeJSON(data []byte, out interface{}) error {
	if data == nil || len(data) < 2 {
		return fmt.Errorf("'data' being decoded is invalid")
	}
	if out == nil {
		return fmt.Errorf("output parameter 'out' is nil")
	}

	decompressedBytes, unencrypted, err := compressutil.Decompress(data, compressutil.CompressionCanaryJSON)
	if err != nil {
		return fmt.Errorf("failed to decompress JSON: err: %v", err)
	}

	// If the data supplied failed to contain the JSON compression canary,
	// it can be inferred that it was not compressed in the first place.
	// Try to JSON decode it.
	if unencrypted {
		return DecodeJSON(data, out)
	}

	if decompressedBytes == nil || len(decompressedBytes) == 0 {
		return fmt.Errorf("decompressed data being decoded is invalid")
	}

	// JSON decode the read out bytes
	return DecodeJSON(decompressedBytes, out)
}

// EncodeJSONAndCompress encodes the given input into JSON and compresses the
// encoded value (lzw). A canary byte is placed at the beginning of the
// returned bytes for the logic in decompression method to identify compressed
// input.
func EncodeJSONAndCompress(in interface{}) ([]byte, error) {
	if in == nil {
		return nil, fmt.Errorf("input for encoding is nil")
	}

	// First JSON encode the given input
	encodedBytes, err := EncodeJSON(in)
	if err != nil {
		return nil, err
	}
	log.Printf("EncodeJSONAndCompress: len(encodedBytes): %d\n", len(encodedBytes))

	return compressutil.Compress(encodedBytes, compressutil.CompressionCanaryJSON)
}
