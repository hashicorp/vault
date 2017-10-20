package vault

import (
	"reflect"
	"testing"
)

func TestIdentityStore_parseMetadata(t *testing.T) {
	goodKVs := []string{
		"key1=value1",
		"key2=value1=value2",
	}
	expectedMap := map[string]string{
		"key1": "value1",
		"key2": "value1=value2",
	}

	actualMap, err := parseMetadata(goodKVs)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedMap, actualMap) {
		t.Fatalf("bad: metadata; expected: %#v\n, actual: %#v\n", expectedMap, actualMap)
	}

	badKV := []string{
		"=world",
	}
	actualMap, err = parseMetadata(badKV)
	if err == nil {
		t.Fatalf("expected an error; got: %#v", actualMap)
	}

	badKV[0] = "world"
	actualMap, err = parseMetadata(badKV)
	if err == nil {
		t.Fatalf("expected an error: %#v", actualMap)
	}
}
