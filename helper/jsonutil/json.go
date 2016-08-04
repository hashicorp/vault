package jsonutil

import (
	"bytes"
	"compress/lzw"
	"encoding/json"
	"fmt"
	"io"
)

const (
	canaryByte byte = 'Z'
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

	// Read the first byte out
	bytesReader := bytes.NewReader(data)
	firstByte, err := bytesReader.ReadByte()
	if err != nil {
		return fmt.Errorf("failed to find the canary in the compressed input")
	}

	// If the first byte doesn't match the canaryByte, it means that the
	// content was not compressed in the first place. Try JSON decoding it.
	if canaryByte != firstByte {
		return DecodeJSON(data, out)
	} else {
		// If the first byte matches the canaryByte, remove the canary
		// byte and try to decompress the data before JSON decoding it.
		data = data[1:]
	}

	// Create a reader to read out the compressed data
	reader := lzw.NewReader(bytes.NewReader(data), lzw.LSB, 8)

	// Close the io.ReadCloser
	defer reader.Close()

	// Read all the compressed data out into a buffer
	var jsonBuf bytes.Buffer
	if _, err := io.Copy(&jsonBuf, reader); err != nil {
		return err
	}

	// JSON decode the read out bytes
	return DecodeJSON(jsonBuf.Bytes(), out)
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

	// Create a buffer and place the canary as its first byte
	var buf bytes.Buffer
	buf.Write([]byte{canaryByte})

	// Create writer to compress the JSON encoded bytes
	writer := lzw.NewWriter(&buf, lzw.LSB, 8)

	// Compress the JSON bytes
	if _, err := writer.Write(encodedBytes); err != nil {
		return nil, fmt.Errorf("failed to compress JSON string; err: %v", err)
	}

	// Close the io.WriteCloser
	if err := writer.Close(); err != nil {
		return nil, err
	}

	// Return the compressed bytes with canary byte at the start
	return buf.Bytes(), nil
}
