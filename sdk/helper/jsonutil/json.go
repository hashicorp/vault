// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jsonutil

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/helper/compressutil"
)

const CustomMaxJSONDepth = 500

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

	jsonBytes, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read JSON input into buffer: %w", err)
	}
	if len(jsonBytes) == 0 {
		return nil
	}

	// We need to read the content of the JSON before decoding it so that we determine the depth of the JSON and fail fast in case it exceeds our internal depth limit.
	// This done because of the DoS found in https://hashicorp.atlassian.net/browse/VAULT-36788, which affects audit logging and other operations such as HMACing of fields.
	depthReader := bytes.NewReader(jsonBytes)
	actualDepth, err := calculateMaxDepthStreaming(depthReader)
	if err != nil {
		return fmt.Errorf("failed to scan JSON for depth: %w", err)
	}
	if actualDepth > CustomMaxJSONDepth {
		return fmt.Errorf("JSON input exceeds allowed nesting depth (%d > %d)", actualDepth, CustomMaxJSONDepth)
	}

	dec := bytes.NewReader(jsonBytes)
	jsonDec := json.NewDecoder(dec)

	// While decoding JSON values, interpret the integer values as `json.Number`s instead of `float64`.
	jsonDec.UseNumber()

	// Since 'out' is an interface representing a pointer, pass it to the decoder without an '&'
	return jsonDec.Decode(out)
}

// This function avoids building the full map in memory, suitable for very deep JSON.
func calculateMaxDepthStreaming(jsonReader io.Reader) (int, error) {
	decoder := json.NewDecoder(jsonReader)
	decoder.UseNumber()
	maxDepth := 0
	currentDepth := 0

	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("error reading JSON token: %w", err)
		}

		if delim, ok := t.(json.Delim); ok {
			if delim == '{' || delim == '[' {
				currentDepth++
				if currentDepth > maxDepth {
					maxDepth = currentDepth
				}
			} else if delim == '}' || delim == ']' {
				currentDepth--
			}
		}
	}
	// Add this check to account for unmatched delimiters
	if currentDepth != 0 {
		return 0, fmt.Errorf("malformed JSON, unmatched delimiters")
	}
	return maxDepth, nil
}
