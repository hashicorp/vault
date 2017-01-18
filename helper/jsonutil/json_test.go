package jsonutil

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/compressutil"
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
	for key, _ := range actual {
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

	// Compress an object
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
