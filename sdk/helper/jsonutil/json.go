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
	"strconv"
	"strings"
	"unicode"

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
	// '{' or '['
	Type json.Delim

	// Number of entries (for objects) or elements for arrays)
	Count int

	// isKey is true if the next expected token in an object is a key.
	isKey bool

	// keys tracks the keys seen in an object to detect duplicates.
	// It is only initialized for objects ('{').
	keys map[string]struct{}
}

// JSONLimits defines the configurable limits for JSON validation.
type JSONLimits struct {
	MaxDepth             int
	MaxStringValueLength int
	MaxObjectEntryCount  int
	MaxArrayElementCount int
	MaxTokens            int
}

// isWhitespace checks if a byte is a JSON whitespace character.
func isWhitespace(b byte) bool {
	// Standard JSON whitespace characters (RFC 8259)
	if b == ' ' || b == '\t' || b == '\n' || b == '\r' {
		return true
	}
	// Custom support for non-standard Unit Separator character (ASCII 31)
	if b == 31 || b == 139 {
		return true
	}
	return false
}

// defaultBufferSize is the default size for the bufio.Reader buffer
const defaultBufferSize = int64(4096)

// VerifyMaxDepthStreaming scans the JSON stream to enforce nesting depth, counts,
// and other limits without decoding the full structure into memory.
func VerifyMaxDepthStreaming(jsonReader io.Reader, limits JSONLimits, maxRequestSize *int64) (int, error) {
	// If the default buffer size is larger than the max request size, use the max request size
	// for the buffer to avoid over-reading.
	bufferSize := defaultBufferSize
	if maxRequestSize != nil && *maxRequestSize < defaultBufferSize {
		bufferSize = *maxRequestSize
	}
	// Use a buffered reader to peek at the stream
	bufReader := bufio.NewReaderSize(jsonReader, int(bufferSize))

	bom, err := bufReader.Peek(3)
	if err == nil && bytes.Equal(bom, []byte{0xEF, 0xBB, 0xBF}) {
		_, _ = bufReader.Discard(3)
	}

	// We use a manual token loop instead of json.Decoder to gain low-level
	// control over the stream. This is necessary to fix a vulnerability where
	// the raw byte length of strings with escape sequences was not correctly limited.
	var (
		maxDepth          int
		currentDepth      int
		tokenCount        int
		lastTokenWasComma bool
	)
	containerInfoStack := make([]containerState, 0, limits.MaxDepth)

	// Prime the loop by finding the first non-whitespace character.
	if err := skipWhitespace(bufReader); err != nil {
		if err == io.EOF {
			// An empty payload or one with only whitespace is valid. Skip verification.
			return 0, nil
		}
		return 0, err
	}

	b, err := bufReader.Peek(1)
	if err != nil {
		// This can happen if there's an I/O error after skipping whitespace.
		return 0, err
	}

	// If the payload doesn't start with a JSON container ('{' or '['), skip
	// verification. The limits are intended for structured data, not primitives
	// or other formats.
	if b[0] != '{' && b[0] != '[' {
		return 0, nil
	}

	for {
		// Check for EOF before peeking.
		b, err := bufReader.Peek(1)
		// Any error from the decoder is now considered a real error.
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, fmt.Errorf("error reading JSON token: %w", err)
		}

		// If the last token was a comma, the next token cannot be a closing delimiter.
		if lastTokenWasComma && (b[0] == '}' || b[0] == ']') {
			if b[0] == '}' {
				return 0, fmt.Errorf("invalid character '}' after object key-value pair")
			}
			return 0, fmt.Errorf("invalid character ']' after array element")
		}

		// After a top-level value, any other character is an error.
		if len(containerInfoStack) == 0 && maxDepth > 0 {
			return 0, fmt.Errorf("invalid character '%c' after top-level value", b[0])
		}

		// Increment and check the total token count limit.
		if limits.MaxTokens > 0 {
			tokenCount++
			if tokenCount > limits.MaxTokens {
				return 0, fmt.Errorf("JSON payload exceeds allowed token count")
			}
		}

		var currentContainer *containerState
		if len(containerInfoStack) > 0 {
			currentContainer = &containerInfoStack[len(containerInfoStack)-1]
		}

		// Before processing the token, reset the comma flag. It will be set
		// again below if the current token is a comma.
		lastTokenWasComma = false

		switch b[0] {
		case '{', '[':
			delim, _ := bufReader.ReadByte()
			if currentContainer != nil {
				if currentContainer.Type == '[' {
					currentContainer.Count++
					if currentContainer.Count > limits.MaxArrayElementCount {
						return 0, fmt.Errorf("JSON array exceeds allowed element count")
					}
				} else {
					currentContainer.isKey = true
				}
			}

			// Handle depth checks and tracking together for clarity.
			currentDepth++
			if currentDepth > maxDepth {
				maxDepth = currentDepth
			}
			// Check depth limit immediately after incrementing.
			if limits.MaxDepth > 0 && currentDepth > limits.MaxDepth {
				return 0, fmt.Errorf("JSON input exceeds allowed nesting depth")
			}

			// For objects, initialize a map to track keys and prevent duplicates.
			var keys map[string]struct{}
			if delim == '{' {
				keys = make(map[string]struct{})
			}

			containerInfoStack = append(containerInfoStack, containerState{Type: json.Delim(delim), isKey: delim == '{', keys: keys})

		case '}', ']':
			// A closing brace cannot follow a colon without a value.
			if currentContainer != nil && currentContainer.Type == '{' && !currentContainer.isKey {
				return 0, fmt.Errorf("invalid character '}' after object key")
			}
			delim, _ := bufReader.ReadByte()
			if currentContainer == nil {
				return 0, fmt.Errorf("malformed JSON: unmatched closing delimiter '%c'", delim)
			}
			if (delim == '}' && currentContainer.Type != '{') || (delim == ']' && currentContainer.Type != '[') {
				return 0, fmt.Errorf("malformed JSON: mismatched closing delimiter '%c'", delim)
			}
			containerInfoStack = containerInfoStack[:len(containerInfoStack)-1]
			currentDepth--
			if len(containerInfoStack) > 0 && containerInfoStack[len(containerInfoStack)-1].Type == '{' {
				containerInfoStack[len(containerInfoStack)-1].isKey = true
			}

		case '"':
			// Manually scan the string to count its raw byte length and get the value.
			val, err := scanString(bufReader, limits.MaxStringValueLength)
			if err != nil {
				return 0, err
			}

			if currentContainer == nil {
				if maxDepth == 0 {
					maxDepth = 1
				}
				break
			}

			if currentContainer.Type == '{' {
				if currentContainer.isKey {
					// Check for duplicate keys.
					if _, ok := currentContainer.keys[val]; ok {
						return 0, fmt.Errorf("duplicate key '%s' in object", val)
					}
					currentContainer.keys[val] = struct{}{}

					currentContainer.Count++
					if currentContainer.Count > limits.MaxObjectEntryCount {
						return 0, fmt.Errorf("JSON object exceeds allowed entry count")
					}
					currentContainer.isKey = false
				} else {
					currentContainer.isKey = true
				}
			} else {
				currentContainer.Count++
				if currentContainer.Count > limits.MaxArrayElementCount {
					return 0, fmt.Errorf("JSON array exceeds allowed element count")
				}
			}

		case 't', 'f', 'n': // true, false, null
			if err := scanLiteral(bufReader); err != nil {
				return 0, err
			}
			if currentContainer == nil {
				if maxDepth == 0 {
					maxDepth = 1
				}
				break
			}
			if currentContainer.Type == '[' {
				currentContainer.Count++
				if currentContainer.Count > limits.MaxArrayElementCount {
					return 0, fmt.Errorf("JSON array exceeds allowed element count")
				}
			} else {
				currentContainer.isKey = true
			}
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if err := scanNumber(bufReader, limits.MaxStringValueLength); err != nil {
				return 0, err
			}
			if currentContainer == nil {
				if maxDepth == 0 {
					maxDepth = 1
				}
				break
			}
			if currentContainer.Type == '[' {
				currentContainer.Count++
				if currentContainer.Count > limits.MaxArrayElementCount {
					return 0, fmt.Errorf("JSON array exceeds allowed element count")
				}
			} else {
				currentContainer.isKey = true
			}

		case ',':
			_, _ = bufReader.ReadByte()
			lastTokenWasComma = true
			if currentContainer != nil && currentContainer.Type == '{' {
				currentContainer.isKey = true
			}

		case ':':
			_, _ = bufReader.ReadByte()

		default:
			return 0, fmt.Errorf("invalid character '%c' looking for beginning of value", b[0])
		}

		if err := skipWhitespace(bufReader); err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}
	}

	if len(containerInfoStack) != 0 {
		return 0, fmt.Errorf("malformed JSON, unclosed containers")
	}

	return maxDepth, nil
}

func skipWhitespace(r *bufio.Reader) error {
	for {
		b, err := r.Peek(1)
		if err != nil {
			return err
		}
		if !isWhitespace(b[0]) {
			return nil
		}
		_, _ = r.ReadByte()
	}
}

// scanString consumes a JSON string from the reader, ensuring the raw byte
// length of its content does not exceed the limit. It returns the unescaped
// string value.
func scanString(r *bufio.Reader, limit int) (string, error) {
	if b, _ := r.ReadByte(); b != '"' {
		return "", fmt.Errorf("expected string")
	}

	var builder strings.Builder
	contentByteCount := 0
	var lastRune rune = -1 // Track the last rune for surrogate pair validation.

	for {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return "", fmt.Errorf("malformed JSON, unclosed string")
			}
			return "", err
		}

		if b == '"' {
			// Before successfully returning, ensure we didn't end on an unpaired high surrogate.
			if lastRune >= 0xD800 && lastRune <= 0xDBFF {
				return "", fmt.Errorf("malformed JSON, unterminated surrogate pair in string")
			}
			return builder.String(), nil
		}

		contentByteCount++
		lastRune = -1 // Reset unless we parse a new rune.

		if b == '\\' {
			escaped, err := r.ReadByte()
			if err != nil {
				return "", fmt.Errorf("malformed JSON, unterminated escape sequence in string")
			}
			contentByteCount++

			switch escaped {
			case '"', '\\', '/':
				builder.WriteByte(escaped)
			case 'b':
				builder.WriteByte('\b')
			case 'f':
				builder.WriteByte('\f')
			case 'n':
				builder.WriteByte('\n')
			case 'r':
				builder.WriteByte('\r')
			case 't':
				builder.WriteByte('\t')
			case 'u':
				hexChars := make([]byte, 4)
				if _, err := io.ReadFull(r, hexChars); err != nil {
					return "", fmt.Errorf("malformed JSON, unterminated unicode escape in string")
				}
				contentByteCount += 4

				hexStr := string(hexChars)
				for _, char := range hexStr {
					if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
						return "", fmt.Errorf("invalid character '%c' in string escape code", char)
					}
				}

				code, _ := strconv.ParseUint(hexStr, 16, 32)
				if code > unicode.MaxRune {
					return "", fmt.Errorf("invalid unicode escape sequence: value out of range")
				}

				r := rune(code)
				builder.WriteRune(r)
				lastRune = r
			default:
				return "", fmt.Errorf("invalid character '%c' in string escape code", escaped)
			}
		} else {
			builder.WriteByte(b)
		}

		if limit > 0 && contentByteCount > limit {
			return "", fmt.Errorf("JSON string value exceeds allowed length")
		}
	}
}

func scanLiteral(r *bufio.Reader) error {
	for {
		b, err := r.Peek(1)
		if err != nil {
			// If we hit EOF after reading part of a literal, that's a clean end.
			if err == io.EOF {
				break
			}
			return err
		}
		if isWhitespace(b[0]) || b[0] == ',' || b[0] == '}' || b[0] == ']' {
			return nil
		}
		_, _ = r.ReadByte()
	}
	return nil
}

func scanNumber(r *bufio.Reader, limit int) error {
	var builder strings.Builder
	byteCount := 0

	// Peek at the first char to check for leading zero issues.
	peeked, err := r.Peek(2)
	if err == nil && len(peeked) > 1 {
		if peeked[0] == '0' && peeked[1] >= '0' && peeked[1] <= '9' {
			_, _ = r.ReadByte()
			return fmt.Errorf("invalid character '%c' after top-level value", peeked[1])
		}
	}

	for {
		b, err := r.Peek(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		char := b[0]
		if isWhitespace(char) || char == ',' || char == '}' || char == ']' {
			break
		}

		// A number token can contain only these characters.
		isNumPart := (char >= '0' && char <= '9') || char == '.' || char == 'e' || char == 'E' || char == '+' || char == '-'
		if !isNumPart {
			// If it's not a valid number character, we stop scanning.
			break
		}

		_, _ = r.ReadByte()
		builder.WriteByte(char)
		byteCount++
		if limit > 0 && byteCount > limit {
			return fmt.Errorf("JSON number value exceeds allowed length")
		}
	}

	if byteCount == 0 {
		return fmt.Errorf("malformed JSON, empty number")
	}

	// Use the standard library for a final, strict validation of the number's syntax.
	// This correctly rejects malformed inputs like "-" or "123.".
	numStr := builder.String()
	if _, err := strconv.ParseFloat(numStr, 64); err != nil {
		return fmt.Errorf("malformed JSON, invalid number syntax for '%s'", numStr)
	}

	return nil
}
