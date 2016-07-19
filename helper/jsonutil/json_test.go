package jsonutil

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_EncodeJSON(t *testing.T) {
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
		t.Fatal("bad: encoded JSON: expected:%s\nactual:%s\n", expected, string(actualBytes))
	}
}

func Test_DecodeJSON(t *testing.T) {
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
		t.Fatal("bad: expected:%#v\nactual:%#v", expected, actual)
	}
}

func Test_DecodeJSONFromReader(t *testing.T) {
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
		t.Fatal("bad: expected:%#v\nactual:%#v", expected, actual)
	}
}
