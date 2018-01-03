package framework

import (
	"reflect"
	"testing"
)

func TestFieldDataGet(t *testing.T) {
	cases := map[string]struct {
		Schema map[string]*FieldSchema
		Raw    map[string]interface{}
		Key    string
		Value  interface{}
	}{
		"string type, string value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeString},
			},
			map[string]interface{}{
				"foo": "bar",
			},
			"foo",
			"bar",
		},

		"string type, int value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeString},
			},
			map[string]interface{}{
				"foo": 42,
			},
			"foo",
			"42",
		},

		"string type, unset value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeString},
			},
			map[string]interface{}{},
			"foo",
			"",
		},

		"string type, unset value with default": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{
					Type:    TypeString,
					Default: "bar",
				},
			},
			map[string]interface{}{},
			"foo",
			"bar",
		},

		"int type, int value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeInt},
			},
			map[string]interface{}{
				"foo": 42,
			},
			"foo",
			42,
		},

		"bool type, bool value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeBool},
			},
			map[string]interface{}{
				"foo": false,
			},
			"foo",
			false,
		},

		"map type, map value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeMap},
			},
			map[string]interface{}{
				"foo": map[string]interface{}{
					"child": true,
				},
			},
			"foo",
			map[string]interface{}{
				"child": true,
			},
		},

		"duration type, string value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeDurationSecond},
			},
			map[string]interface{}{
				"foo": "42",
			},
			"foo",
			42,
		},

		"duration type, string duration value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeDurationSecond},
			},
			map[string]interface{}{
				"foo": "42m",
			},
			"foo",
			2520,
		},

		"duration type, int value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeDurationSecond},
			},
			map[string]interface{}{
				"foo": 42,
			},
			"foo",
			42,
		},

		"duration type, float value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeDurationSecond},
			},
			map[string]interface{}{
				"foo": 42.0,
			},
			"foo",
			42,
		},

		"duration type, nil value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeDurationSecond},
			},
			map[string]interface{}{
				"foo": nil,
			},
			"foo",
			0,
		},

		"slice type, empty slice": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeSlice},
			},
			map[string]interface{}{
				"foo": []interface{}{},
			},
			"foo",
			[]interface{}{},
		},

		"slice type, filled, mixed slice": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeSlice},
			},
			map[string]interface{}{
				"foo": []interface{}{123, "abc"},
			},
			"foo",
			[]interface{}{123, "abc"},
		},

		"string slice type, filled slice": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeStringSlice},
			},
			map[string]interface{}{
				"foo": []interface{}{123, "abc"},
			},
			"foo",
			[]string{"123", "abc"},
		},

		"string slice type, single value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeStringSlice},
			},
			map[string]interface{}{
				"foo": "abc",
			},
			"foo",
			[]string{"abc"},
		},

		"comma string slice type, comma string with one value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeCommaStringSlice},
			},
			map[string]interface{}{
				"foo": "value1",
			},
			"foo",
			[]string{"value1"},
		},

		"comma string slice type, comma string with multi value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeCommaStringSlice},
			},
			map[string]interface{}{
				"foo": "value1,value2,value3",
			},
			"foo",
			[]string{"value1", "value2", "value3"},
		},

		"comma string slice type, nil string slice value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeCommaStringSlice},
			},
			map[string]interface{}{
				"foo": "",
			},
			"foo",
			[]string{},
		},

		"commma string slice type, string slice with one value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeCommaStringSlice},
			},
			map[string]interface{}{
				"foo": []interface{}{"value1"},
			},
			"foo",
			[]string{"value1"},
		},

		"comma string slice type, string slice with multi value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeCommaStringSlice},
			},
			map[string]interface{}{
				"foo": []interface{}{"value1", "value2", "value3"},
			},
			"foo",
			[]string{"value1", "value2", "value3"},
		},

		"comma string slice type, empty string slice value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeCommaStringSlice},
			},
			map[string]interface{}{
				"foo": []interface{}{},
			},
			"foo",
			[]string{},
		},

		"name string type, valid string": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeNameString},
			},
			map[string]interface{}{
				"foo": "bar",
			},
			"foo",
			"bar",
		},

		"name string type, valid value with special characters": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeNameString},
			},
			map[string]interface{}{
				"foo": "bar.baz-bay123",
			},
			"foo",
			"bar.baz-bay123",
		},

		"keypair type, valid value map type": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeKVPairs},
			},
			map[string]interface{}{
				"foo": map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
					"key3": 1,
				},
			},
			"foo",
			map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "1",
			},
		},

		"keypair type, list of equal sign delim key pairs type": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeKVPairs},
			},
			map[string]interface{}{
				"foo": []interface{}{"key1=value1", "key2=value2", "key3=1"},
			},
			"foo",
			map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "1",
			},
		},

		"keypair type, single equal sign delim value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeKVPairs},
			},
			map[string]interface{}{
				"foo": "key1=value1",
			},
			"foo",
			map[string]string{
				"key1": "value1",
			},
		},

		"name string type, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeNameString},
			},
			map[string]interface{}{},
			"foo",
			"",
		},

		"string type, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeString},
			},
			map[string]interface{}{},
			"foo",
			"",
		},

		"type int, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeInt},
			},
			map[string]interface{}{},
			"foo",
			0,
		},

		"type bool, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeBool},
			},
			map[string]interface{}{},
			"foo",
			false,
		},

		"type map, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeMap},
			},
			map[string]interface{}{},
			"foo",
			map[string]interface{}{},
		},

		"type duration second, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeDurationSecond},
			},
			map[string]interface{}{},
			"foo",
			0,
		},

		"type slice, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeSlice},
			},
			map[string]interface{}{},
			"foo",
			[]interface{}{},
		},

		"type string slice, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeStringSlice},
			},
			map[string]interface{}{},
			"foo",
			[]string{},
		},

		"type comma string slice, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeCommaStringSlice},
			},
			map[string]interface{}{},
			"foo",
			[]string{},
		},

		"type kv pair, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeKVPairs},
			},
			map[string]interface{}{},
			"foo",
			map[string]string{},
		},
	}

	for name, tc := range cases {
		data := &FieldData{
			Raw:    tc.Raw,
			Schema: tc.Schema,
		}

		if err := data.Validate(); err != nil {
			t.Fatalf("bad: %#v", err)
		}

		actual := data.Get(tc.Key)
		if !reflect.DeepEqual(actual, tc.Value) {
			t.Fatalf(
				"bad: %s\n\nExpected: %#v\nGot: %#v",
				name, tc.Value, actual)
		}
	}
}

func TestFieldDataGet_Error(t *testing.T) {
	cases := map[string]struct {
		Schema map[string]*FieldSchema
		Raw    map[string]interface{}
		Key    string
	}{
		"name string type, invalid value with invalid characters": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeNameString},
			},
			map[string]interface{}{
				"foo": "bar baz",
			},
			"foo",
		},
		"name string type, invalid value with special characters at beginning": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeNameString},
			},
			map[string]interface{}{
				"foo": ".barbaz",
			},
			"foo",
		},
		"name string type, invalid value with special characters at end": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeNameString},
			},
			map[string]interface{}{
				"foo": "barbaz-",
			},
			"foo",
		},
		"name string type, empty string": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeNameString},
			},
			map[string]interface{}{
				"foo": "",
			},
			"foo",
		},
		"keypair type, csv version empty key name": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeKVPairs},
			},
			map[string]interface{}{
				"foo": []interface{}{"=value1", "key2=value2", "key3=1"},
			},
			"foo",
		},
	}

	for _, tc := range cases {
		data := &FieldData{
			Raw:    tc.Raw,
			Schema: tc.Schema,
		}

		_, _, err := data.GetOkErr(tc.Key)
		if err == nil {
			t.Fatalf("error expected, none received")
		}
	}
}
