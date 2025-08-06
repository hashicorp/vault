// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jsonutil

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/helper/compressutil"
)

// EncodeJSON encodes/marshals the given object into JSON
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

// EncodeJSONAndCompress encodes the given input into JSON and compresses the
// encoded value (using Gzip format BestCompression level, by default). A
// canary byte is placed at the beginning of the returned bytes for the logic
// in decompression method to identify compressed input.
func EncodeJSONAndCompress(in interface{}, config *compressutil.CompressionConfig) ([]byte, error) {
	if in == nil {
		return nil, fmt.Errorf("input for encoding is nil")
	}

	// First JSON encode the given input
	encodedBytes, err := EncodeJSON(in)
	if err != nil {
		return nil, err
	}

	if config == nil {
		config = &compressutil.CompressionConfig{
			Type:                 compressutil.CompressionTypeGzip,
			GzipCompressionLevel: gzip.BestCompression,
		}
	}

	return compressutil.Compress(encodedBytes, config)
}

// DecodeJSON tries to decompress the given data. The call to decompress, fails
// if the content was not compressed in the first place, which is identified by
// a canary byte before the compressed data. If the data is not compressed, it
// is JSON decoded directly. Otherwise the decompressed data will be JSON
// decoded.
func DecodeJSON(data []byte, out interface{}) error {
	if data == nil || len(data) == 0 {
		return fmt.Errorf("'data' being decoded is nil")
	}
	if out == nil {
		return fmt.Errorf("output parameter 'out' is nil")
	}

	// Decompress the data if it was compressed in the first place
	decompressedBytes, uncompressed, err := compressutil.Decompress(data)
	if err != nil {
		return errwrap.Wrapf("failed to decompress JSON: {{err}}", err)
	}
	if !uncompressed && (decompressedBytes == nil || len(decompressedBytes) == 0) {
		return fmt.Errorf("decompressed data being decoded is invalid")
	}

	// If the input supplied failed to contain the compression canary, it
	// will be notified by the compression utility. Decode the decompressed
	// input.
	if !uncompressed {
		data = decompressedBytes
	}

	return DecodeJSONFromReader(bytes.NewReader(data), out)
}

// DecodeJSONFromReader Decodes/Unmarshals the given io.Reader pointing to a JSON, into a desired object
func DecodeJSONFromReader(r io.Reader, out interface{}) error {
	if r == nil {
		return fmt.Errorf("'io.Reader' being decoded is nil")
	}
	if out == nil {
		return fmt.Errorf("output parameter 'out' is nil")
	}

	dec := json.NewDecoder(r)

	// While decoding JSON values, interpret the integer values as `json.Number`s instead of `float64`.
	dec.UseNumber()

	// Since 'out' is an interface representing a pointer, pass it to the decoder without an '&'
	return dec.Decode(out)
}

// containerState holds information about an open JSON container (object or array).
type containerState struct {
	Type  json.Delim // '{' or '['
	Count int        // Number of entries (for objects) or elements for arrays)
}

// JSONLimits defines the configurable limits for JSON validation.
type JSONLimits struct {
	MaxDepth             int
	MaxStringValueLength int
	MaxObjectEntryCount  int
	MaxArrayElementCount int
}

// isWhitespace checks if a byte is a JSON whitespace character.
func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

// VerifyMaxDepthStreaming scans the JSON stream to determine its maximum nesting depth
// and enforce various limits. It first checks if the stream is likely JSON before proceeding.
func VerifyMaxDepthStreaming(jsonReader io.Reader, limits JSONLimits) (int, error) {
	// Use a buffered reader to peek at the stream without consuming it from the original reader.
	bufReader := bufio.NewReader(jsonReader)

	// Find the first non-whitespace character.
	var firstByte byte
	var err error
	for {
		firstByte, err = bufReader.ReadByte()
		if err != nil {
			// If we hit EOF before finding a real character, it's an empty or whitespace-only payload.
			if err == io.EOF {
				return 0, nil
			}
			return 0, err // A different I/O error occurred.
		}
		if !isWhitespace(firstByte) {
			break // Found the first significant character.
		}
	}

	// If the payload doesn't start with '{' or '[', assume it's not a JSON object or array
	// and that our limits do not apply.
	if firstByte != '{' && firstByte != '[' {
		return 0, nil
	}

	fullStreamReader := io.MultiReader(bytes.NewReader([]byte{firstByte}), bufReader)
	decoder := json.NewDecoder(fullStreamReader)
	decoder.UseNumber()

	var (
		maxDepth      = 0
		currentDepth  = 0
		isKeyExpected bool
	)
	containerInfoStack := make([]containerState, 0, limits.MaxDepth)

	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			// Any error from the decoder is now considered a real error.
			return 0, fmt.Errorf("error reading JSON token: %w", err)
		}

		switch v := t.(type) {
		case json.Delim:
			switch v {
			case '{', '[':
				currentDepth++
				// Check against the limit directly.
				if currentDepth > limits.MaxDepth {
					return 0, fmt.Errorf("JSON input exceeds allowed nesting depth")
				}
				if currentDepth > maxDepth {
					maxDepth = currentDepth
				}

				containerInfoStack = append(containerInfoStack, containerState{Type: v, Count: 0})
				if v == '{' {
					isKeyExpected = true
				}
			case '}', ']':
				if len(containerInfoStack) == 0 {
					return 0, fmt.Errorf("malformed JSON: unmatched closing delimiter '%c'", v)
				}
				top := containerInfoStack[len(containerInfoStack)-1]
				containerInfoStack = containerInfoStack[:len(containerInfoStack)-1]
				currentDepth--
				if (v == '}' && top.Type != '{') || (v == ']' && top.Type != '[') {
					return 0, fmt.Errorf("malformed JSON: mismatched closing delimiter '%c' for opening '%c'", v, top.Type)
				}
				if len(containerInfoStack) > 0 && containerInfoStack[len(containerInfoStack)-1].Type == '{' {
					isKeyExpected = false
				}
			}
		case string:
			if len(v) > limits.MaxStringValueLength {
				return 0, fmt.Errorf("JSON string value exceeds allowed length")
			}
			if len(containerInfoStack) > 0 {
				top := &containerInfoStack[len(containerInfoStack)-1]
				if top.Type == '{' {
					if isKeyExpected {
						top.Count++
						if top.Count > limits.MaxObjectEntryCount {
							return 0, fmt.Errorf("JSON object exceeds allowed entry count")
						}
						isKeyExpected = false
					}
				} else if top.Type == '[' {
					top.Count++
					if top.Count > limits.MaxArrayElementCount {
						return 0, fmt.Errorf("JSON array exceeds allowed element count")
					}
				}
			}
		default: // Handles numbers, booleans, and nulls
			if len(containerInfoStack) > 0 {
				top := &containerInfoStack[len(containerInfoStack)-1]
				if top.Type == '[' {
					top.Count++
					if top.Count > limits.MaxArrayElementCount {
						return 0, fmt.Errorf("JSON array exceeds allowed element count")
					}
				} else if top.Type == '{' {
					isKeyExpected = true
				}
			}
		}
	}

	if len(containerInfoStack) != 0 {
		return 0, fmt.Errorf("malformed JSON, unclosed containers")
	}

	return maxDepth, nil
}
