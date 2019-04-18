package framework

import (
	"net/http"
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

		"lowercase string type, lowercase string value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeLowerCaseString},
			},
			map[string]interface{}{
				"foo": "bar",
			},
			"foo",
			"bar",
		},

		"lowercase string type, mixed-case string value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeLowerCaseString},
			},
			map[string]interface{}{
				"foo": "BaR",
			},
			"foo",
			"bar",
		},

		"lowercase string type, int value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeLowerCaseString},
			},
			map[string]interface{}{
				"foo": 42,
			},
			"foo",
			"42",
		},

		"lowercase string type, unset value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeLowerCaseString},
			},
			map[string]interface{}{},
			"foo",
			"",
		},

		"lowercase string type, unset value with lowercase default": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{
					Type:    TypeLowerCaseString,
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

		"string slice type, empty string": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeStringSlice},
			},
			map[string]interface{}{
				"foo": "",
			},
			"foo",
			[]string{},
		},

		"comma string slice type, empty string": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeCommaStringSlice},
			},
			map[string]interface{}{
				"foo": "",
			},
			"foo",
			[]string{},
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

		"comma string slice type, string slice with one value": {
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

		"comma int slice type, comma int with one value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeCommaIntSlice},
			},
			map[string]interface{}{
				"foo": 1,
			},
			"foo",
			[]int{1},
		},

		"comma int slice type, comma int with multi value slice": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeCommaIntSlice},
			},
			map[string]interface{}{
				"foo": []int{1, 2, 3},
			},
			"foo",
			[]int{1, 2, 3},
		},

		"comma int slice type, comma int with multi value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeCommaIntSlice},
			},
			map[string]interface{}{
				"foo": "1,2,3",
			},
			"foo",
			[]int{1, 2, 3},
		},

		"comma int slice type, nil int slice value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeCommaIntSlice},
			},
			map[string]interface{}{
				"foo": "",
			},
			"foo",
			[]int{},
		},

		"comma int slice type, int slice with one value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeCommaIntSlice},
			},
			map[string]interface{}{
				"foo": []interface{}{"1"},
			},
			"foo",
			[]int{1},
		},

		"comma int slice type, int slice with multi value strings": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeCommaIntSlice},
			},
			map[string]interface{}{
				"foo": []interface{}{"1", "2", "3"},
			},
			"foo",
			[]int{1, 2, 3},
		},

		"comma int slice type, int slice with multi value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeCommaIntSlice},
			},
			map[string]interface{}{
				"foo": []interface{}{1, 2, 3},
			},
			"foo",
			[]int{1, 2, 3},
		},

		"comma int slice type, empty int slice value": {
			map[string]*FieldSchema{
				"foo": &FieldSchema{Type: TypeCommaIntSlice},
			},
			map[string]interface{}{
				"foo": []interface{}{},
			},
			"foo",
			[]int{},
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

		"type header, keypair string array": {
			map[string]*FieldSchema{
				"foo": {Type: TypeHeader},
			},
			map[string]interface{}{
				"foo": []interface{}{"key1:value1", "key2:value2", "key3:1"},
			},
			"foo",
			http.Header{
				"Key1": []string{"value1"},
				"Key2": []string{"value2"},
				"Key3": []string{"1"},
			},
		},

		"type header, b64 string": {
			map[string]*FieldSchema{
				"foo": {Type: TypeHeader},
			},
			map[string]interface{}{
				"foo": "eyJDb250ZW50LUxlbmd0aCI6IFsiNDMiXSwgIlVzZXItQWdlbnQiOiBbImF3cy1zZGstZ28vMS40LjEyIChnbzEuNy4xOyBsaW51eDsgYW1kNjQpIl0sICJYLVZhdWx0LUFXU0lBTS1TZXJ2ZXItSWQiOiBbInZhdWx0LmV4YW1wbGUuY29tIl0sICJYLUFtei1EYXRlIjogWyIyMDE2MDkzMFQwNDMxMjFaIl0sICJDb250ZW50LVR5cGUiOiBbImFwcGxpY2F0aW9uL3gtd3d3LWZvcm0tdXJsZW5jb2RlZDsgY2hhcnNldD11dGYtOCJdLCAiQXV0aG9yaXphdGlvbiI6IFsiQVdTNC1ITUFDLVNIQTI1NiBDcmVkZW50aWFsPWZvby8yMDE2MDkzMC91cy1lYXN0LTEvc3RzL2F3czRfcmVxdWVzdCwgU2lnbmVkSGVhZGVycz1jb250ZW50LWxlbmd0aDtjb250ZW50LXR5cGU7aG9zdDt4LWFtei1kYXRlO3gtdmF1bHQtc2VydmVyLCBTaWduYXR1cmU9YTY5ZmQ3NTBhMzQ0NWM0ZTU1M2UxYjNlNzlkM2RhOTBlZWY1NDA0N2YxZWI0ZWZlOGZmYmM5YzQyOGMyNjU1YiJdLCAiRm9vIjogNDJ9",
			},
			"foo",
			http.Header{
				"Content-Length":           []string{"43"},
				"User-Agent":               []string{"aws-sdk-go/1.4.12 (go1.7.1; linux; amd64)"},
				"X-Vault-Awsiam-Server-Id": []string{"vault.example.com"},
				"X-Amz-Date":               []string{"20160930T043121Z"},
				"Content-Type":             []string{"application/x-www-form-urlencoded; charset=utf-8"},
				"Authorization":            []string{"AWS4-HMAC-SHA256 Credential=foo/20160930/us-east-1/sts/aws4_request, SignedHeaders=content-length;content-type;host;x-amz-date;x-vault-server, Signature=a69fd750a3445c4e553e1b3e79d3da90eef54047f1eb4efe8ffbc9c428c2655b"},
				"Foo":                      []string{"42"},
			},
		},

		"type header, json string": {
			map[string]*FieldSchema{
				"foo": {Type: TypeHeader},
			},
			map[string]interface{}{
				"foo": `{"hello":"world","bonjour":["monde","dieu"], "Guten Tag": 42, "你好": ["10", 20, 3.14]}`,
			},
			"foo",
			http.Header{
				"Hello":     []string{"world"},
				"Bonjour":   []string{"monde", "dieu"},
				"Guten Tag": []string{"42"},
				"你好":        []string{"10", "20", "3.14"},
			},
		},

		"type header, keypair string array with dupe key": {
			map[string]*FieldSchema{
				"foo": {Type: TypeHeader},
			},
			map[string]interface{}{
				"foo": []interface{}{"key1:value1", "key2:value2", "key3:1", "key3:true"},
			},
			"foo",
			http.Header{
				"Key1": []string{"value1"},
				"Key2": []string{"value2"},
				"Key3": []string{"1", "true"},
			},
		},

		"type header, map string slice": {
			map[string]*FieldSchema{
				"foo": {Type: TypeHeader},
			},
			map[string]interface{}{
				"foo": map[string][]string{
					"key1": {"value1"},
					"key2": {"value2"},
					"key3": {"1"},
				},
			},
			"foo",
			http.Header{
				"Key1": []string{"value1"},
				"Key2": []string{"value2"},
				"Key3": []string{"1"},
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

		"type header, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeHeader},
			},
			map[string]interface{}{},
			"foo",
			http.Header{},
		},
	}

	for name, tc := range cases {
		data := &FieldData{
			Raw:    tc.Raw,
			Schema: tc.Schema,
		}

		if err := data.Validate(); err != nil {
			t.Fatalf("bad: %s", err)
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

func TestFieldDataGetFirst(t *testing.T) {
	data := &FieldData{
		Raw: map[string]interface{}{
			"foo":  "bar",
			"fizz": "buzz",
		},
		Schema: map[string]*FieldSchema{
			"foo":  {Type: TypeNameString},
			"fizz": {Type: TypeNameString},
		},
	}

	result, ok := data.GetFirst("foo", "fizz")
	if !ok {
		t.Fatal("should have found value for foo")
	}
	if result.(string) != "bar" {
		t.Fatal("should have gotten bar for foo")
	}

	result, ok = data.GetFirst("fizz", "foo")
	if !ok {
		t.Fatal("should have found value for fizz")
	}
	if result.(string) != "buzz" {
		t.Fatal("should have gotten buzz for fizz")
	}

	result, ok = data.GetFirst("cats")
	if ok {
		t.Fatal("shouldn't have gotten anything for cats")
	}
}
