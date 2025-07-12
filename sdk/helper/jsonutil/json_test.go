// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jsonutil

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/compressutil"
)

func TestJSONUtil_CompressDecompressJSON(t *testing.T) {
	expected := map[string]interface{}{
		"test":       "data",
		"validation": "process",
	}

	// Compress an object
	compressedBytes, err := EncodeJSONAndCompress(expected, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(compressedBytes) == 0 {
		t.Fatal("expected compressed data")
	}

	// Check if canary is present in the compressed data
	if compressedBytes[0] != compressutil.CompressionCanaryGzip {
		t.Fatalf("canary missing in compressed data")
	}

	// Decompress and decode the compressed information and verify the functional
	// behavior
	var actual map[string]interface{}
	if err = DecodeJSON(compressedBytes, &actual); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: expected: %#v\nactual: %#v", expected, actual)
	}
	for key := range actual {
		delete(actual, key)
	}

	// Test invalid data
	if err = DecodeJSON([]byte{}, &actual); err == nil {
		t.Fatalf("expected a failure")
	}

	// Test invalid data after the canary byte
	var buf bytes.Buffer
	buf.Write([]byte{compressutil.CompressionCanaryGzip})
	if err = DecodeJSON(buf.Bytes(), &actual); err == nil {
		t.Fatalf("expected a failure")
	}

	// Compress an object with BestSpeed
	compressedBytes, err = EncodeJSONAndCompress(expected, &compressutil.CompressionConfig{
		Type:                 compressutil.CompressionTypeGzip,
		GzipCompressionLevel: gzip.BestSpeed,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(compressedBytes) == 0 {
		t.Fatal("expected compressed data")
	}

	// Check if canary is present in the compressed data
	if compressedBytes[0] != compressutil.CompressionCanaryGzip {
		t.Fatalf("canary missing in compressed data")
	}

	// Decompress and decode the compressed information and verify the functional
	// behavior
	if err = DecodeJSON(compressedBytes, &actual); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: expected: %#v\nactual: %#v", expected, actual)
	}
}

func TestJSONUtil_EncodeJSON(t *testing.T) {
	input := map[string]interface{}{
		"test":       "data",
		"validation": "process",
	}

	actualBytes, err := EncodeJSON(input)
	if err != nil {
		t.Fatalf("failed to encode JSON: %v", err)
	}

	actual := strings.TrimSpace(string(actualBytes))
	expected := `{"test":"data","validation":"process"}`

	if actual != expected {
		t.Fatalf("bad: encoded JSON: expected:%s\nactual:%s\n", expected, string(actualBytes))
	}
}

func TestJSONUtil_DecodeJSON(t *testing.T) {
	input := `{"test":"data","validation":"process"}`

	var actual map[string]interface{}

	err := DecodeJSON([]byte(input), &actual)
	if err != nil {
		fmt.Printf("decoding err: %v\n", err)
	}

	expected := map[string]interface{}{
		"test":       "data",
		"validation": "process",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, actual)
	}
}

func TestJSONUtil_DecodeJSONFromReader(t *testing.T) {
	input := `{"test":"data","validation":"process"}`

	var actual map[string]interface{}

	err := DecodeJSONFromReader(bytes.NewReader([]byte(input)), &actual)
	if err != nil {
		fmt.Printf("decoding err: %v\n", err)
	}

	expected := map[string]interface{}{
		"test":       "data",
		"validation": "process",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, actual)
	}
}

func TestJSONUtil_DecodeJSON_DepthChecks(t *testing.T) {
	tests := []struct {
		name        string
		totalDepth  int    // This 'totalDepth' is the parameter passed to generateComplexJSON
		jsonInput   string // Raw JSON string for direct input, overrides 'totalDepth' if not empty
		decodeFunc  func([]byte, interface{}) error
		expectError bool
		errorMsg    string
	}{
		{
			name:        "DecodeJSON: JSON exceeding max depth (total depth CustomMaxJSONDepth + 1)",
			totalDepth:  CustomMaxJSONDepth + 1,
			decodeFunc:  DecodeJSON,
			expectError: true,
			errorMsg:    "JSON input exceeds allowed nesting depth",
		},
		{
			name:        "DecodeJSON: JSON at exact max allowed depth (total depth CustomMaxJSONDepth)",
			totalDepth:  CustomMaxJSONDepth,
			decodeFunc:  DecodeJSON,
			expectError: false,
		},
		{
			name:        "DecodeJSON: JSON well within max depth",
			totalDepth:  CustomMaxJSONDepth / 2,
			decodeFunc:  DecodeJSON,
			expectError: false,
		},
		{
			name:       "DecodeJSONFromReader: JSON exceeding max depth (total depth CustomMaxJSONDepth + 1)",
			totalDepth: CustomMaxJSONDepth + 1,
			decodeFunc: func(data []byte, out interface{}) error {
				return DecodeJSONFromReader(bytes.NewReader(data), out)
			},
			expectError: true,
			errorMsg:    "JSON input exceeds allowed nesting depth",
		},
		{
			name:       "DecodeJSONFromReader: JSON at exact max allowed depth (total depth CustomMaxJSONDepth)",
			totalDepth: CustomMaxJSONDepth,
			decodeFunc: func(data []byte, out interface{}) error {
				return DecodeJSONFromReader(bytes.NewReader(data), out)
			},
			expectError: false,
		},
		{
			name:       "DecodeJSONFromReader: JSON well within max depth",
			totalDepth: CustomMaxJSONDepth / 2,
			decodeFunc: func(data []byte, out interface{}) error {
				return DecodeJSONFromReader(bytes.NewReader(data), out)
			},
			expectError: false,
		},
		{
			name:      "DecodeJSONFromReader: Empty JSON input (expected nil error)",
			jsonInput: "",
			decodeFunc: func(data []byte, out interface{}) error {
				return DecodeJSONFromReader(bytes.NewReader(data), out)
			},
			expectError: false,
			errorMsg:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var jsonStr string
			if tt.jsonInput != "" {
				jsonStr = tt.jsonInput
			} else {
				jsonStr = generateComplexJSON(tt.totalDepth) // Pass totalDepth directly
			}

			var actual interface{} // Use interface{} as we don't care about the content, just if it decodes

			err := tt.decodeFunc([]byte(jsonStr), &actual)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected an error containing '%s', but got nil", tt.errorMsg)
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Fatalf("expected error message to contain '%s', but got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("did not expect an error, but got: %v", err)
				}
			}
		})
	}
}

func TestJSONUtil_DecodeJSON_Bounded(t *testing.T) {
	tests := []struct {
		name        string
		depthLimit  int
		byteLimit   int
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Bounded JSON within limits",
			depthLimit:  CustomMaxJSONDepth - 1, // Ensure usableDepth is positive and within limit
			byteLimit:   1024,
			expectError: false,
		},
		{
			name:        "Bounded JSON exceeding depth limit",
			depthLimit:  CustomMaxJSONDepth + 1, // This will make usableDepth also exceed
			byteLimit:   1024,
			expectError: true,
			errorMsg:    "JSON input exceeds allowed nesting depth",
		},
		{
			name:        "Bounded JSON with minimal depth",
			depthLimit:  4, // usableDepth = 1
			byteLimit:   100,
			expectError: false,
		},
		{
			name:        "Bounded JSON with very small byte limit (should still decode if depth is okay)",
			depthLimit:  4,
			byteLimit:   50, // This might result in an empty inner array, but should still be valid JSON
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonStr, err := generateBoundedJSON(tt.depthLimit, tt.byteLimit)
			if err != nil {
				// This error is from the generator, not the decoder.
				t.Fatalf("failed to generate bounded JSON: %v", err)
			}

			var actual interface{}
			err = DecodeJSON([]byte(jsonStr), &actual)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected an error containing '%s', but got nil", tt.errorMsg)
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Fatalf("expected error message to contain '%s', but got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("did not expect an error, but got: %v", err)
				}
			}
		})
	}
}

func TestJSONUtil_CalculateMaxDepthStreaming(t *testing.T) {
	tests := []struct {
		name          string
		jsonInput     string
		expectedDepth int
		expectError   bool
		errorMsg      string
	}{
		{
			name:          "Empty JSON",
			jsonInput:     "",
			expectedDepth: 0,
			expectError:   false,
		},
		{
			name:          "Simple object",
			jsonInput:     `{"key": "value"}`,
			expectedDepth: 1,
			expectError:   false,
		},
		{
			name:          "Simple array",
			jsonInput:     `["item1", "item2"]`,
			expectedDepth: 1,
			expectError:   false,
		},
		{
			name:          "Nested object",
			jsonInput:     `{"a": {"b": "c"}}`,
			expectedDepth: 2,
			expectError:   false,
		},
		{
			name:          "Nested array",
			jsonInput:     `[["item"]]`,
			expectedDepth: 2,
			expectError:   false,
		},
		{
			name:          "Mixed nesting",
			jsonInput:     `{"a": [{"b": {"c": 1}}]}`,
			expectedDepth: 4,
			expectError:   false,
		},
		{
			name:          "JSON at CustomMaxJSONDepth",
			jsonInput:     generateComplexJSON(CustomMaxJSONDepth), // Pass total depth
			expectedDepth: CustomMaxJSONDepth,                      // Expected total depth
			expectError:   false,
		},
		{
			name:          "JSON exceeding CustomMaxJSONDepth",
			jsonInput:     generateComplexJSON(CustomMaxJSONDepth + 1), // Pass total depth
			expectedDepth: CustomMaxJSONDepth + 1,                      // Expected total depth
			expectError:   false,                                       // calculateMaxDepthStreaming itself doesn't error on depth, DecodeJSONFromReader does
		},
		{
			name:          "Invalid JSON - unclosed object",
			jsonInput:     `{"a": {`,
			expectedDepth: 2, // Actual depth reached before EOF
			expectError:   true,
			errorMsg:      "malformed JSON, unmatched delimiters",
		},
		{
			name:          "Invalid JSON - unclosed array",
			jsonInput:     `[`,
			expectedDepth: 1, // Actual depth reached before EOF
			expectError:   true,
			errorMsg:      "malformed JSON, unmatched delimiters",
		},
		{
			name:          "Invalid JSON - malformed",
			jsonInput:     `{"a":`,
			expectedDepth: 0, // Actual depth doesn't matter, error is expected
			expectError:   true,
			errorMsg:      "malformed JSON, unmatched delimiters",
		},
		{
			name:          "JSON with numbers and strings",
			jsonInput:     `{"key1": 123, "key2": "value", "key3": [1, 2, {"nested": 3}]}`,
			expectedDepth: 3,
			expectError:   false,
		},
		{
			name:          "JSON with null and boolean",
			jsonInput:     `{"a": null, "b": true}`,
			expectedDepth: 1,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader([]byte(tt.jsonInput))
			actualDepth, err := calculateMaxDepthStreaming(reader)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected an error containing '%s', but got nil", tt.errorMsg)
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Fatalf("expected error message to contain '%s', but got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("did not expect an error, but got: %v", err)
				}
				if actualDepth != tt.expectedDepth {
					t.Fatalf("bad depth: expected %d, got %d for JSON: %s", tt.expectedDepth, actualDepth, tt.jsonInput)
				}
			}
		})
	}
}

// generateComplexJSON generates a JSON string with a specified total nesting depth.
// It accounts for the fixed outer {"data":{"data":...}} structure.
func generateComplexJSON(totalDepth int) string {
	if totalDepth < 2 {
		return "" // Cannot generate a valid JSON with less than 2 depth due to fixed outer structure
	}

	// Calculate the depth for the inner repeating structure
	innerRepeatDepth := totalDepth - 2

	var innerBuilder strings.Builder
	innerBuilder.WriteString("{")
	i := 0
	prefixes := []string{"", "_", ".", ","}
	for _, prefix := range prefixes {
		for x := 32; x < 126; x++ {
			if x == 34 || x == 92 { // Avoid " and \
				continue
			}
			innerBuilder.WriteString(fmt.Sprintf(`"%s%c":%d,`, prefix, rune(x), i))
			i++
		}
	}
	innerBuilder.WriteString(`"~":`) // The key that will contain the next nested object

	var jsonBuilder strings.Builder
	jsonBuilder.WriteString(`{"data":{"data":`) // Initial nesting for the structure (2 levels)

	for k := 0; k < innerRepeatDepth; k++ {
		jsonBuilder.WriteString(innerBuilder.String())
	}
	jsonBuilder.WriteString(`0.1`) // The innermost value
	for k := 0; k < innerRepeatDepth; k++ {
		jsonBuilder.WriteString(`}`) // Close all nested objects from the innerRepeatDepth
	}
	jsonBuilder.WriteString(`}}`) // Close initial nesting

	return jsonBuilder.String()
}

// generateBoundedJSON generates a JSON string that respects a maximum depth and byte limit.
// It constructs a structure like {"data":{"data":[[],[],...[]]}} where the inner arrays
// are nested to `depthLimit - 3` and repeated to stay within `byteLimit`.
func generateBoundedJSON(depthLimit, byteLimit int) (string, error) {
	// usableDepth accounts for the outer {"data":{"data":...}} structure
	usableDepth := depthLimit - 3
	if usableDepth <= 0 {
		return "", fmt.Errorf("error: depth limit too small (must be > 2 to allow usable inner depth)")
	}

	prefix := `{"data":{"data":[`
	suffix := `]}}`

	// Calculate bytes for one instance of the nested array structure, e.g., "[[]]" for usableDepth=2
	// Each level adds 2 bytes (opening and closing bracket)
	nestedArrayBytes := 2 * usableDepth

	// Calculate remaining bytes available for the repeating arrays
	remainingBytes := byteLimit - len(prefix) - len(suffix)

	arrayCount := 0
	if nestedArrayBytes > 0 { // Avoid division by zero
		// +1 for the comma separator between arrays
		arrayCount = (remainingBytes + 1) / (nestedArrayBytes + 1)
		if remainingBytes < nestedArrayBytes { // If not even one full array fits
			arrayCount = 0
		}
	}

	var sb strings.Builder
	sb.WriteString(prefix)

	for i := 0; i < arrayCount; i++ {
		if i > 0 {
			sb.WriteRune(',') // Add comma separator for subsequent arrays
		}
		// Write opening brackets for the nested array
		for j := 0; j < usableDepth; j++ {
			sb.WriteRune('[')
		}
		// Write closing brackets for the nested array
		for j := 0; j < usableDepth; j++ {
			sb.WriteRune(']')
		}
	}
	sb.WriteString(suffix)

	// Basic check to ensure we didn't exceed the byte limit due to calculation
	if sb.Len() > byteLimit {
		return "", fmt.Errorf("generated JSON length (%d) exceeds byte limit (%d)", sb.Len(), byteLimit)
	}

	return sb.String(), nil
}
