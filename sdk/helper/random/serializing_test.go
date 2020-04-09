package random

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestJSONMarshalling(t *testing.T) {
	expected := serializableRules{
		&CharsetRestriction{
			Charset:  LowercaseRuneset,
			MinChars: 1,
		},
		&CharsetRestriction{
			Charset:  UppercaseRuneset,
			MinChars: 1,
		},
		&CharsetRestriction{
			Charset:  NumericRuneset,
			MinChars: 1,
		},
		&CharsetRestriction{
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
