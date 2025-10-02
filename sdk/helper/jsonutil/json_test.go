// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jsonutil

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
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

func TestJSONUtil_Limits_DefaultLimits(t *testing.T) {
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
			errorMsg:    "invalid character '}' after top-level value",
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

			_, err := VerifyMaxDepthStreaming(bytes.NewReader([]byte(tt.jsonInput)), limits, nil)

			if tt.expectError {
				require.Error(t, err, "expected an error but got nil")
				require.Contains(t, err.Error(), tt.errorMsg, "error message mismatch")
			} else {
				require.NoError(t, err, "did not expect an error but got one")
			}
		})
	}
}

func TestJSONUtil_Limits_ConfiguredLimits(t *testing.T) {
	limits := JSONLimits{
		MaxDepth:             64,
		MaxStringValueLength: 1024,
		MaxObjectEntryCount:  3,
		MaxArrayElementCount: 3,
		MaxTokens:            20,
	}

	bom := []byte{0xEF, 0xBB, 0xBF}

	tests := []struct {
		name     string
		payload  []byte
		errorMsg string
	}{
		{
			name:     "object entries with string values",
			payload:  []byte(`{"k0":"v0","k1":"v1","k2":"v2","k3":"v3"}`),
			errorMsg: "JSON object exceeds allowed entry count",
		},
		{
			name:     "object entries with array values",
			payload:  []byte(`{"k0":[],"k1":[],"k2":[],"k3":[]}`),
			errorMsg: "JSON object exceeds allowed entry count",
		},
		{
			name:     "object entries with object values",
			payload:  []byte(`{"k0":{},"k1":{},"k2":{},"k3":{}}`),
			errorMsg: "JSON object exceeds allowed entry count",
		},
		{
			name:     "array elements as objects",
			payload:  []byte(`[{}, {}, {}, {}]`),
			errorMsg: "JSON array exceeds allowed element count",
		},
		{
			name:     "BOM-prefixed over-limit object",
			payload:  append(bom, []byte(`{"k0":"v0","k1":"v1","k2":"v2","k3":"v3"}`)...),
			errorMsg: "JSON object exceeds allowed entry count",
		},
		{
			name:     "object key exceeds string length limit",
			payload:  []byte(fmt.Sprintf(`{"%s": 0}`, strings.Repeat("a", limits.MaxStringValueLength+1))),
			errorMsg: "JSON string value exceeds allowed length",
		},
		{
			name:     "trailing data after valid JSON",
			payload:  []byte(`{"k0":"v0"} "invalid"`),
			errorMsg: "invalid character '\"' after top-level value",
		},
		{
			name:     "object with embedded null byte in key",
			payload:  []byte(`{"k0\u0000":0, "k1":1, "k2":2, "k3":3}`),
			errorMsg: "JSON object exceeds allowed entry count",
		},
		{
			name:     "incomplete JSON stream",
			payload:  []byte(`{"k0":"v0",`),
			errorMsg: "malformed JSON, unclosed containers",
		},
		{
			name:     "deeply nested object exceeds depth limit",
			payload:  []byte(strings.Repeat(`{"a":`, limits.MaxDepth+1) + "null" + strings.Repeat(`}`, limits.MaxDepth+1)),
			errorMsg: "JSON payload exceeds allowed token count",
		},
		{
			name:     "payload exceeds token limit",
			payload:  []byte(`{"a":1,"b":2,"c":3,"d":4,"e":5,"f":6,"g":7,"h":8,"i":9,"j":10}`),
			errorMsg: "JSON object exceeds allowed entry count",
		},
		{
			name:     "string with many escapes exceeds length limit",
			payload:  []byte(fmt.Sprintf(`{"k":"%s"}`, strings.Repeat(`\"`, limits.MaxStringValueLength/2+1))),
			errorMsg: "JSON string value exceeds allowed length",
		},
		{
			name: "deeply nested string exceeds length limit",
			payload: []byte(fmt.Sprintf(`%s{"key":"%s"}%s`,
				strings.Repeat(`{"a":`, 60),
				strings.Repeat("b", limits.MaxStringValueLength+1),
				strings.Repeat(`}`, 60))),
			errorMsg: "JSON payload exceeds allowed token count",
		},
		{
			name:     "very long number exceeds length limit",
			payload:  []byte(fmt.Sprintf(`{"key":%s}`, strings.Repeat("1", limits.MaxStringValueLength+1))),
			errorMsg: "JSON number value exceeds allowed length",
		},
		{
			name: "string with invalid unicode escape",
			// 'X' is not a valid hex digit
			payload:  []byte(`{"key":"\u123X"}`),
			errorMsg: "invalid character 'X' in string escape code",
		},
		{
			name:     "object with trailing comma",
			payload:  []byte(`{"k0":"v0",}`),
			errorMsg: "invalid character '}' after object key-value pair",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := VerifyMaxDepthStreaming(bytes.NewReader(tt.payload), limits, nil)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.errorMsg)
		})
	}
}

func TestVerifyMaxDepthStreaming_MaxTokens(t *testing.T) {
	t.Run("payload exceeds limit", func(t *testing.T) {
		limits := JSONLimits{
			MaxTokens:            5, // Set a small, specific limit.
			MaxDepth:             CustomMaxJSONDepth,
			MaxStringValueLength: CustomMaxJSONStringValueLength,
			MaxObjectEntryCount:  CustomMaxJSONObjectEntryCount,
			MaxArrayElementCount: CustomMaxJSONArrayElementCount,
		}

		// This payload contains 6 tokens: {, "k0", "v0", "k1", "v1", }
		payload := []byte(`{"k0":"v0","k1":"v1"}`)
		expectedErrorMsg := "JSON payload exceeds allowed token count"

		_, err := VerifyMaxDepthStreaming(bytes.NewReader(payload), limits, nil)

		// We expect an error because the token count (6) is greater than the limit (5).
		require.Error(t, err)
		require.Contains(t, err.Error(), expectedErrorMsg)
	})

	t.Run("payload within limit", func(t *testing.T) {
		limits := JSONLimits{
			MaxTokens:            5,
			MaxDepth:             CustomMaxJSONDepth,
			MaxStringValueLength: CustomMaxJSONStringValueLength,
			MaxObjectEntryCount:  CustomMaxJSONObjectEntryCount,
			MaxArrayElementCount: CustomMaxJSONArrayElementCount,
		}

		// This payload contains 3 tokens: {, "key", }
		payload := []byte(`{"key":null}`)

		_, err := VerifyMaxDepthStreaming(bytes.NewReader(payload), limits, nil)

		// We expect no error because the token count (3) is less than the limit (5).
		require.NoError(t, err)
	})
}

// TestJSONUtil_Limits_Strictness adds tests for cases that a lenient parser
// might accept but a security-focused one should reject.
func TestJSONUtil_Limits_Strictness(t *testing.T) {
	limits := JSONLimits{
		MaxDepth:             64,
		MaxStringValueLength: 1024,
		MaxObjectEntryCount:  3,
		MaxArrayElementCount: 3,
		MaxTokens:            100,
	}

	tests := []struct {
		name     string
		payload  []byte
		errorMsg string
	}{
		// RFC 8259 states that object key names SHOULD be unique, but doesn't
		// require it. A strict parser should reject duplicates to prevent ambiguity.
		{
			name:     "object with duplicate keys",
			payload:  []byte(`{"key":"v1", "key":"v2"}`),
			errorMsg: "duplicate key 'key' in object",
		},
		{
			name:     "array with trailing comma",
			payload:  []byte(`[1, 2, 3,]`),
			errorMsg: "invalid character ']' after array element",
		},
		// A robust parser should reject any invalid escape sequence, not just unicode.
		{
			name:     "string with invalid escape sequence",
			payload:  []byte(`{"key":"\q"}`),
			errorMsg: "invalid character 'q' in string escape code",
		},
		// A key must be followed by a colon and a value.
		{
			name:     "object with missing value after key",
			payload:  []byte(`{"key":}`),
			errorMsg: "invalid character '}' after object key",
		},
		// Numbers starting with zero (unless they are just "0") are not standard.
		{
			name:     "number with leading zero",
			payload:  []byte(`[0123]`),
			errorMsg: "invalid character '1' after top-level value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := VerifyMaxDepthStreaming(bytes.NewReader(tt.payload), limits, nil)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.errorMsg)
		})
	}
}

// TestVerifyMaxDepthStreaming_NonContainerBypass ensures that top-level values
// that are not objects or arrays are correctly ignored, as the limits are not
// intended to apply to them.
func TestVerifyMaxDepthStreaming_NonContainerBypass(t *testing.T) {
	limits := JSONLimits{MaxDepth: 1, MaxTokens: 1}

	tests := map[string][]byte{
		"top-level string": []byte(`"this is a string"`),
		"top-level number": []byte(`12345`),
		"top-level bool":   []byte(`true`),
		"top-level null":   []byte(`null`),
	}

	for name, payload := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := VerifyMaxDepthStreaming(bytes.NewReader(payload), limits, nil)
			require.NoError(t, err, "expected no error for non-container top-level value")
		})
	}
}

// TestVerifyMaxDepthStreaming_ValidVaultPayloads ensures the parser
// correctly accepts legitimate, common JSON payloads from the Vault ecosystem.
func TestVerifyMaxDepthStreaming_ValidVaultPayloads(t *testing.T) {
	// Use reasonable limits that are well above what these payloads require,
	// ensuring that the parser doesn't fail on valid structures.
	limits := JSONLimits{
		MaxDepth:             10,
		MaxStringValueLength: 4096,
		MaxObjectEntryCount:  100,
		MaxArrayElementCount: 100,
		MaxTokens:            500,
	}

	tests := map[string]string{
		"KVv2 secret read response": `{
			"request_id": "a5a1c058-305f-3576-2f1d-f8f9a46a742c",
			"lease_id": "",
			"lease_duration": 0,
			"renewable": false,
			"data": {
				"data": {
					"foo": "bar"
				},
				"metadata": {
					"created_time": "2025-09-19T08:10:00.123456789Z",
					"custom_metadata": null,
					"deletion_time": "",
					"destroyed": false,
					"version": 1
				}
			},
			"warnings": null,
			"wrap_info": null
		}`,
		"Auth token lookup response": `{
			"request_id": "6d1f2b3e-7c3a-4e2b-8c6a-1b7d5f0e3a1b",
			"data": {
				"accessor": "St8oY1x3x6z5y9p6q3r8s7t2",
				"creation_time": 1663242230,
				"display_name": "userpass-user",
				"entity_id": "e-12345-67890-abcdef",
				"expire_time": "2025-10-19T10:10:00.000Z",
				"explicit_max_ttl": 0,
				"id": "h.123abcde456fghij789klmno",
				"identity_policies": ["default", "dev-policy"],
				"issue_time": "2025-09-19T10:10:00.000Z",
				"meta": {
					"username": "test-user"
				},
				"num_uses": 0,
				"orphan": true,
				"path": "auth/userpass/login/test-user",
				"policies": ["default", "dev-policy"],
				"renewable": true,
				"ttl": 2764799,
				"type": "service"
			}
		}`,
		"LIST operation response": `{
			"request_id": "c1a2b3d4-e5f6-a7b8-c9d0-e1f2a3b4c5d6",
			"data": {
				"keys": [
					"secret1",
					"secret2/",
					"another-secret"
				]
			}
		}`,
		"Policy write request": `{
			"policy": "path \"secret/data/foo\" {\n  capabilities = [\"read\", \"list\"]\n}\n\npath \"secret/data/bar\" {\n  capabilities = [\"create\", \"update\"]\n}"
		}`,
		"Transit batch encryption request": `{
			"batch_input": [
				{
					"plaintext": "aGVsbG8gd29ybGQ=",
					"context": "Y29udGV4dDE="
				},
				{
					"plaintext": "dGhpcyBpcyBhIHRlc3Q="
				}
			]
		}`,
	}

	for name, payload := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := VerifyMaxDepthStreaming(bytes.NewReader([]byte(payload)), limits, nil)
			require.NoError(t, err, "expected valid Vault payload to parse without error")
		})
	}
}

// TestVerifyMaxDepthStreaming_InvalidVaultPayloads ensures the parser correctly
// rejects real-world Vault payloads that have been crafted to violate specific
// security limits.
func TestVerifyMaxDepthStreaming_InvalidVaultPayloads(t *testing.T) {
	tests := []struct {
		name     string
		payload  []byte
		errorMsg string
		limits   JSONLimits
	}{
		// This KVv2 secret is valid, but the metadata object contains 5 keys.
		// It should be rejected by the MaxObjectEntryCount limit of 4.
		{
			name: "KVv2 secret with too many metadata entries",
			payload: []byte(`{
				"data": {
					"data": {"foo": "bar"},
					"metadata": {
						"created_time": "2025-09-19T08:10:00.123456789Z",
						"custom_metadata": null,
						"deletion_time": "",
						"destroyed": false,
						"version": 1
					}
				}
			}`),
			errorMsg: "JSON input exceeds allowed nesting depth",
			limits: JSONLimits{
				MaxDepth:             2,
				MaxStringValueLength: 100,
				MaxObjectEntryCount:  4,
				MaxArrayElementCount: 2,
				MaxTokens:            40,
			},
		},
		// This auth token response is flat but contains many key-value pairs,
		// resulting in over 50 tokens. It should be rejected by the MaxTokens limit of 40.
		{
			name: "Auth token response with too many tokens",
			payload: []byte(`{
				"data": {
					"accessor": "St8oY1x3x6z5y9p6q3r8s7t2",
					"creation_time": 1663242230,
					"display_name": "userpass-user",
					"entity_id": "e-12345-67890-abcdef",
					"expire_time": "2025-10-19T10:10:00.000Z",
					"explicit_max_ttl": 0,
					"id": "h.123abcde456fghij789klmno",
					"identity_policies": ["default", "dev-policy"],
					"issue_time": "2025-09-19T10:10:00.000Z",
					"meta": {"username": "test-user"},
					"num_uses": 0,
					"orphan": true,
					"path": "auth/userpass/login/test-user",
					"policies": ["default", "dev-policy"],
					"renewable": true,
					"ttl": 2764799,
					"type": "service"
				}
			}`),
			errorMsg: "JSON payload exceeds allowed token count",
			limits: JSONLimits{
				MaxDepth:             100,
				MaxStringValueLength: 100,
				MaxObjectEntryCount:  50,
				MaxArrayElementCount: 2,
				MaxTokens:            40,
			},
		},
		// This LIST response is valid but contains 3 elements in the "keys" array.
		// It should be rejected by the MaxArrayElementCount limit of 2.
		{
			name: "LIST response with too many keys in array",
			payload: []byte(`{
				"data": {
					"keys": [
						"secret1",
						"secret2/",
						"another-secret"
					]
				}
			}`),
			errorMsg: "JSON input exceeds allowed nesting depth",
			limits: JSONLimits{
				MaxDepth:             2,
				MaxStringValueLength: 100,
				MaxObjectEntryCount:  4,
				MaxArrayElementCount: 2,
				MaxTokens:            40,
			},
		},
		// The policy string in this payload is over 100 bytes long.
		// It should be rejected by the MaxStringValueLength limit of 100.
		{
			name: "Policy write request with oversized policy string",
			payload: []byte(`{
				"policy": "path \"secret/data/foo\" {\n  capabilities = [\"read\", \"list\"]\n}\n\npath \"secret/data/bar\" {\n  capabilities = [\"create\", \"update\"]\n}"
			}`),
			errorMsg: "JSON string value exceeds allowed length",
			limits: JSONLimits{
				MaxDepth:             2,
				MaxStringValueLength: 100,
				MaxObjectEntryCount:  4,
				MaxArrayElementCount: 2,
				MaxTokens:            40,
			},
		},
		// This transit request has a nesting depth of 3 ({ -> [ -> {).
		// It should be rejected by the MaxDepth limit of 2.
		{
			name: "Transit batch request with excessive depth",
			payload: []byte(`{
				"batch_input": [
					{
						"plaintext": "aGVsbG8gd29ybGQ="
					}
				]
			}`),
			errorMsg: "JSON input exceeds allowed nesting depth",
			limits: JSONLimits{
				MaxDepth:             2,
				MaxStringValueLength: 100,
				MaxObjectEntryCount:  4,
				MaxArrayElementCount: 2,
				MaxTokens:            40,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := VerifyMaxDepthStreaming(bytes.NewReader(tt.payload), tt.limits, nil)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.errorMsg)
		})
	}
}

// TestVerifyMaxDepthStreaming_MaxRequestSize ensures that the buffered reader
// used in VerifyMaxDepthStreaming correctly enforces the maxRequestSize limit.
// The test checks that a non-JSON payload can be processed without error and
// that no more than the maxRequestSize bytes are read from the reader.
func TestVerifyMaxDepthStreaming_MaxRequestSize(t *testing.T) {
	input := `this is test data that isn't json`
	bodyReader := strings.NewReader(input)
	read := new(bytes.Buffer)
	maxRequestSize := new(int64)
	*maxRequestSize = 16
	// this is meant to mimic http request body reading
	// the TeeReader will capture what is read from the bodyReader
	reader := io.TeeReader(http.MaxBytesReader(nil, io.NopCloser(bodyReader), *maxRequestSize), read)
	_, err := VerifyMaxDepthStreaming(reader, JSONLimits{MaxDepth: 5, MaxTokens: 5}, maxRequestSize)
	require.NoError(t, err)
	// then, check to see that combining the read buffer and the remaining body
	// results in the original input
	// if VerifyMaxDepthStreaming tried to read more than the
	// maxRequestSize, this check will fail because the http.MaxBytesReader will
	// have read an extra byte without writing it to the read buffer
	full := io.MultiReader(read, bodyReader)
	all, err := io.ReadAll(full)
	require.NoError(t, err)
	require.Equal(t, input, string(all))
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
