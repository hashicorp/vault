package random

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestJSONMarshalling(t *testing.T) {
	expected := serializableRules{
		Charset{
			Charset:  LowercaseRuneset,
			MinChars: 1,
		},
		Charset{
			Charset:  UppercaseRuneset,
			MinChars: 1,
		},
		Charset{
			Charset:  NumericRuneset,
			MinChars: 1,
		},
		Charset{
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
