// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package random

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestJSONMarshalling(t *testing.T) {
	expected := serializableRules{
		CharsetRule{
			Charset:  LowercaseRuneset,
			MinChars: 1,
		},
		CharsetRule{
			Charset:  UppercaseRuneset,
			MinChars: 1,
		},
		CharsetRule{
			Charset:  NumericRuneset,
			MinChars: 1,
		},
		CharsetRule{
			Charset:  ShortSymbolRuneset,
			MinChars: 1,
		},
	}

	marshalled, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("no error expected, got: %s", err)
	}

	actual := serializableRules{}
	err = json.Unmarshal(marshalled, &actual)
	if err != nil {
		t.Fatalf("no error expected, got: %s", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Actual: %#v\nExpected: %#v", actual, expected)
	}
}

func TestRunes_UnmarshalJSON(t *testing.T) {
	data := []byte(`"noaw8hgfsdjlkfsj3"`)

	expected := runes([]rune("noaw8hgfsdjlkfsj3"))
	actual := runes{}
	err := (&actual).UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("no error expected, got: %s", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Actual: %#v\nExpected: %#v", actual, expected)
	}
}
