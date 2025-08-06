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
	"github.com/stretchr/testify/require"
)

const (
	// CustomMaxJSONDepth specifies the maximum nesting depth of a JSON object.
	// This limit is designed to prevent stack exhaustion attacks from deeply
	// nested JSON payloads, which could otherwise lead to a denial-of-service
	// (DoS) vulnerability. The default value of 500 is intentionally generous
	// to support complex but legitimate configurations, while still providing
	// a safeguard against malicious or malformed input. This value is
	// configurable to accommodate unique environmental requirements.
	CustomMaxJSONDepth = 500

	// CustomMaxJSONStringValueLength defines the maximum allowed length for a single
	// string value within a JSON payload, in bytes. This is a critical defense
	// against excessive memory allocation attacks where a client might send a
	// very large string value to exhaust server memory. The default of 1MB
	// (1024 * 1024 bytes) is chosen to comfortably accommodate large secrets
	// such as private keys, certificate chains, or detailed configuration data,
	// without permitting unbounded allocation. This value is configurable.
	CustomMaxJSONStringValueLength = 1024 * 1024 // 1MB

	// CustomMaxJSONObjectEntryCount sets the maximum number of key-value pairs
	// allowed in a single JSON object. This limit helps mitigate the risk of
	// hash-collision denial-of-service (HashDoS) attacks and prevents general
	// resource exhaustion from parsing objects with an excessive number of
	// entries. A default of 10,000 entries is well beyond the scope of typical
	// Vault secrets or configurations, providing a high ceiling for normal
	// operations while ensuring stability. This value is configurable.
	CustomMaxJSONObjectEntryCount = 10000

	// CustomMaxJSONArrayElementCount determines the maximum number of elements
	// permitted in a single JSON array. This is particularly relevant for API
	// endpoints that can return large lists, such as the result of a `LIST`
	// operation on a secrets engine path. The default limit of 10,000 elements
	// prevents a single request from causing excessive memory consumption. While
	// most environments will fall well below this limit, it is configurable for
	// systems that require handling larger datasets, though pagination is the
	// recommended practice for such cases.
	CustomMaxJSONArrayElementCount = 10000
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

func TestJSONUtil_Limits(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		expectError bool
		errorMsg    string
	}{
		// Depth Limits
		{
			name:        "JSON exceeding max depth",
			jsonInput:   generateComplexJSON(CustomMaxJSONDepth + 1),
			expectError: true,
			errorMsg:    "JSON input exceeds allowed nesting depth",
		},
		{
			name:        "JSON at max allowed depth",
			jsonInput:   generateComplexJSON(CustomMaxJSONDepth),
			expectError: false,
		},
		// Malformed JSON
		{
			name:        "Malformed - Unmatched opening brace",
			jsonInput:   `{"a": {`,
			expectError: true,
			errorMsg:    "malformed JSON, unclosed containers",
		},
		{
			name:        "Malformed - Unmatched closing brace",
			jsonInput:   `{}}`,
			expectError: true,
			errorMsg:    "error reading JSON token: invalid character '}' looking for beginning of value",
		},
		// String Length Limits
		{
			name:        "String value exceeding max length",
			jsonInput:   fmt.Sprintf(`{"key": "%s"}`, strings.Repeat("a", CustomMaxJSONStringValueLength+1)),
			expectError: true,
			errorMsg:    "JSON string value exceeds allowed length",
		},
		{
			name:        "String at max length",
			jsonInput:   fmt.Sprintf(`{"key": "%s"}`, strings.Repeat("a", CustomMaxJSONStringValueLength)),
			expectError: false,
		},
		// Object Entry Count Limits
		{
			name:        "Object exceeding max entry count",
			jsonInput:   fmt.Sprintf(`{%s}`, generateObjectEntries(CustomMaxJSONObjectEntryCount+1)),
			expectError: true,
			errorMsg:    "JSON object exceeds allowed entry count",
		},
		{
			name:        "Object at max entry count",
			jsonInput:   fmt.Sprintf(`{%s}`, generateObjectEntries(CustomMaxJSONObjectEntryCount)),
			expectError: false,
		},
		// Array Element Count Limits
		{
			name:        "Array exceeding max element count",
			jsonInput:   fmt.Sprintf(`[%s]`, generateArrayElements(CustomMaxJSONArrayElementCount+1)),
			expectError: true,
			errorMsg:    "JSON array exceeds allowed element count",
		},
		{
			name:        "Array at max element count",
			jsonInput:   fmt.Sprintf(`[%s]`, generateArrayElements(CustomMaxJSONArrayElementCount)),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limits := JSONLimits{
				MaxDepth:             CustomMaxJSONDepth,
				MaxStringValueLength: CustomMaxJSONStringValueLength,
				MaxObjectEntryCount:  CustomMaxJSONObjectEntryCount,
				MaxArrayElementCount: CustomMaxJSONArrayElementCount,
			}

			_, err := VerifyMaxDepthStreaming(bytes.NewReader([]byte(tt.jsonInput)), limits)

			if tt.expectError {
				require.Error(t, err, "expected an error but got nil")
				require.Contains(t, err.Error(), tt.errorMsg, "error message mismatch")
			} else {
				require.NoError(t, err, "did not expect an error but got one")
			}
		})
	}
}

// generateComplexJSON generates a valid JSON string with a specified nesting depth.
func generateComplexJSON(depth int) string {
	if depth <= 0 {
		return "{}"
	}
	// Build the nested structure from the inside out.
	json := "1"
	for i := 0; i < depth; i++ {
		json = fmt.Sprintf(`{"a":%s}`, json)
	}
	return json
}

// generateObjectEntries creates a string of object entries for testing.
func generateObjectEntries(count int) string {
	var sb strings.Builder
	for i := 0; i < count; i++ {
		sb.WriteString(fmt.Sprintf(`"key%d":%d`, i, i))
		if i < count-1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}

// generateArrayElements creates a string of array elements for testing.
func generateArrayElements(count int) string {
	var sb strings.Builder
	for i := 0; i < count; i++ {
		sb.WriteString(fmt.Sprintf("%d", i))
		if i < count-1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}
