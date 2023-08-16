package framework

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestFieldDataGet(t *testing.T) {
	cases := map[string]struct {
		Schema      map[string]*FieldSchema
		Raw         map[string]interface{}
		Key         string
		Value       interface{}
		ExpectError bool
	}{
		"string type, string value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeString},
			},
			map[string]interface{}{
				"foo": "bar",
			},
			"foo",
			"bar",
			false,
		},

		"string type, int value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeString},
			},
			map[string]interface{}{
				"foo": 42,
			},
			"foo",
			"42",
			false,
		},

		"string type, unset value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeString},
			},
			map[string]interface{}{},
			"foo",
			"",
			false,
		},

		"string type, unset value with default": {
			map[string]*FieldSchema{
				"foo": {
					Type:    TypeString,
					Default: "bar",
				},
			},
			map[string]interface{}{},
			"foo",
			"bar",
			false,
		},

		"lowercase string type, lowercase string value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeLowerCaseString},
			},
			map[string]interface{}{
				"foo": "bar",
			},
			"foo",
			"bar",
			false,
		},

		"lowercase string type, mixed-case string value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeLowerCaseString},
			},
			map[string]interface{}{
				"foo": "BaR",
			},
			"foo",
			"bar",
			false,
		},

		"lowercase string type, int value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeLowerCaseString},
			},
			map[string]interface{}{
				"foo": 42,
			},
			"foo",
			"42",
			false,
		},

		"lowercase string type, unset value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeLowerCaseString},
			},
			map[string]interface{}{},
			"foo",
			"",
			false,
		},

		"lowercase string type, unset value with lowercase default": {
			map[string]*FieldSchema{
				"foo": {
					Type:    TypeLowerCaseString,
					Default: "bar",
				},
			},
			map[string]interface{}{},
			"foo",
			"bar",
			false,
		},

		"int type, int value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeInt},
			},
			map[string]interface{}{
				"foo": 42,
			},
			"foo",
			42,
			false,
		},

		"bool type, bool value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeBool},
			},
			map[string]interface{}{
				"foo": false,
			},
			"foo",
			false,
			false,
		},

		"map type, map value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeMap},
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
			false,
		},

		"duration type, string value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeDurationSecond},
			},
			map[string]interface{}{
				"foo": "42",
			},
			"foo",
			42,
			false,
		},

		"duration type, string duration value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeDurationSecond},
			},
			map[string]interface{}{
				"foo": "42m",
			},
			"foo",
			2520,
			false,
		},

		"duration type, int value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeDurationSecond},
			},
			map[string]interface{}{
				"foo": 42,
			},
			"foo",
			42,
			false,
		},

		"duration type, float value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeDurationSecond},
			},
			map[string]interface{}{
				"foo": 42.0,
			},
			"foo",
			42,
			false,
		},

		"duration type, nil value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeDurationSecond},
			},
			map[string]interface{}{
				"foo": nil,
			},
			"foo",
			0,
			false,
		},

		"duration type, 0 value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeDurationSecond},
			},
			map[string]interface{}{
				"foo": 0,
			},
			"foo",
			0,
			false,
		},

		"signed duration type, positive string value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeSignedDurationSecond},
			},
			map[string]interface{}{
				"foo": "42",
			},
			"foo",
			42,
			false,
		},

		"signed duration type, positive string duration value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeSignedDurationSecond},
			},
			map[string]interface{}{
				"foo": "42m",
			},
			"foo",
			2520,
			false,
		},

		"signed duration type, positive int value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeSignedDurationSecond},
			},
			map[string]interface{}{
				"foo": 42,
			},
			"foo",
			42,
			false,
		},

		"signed duration type, positive float value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeSignedDurationSecond},
			},
			map[string]interface{}{
				"foo": 42.0,
			},
			"foo",
			42,
			false,
		},

		"signed duration type, negative string value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeSignedDurationSecond},
			},
			map[string]interface{}{
				"foo": "-42",
			},
			"foo",
			-42,
			false,
		},

		"signed duration type, negative string duration value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeSignedDurationSecond},
			},
			map[string]interface{}{
				"foo": "-42m",
			},
			"foo",
			-2520,
			false,
		},

		"signed duration type, negative int value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeSignedDurationSecond},
			},
			map[string]interface{}{
				"foo": -42,
			},
			"foo",
			-42,
			false,
		},

		"signed duration type, negative float value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeSignedDurationSecond},
			},
			map[string]interface{}{
				"foo": -42.0,
			},
			"foo",
			-42,
			false,
		},

		"signed duration type, nil value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeSignedDurationSecond},
			},
			map[string]interface{}{
				"foo": nil,
			},
			"foo",
			0,
			false,
		},

		"signed duration type, 0 value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeSignedDurationSecond},
			},
			map[string]interface{}{
				"foo": 0,
			},
			"foo",
			0,
			false,
		},

		"slice type, empty slice": {
			map[string]*FieldSchema{
				"foo": {Type: TypeSlice},
			},
			map[string]interface{}{
				"foo": []interface{}{},
			},
			"foo",
			[]interface{}{},
			false,
		},

		"slice type, filled, mixed slice": {
			map[string]*FieldSchema{
				"foo": {Type: TypeSlice},
			},
			map[string]interface{}{
				"foo": []interface{}{123, "abc"},
			},
			"foo",
			[]interface{}{123, "abc"},
			false,
		},

		"string slice type, filled slice": {
			map[string]*FieldSchema{
				"foo": {Type: TypeStringSlice},
			},
			map[string]interface{}{
				"foo": []interface{}{123, "abc"},
			},
			"foo",
			[]string{"123", "abc"},
			false,
		},

		"string slice type, single value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeStringSlice},
			},
			map[string]interface{}{
				"foo": "abc",
			},
			"foo",
			[]string{"abc"},
			false,
		},

		"string slice type, empty string": {
			map[string]*FieldSchema{
				"foo": {Type: TypeStringSlice},
			},
			map[string]interface{}{
				"foo": "",
			},
			"foo",
			[]string{},
			false,
		},

		"comma string slice type, empty string": {
			map[string]*FieldSchema{
				"foo": {Type: TypeCommaStringSlice},
			},
			map[string]interface{}{
				"foo": "",
			},
			"foo",
			[]string{},
			false,
		},

		"comma string slice type, comma string with one value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeCommaStringSlice},
			},
			map[string]interface{}{
				"foo": "value1",
			},
			"foo",
			[]string{"value1"},
			false,
		},

		"comma string slice type, comma string with multi value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeCommaStringSlice},
			},
			map[string]interface{}{
				"foo": "value1,value2,value3",
			},
			"foo",
			[]string{"value1", "value2", "value3"},
			false,
		},

		"comma string slice type, nil string slice value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeCommaStringSlice},
			},
			map[string]interface{}{
				"foo": "",
			},
			"foo",
			[]string{},
			false,
		},

		"comma string slice type, string slice with one value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeCommaStringSlice},
			},
			map[string]interface{}{
				"foo": []interface{}{"value1"},
			},
			"foo",
			[]string{"value1"},
			false,
		},

		"comma string slice type, string slice with multi value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeCommaStringSlice},
			},
			map[string]interface{}{
				"foo": []interface{}{"value1", "value2", "value3"},
			},
			"foo",
			[]string{"value1", "value2", "value3"},
			false,
		},

		"comma string slice type, empty string slice value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeCommaStringSlice},
			},
			map[string]interface{}{
				"foo": []interface{}{},
			},
			"foo",
			[]string{},
			false,
		},

		"comma int slice type, comma int with one value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeCommaIntSlice},
			},
			map[string]interface{}{
				"foo": 1,
			},
			"foo",
			[]int{1},
			false,
		},

		"comma int slice type, comma int with multi value slice": {
			map[string]*FieldSchema{
				"foo": {Type: TypeCommaIntSlice},
			},
			map[string]interface{}{
				"foo": []int{1, 2, 3},
			},
			"foo",
			[]int{1, 2, 3},
			false,
		},

		"comma int slice type, comma int with multi value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeCommaIntSlice},
			},
			map[string]interface{}{
				"foo": "1,2,3",
			},
			"foo",
			[]int{1, 2, 3},
			false,
		},

		"comma int slice type, nil int slice value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeCommaIntSlice},
			},
			map[string]interface{}{
				"foo": "",
			},
			"foo",
			[]int{},
			false,
		},

		"comma int slice type, int slice with one value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeCommaIntSlice},
			},
			map[string]interface{}{
				"foo": []interface{}{"1"},
			},
			"foo",
			[]int{1},
			false,
		},

		"comma int slice type, int slice with multi value strings": {
			map[string]*FieldSchema{
				"foo": {Type: TypeCommaIntSlice},
			},
			map[string]interface{}{
				"foo": []interface{}{"1", "2", "3"},
			},
			"foo",
			[]int{1, 2, 3},
			false,
		},

		"comma int slice type, int slice with multi value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeCommaIntSlice},
			},
			map[string]interface{}{
				"foo": []interface{}{1, 2, 3},
			},
			"foo",
			[]int{1, 2, 3},
			false,
		},

		"comma int slice type, empty int slice value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeCommaIntSlice},
			},
			map[string]interface{}{
				"foo": []interface{}{},
			},
			"foo",
			[]int{},
			false,
		},

		"comma int slice type, json number": {
			map[string]*FieldSchema{
				"foo": {Type: TypeCommaIntSlice},
			},
			map[string]interface{}{
				"foo": json.Number("1"),
			},
			"foo",
			[]int{1},
			false,
		},

		"name string type, valid string": {
			map[string]*FieldSchema{
				"foo": {Type: TypeNameString},
			},
			map[string]interface{}{
				"foo": "bar",
			},
			"foo",
			"bar",
			false,
		},

		"name string type, valid value with special characters": {
			map[string]*FieldSchema{
				"foo": {Type: TypeNameString},
			},
			map[string]interface{}{
				"foo": "bar.baz-bay123",
			},
			"foo",
			"bar.baz-bay123",
			false,
		},

		"keypair type, valid value map type": {
			map[string]*FieldSchema{
				"foo": {Type: TypeKVPairs},
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
			false,
		},

		"keypair type, list of equal sign delim key pairs type": {
			map[string]*FieldSchema{
				"foo": {Type: TypeKVPairs},
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
			false,
		},

		"keypair type, single equal sign delim value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeKVPairs},
			},
			map[string]interface{}{
				"foo": "key1=value1",
			},
			"foo",
			map[string]string{
				"key1": "value1",
			},
			false,
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
			false,
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
			false,
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
			false,
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
			false,
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
			false,
		},

		"name string type, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeNameString},
			},
			map[string]interface{}{},
			"foo",
			"",
			false,
		},

		"string type, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeString},
			},
			map[string]interface{}{},
			"foo",
			"",
			false,
		},

		"type int, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeInt},
			},
			map[string]interface{}{},
			"foo",
			0,
			false,
		},

		"type bool, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeBool},
			},
			map[string]interface{}{},
			"foo",
			false,
			false,
		},

		"type map, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeMap},
			},
			map[string]interface{}{},
			"foo",
			map[string]interface{}{},
			false,
		},

		"type duration second, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeDurationSecond},
			},
			map[string]interface{}{},
			"foo",
			0,
			false,
		},

		"type signed duration second, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeSignedDurationSecond},
			},
			map[string]interface{}{},
			"foo",
			0,
			false,
		},

		"type slice, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeSlice},
			},
			map[string]interface{}{},
			"foo",
			[]interface{}{},
			false,
		},

		"type string slice, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeStringSlice},
			},
			map[string]interface{}{},
			"foo",
			[]string{},
			false,
		},

		"type comma string slice, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeCommaStringSlice},
			},
			map[string]interface{}{},
			"foo",
			[]string{},
			false,
		},

		"comma string slice type, single JSON number value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeCommaStringSlice},
			},
			map[string]interface{}{
				"foo": json.Number("123"),
			},
			"foo",
			[]string{"123"},
			false,
		},

		"type kv pair, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeKVPairs},
			},
			map[string]interface{}{},
			"foo",
			map[string]string{},
			false,
		},

		"type header, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeHeader},
			},
			map[string]interface{}{},
			"foo",
			http.Header{},
			false,
		},

		"float type, positive with decimals, as string": {
			map[string]*FieldSchema{
				"foo": {Type: TypeFloat},
			},
			map[string]interface{}{
				"foo": "1234567.891234567",
			},
			"foo",
			1234567.891234567,
			false,
		},

		"float type, negative with decimals, as string": {
			map[string]*FieldSchema{
				"foo": {Type: TypeFloat},
			},
			map[string]interface{}{
				"foo": "-1234567.891234567",
			},
			"foo",
			-1234567.891234567,
			false,
		},

		"float type, positive without decimals": {
			map[string]*FieldSchema{
				"foo": {Type: TypeFloat},
			},
			map[string]interface{}{
				"foo": 1234567,
			},
			"foo",
			1234567.0,
			false,
		},

		"type float, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeFloat},
			},
			map[string]interface{}{},
			"foo",
			0.0,
			false,
		},

		"type float, invalid value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeFloat},
			},
			map[string]interface{}{
				"foo": "invalid0.0",
			},
			"foo",
			0.0,
			true,
		},

		"type time, not supplied": {
			map[string]*FieldSchema{
				"foo": {Type: TypeTime},
			},
			map[string]interface{}{},
			"foo",
			time.Time{},
			false,
		},
		"type time, string value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeTime},
			},
			map[string]interface{}{
				"foo": "2021-12-11T09:08:07Z",
			},
			"foo",
			// Comparison uses DeepEqual() so better match exactly,
			// can't have a different location.
			time.Date(2021, 12, 11, 9, 8, 7, 0, time.UTC),
			false,
		},
		"type time, invalid value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeTime},
			},
			map[string]interface{}{
				"foo": "2021-13-11T09:08:07+02:00",
			},
			"foo",
			time.Time{},
			true,
		},
	}

	for name, tc := range cases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			data := &FieldData{
				Raw:    tc.Raw,
				Schema: tc.Schema,
			}

			err := data.Validate()
			switch {
			case tc.ExpectError && err == nil:
				t.Fatalf("expected error")
			case tc.ExpectError && err != nil:
				return
			case !tc.ExpectError && err != nil:
				t.Fatal(err)
			default:
				// Continue if !tc.ExpectError && err == nil
			}

			actual := data.Get(tc.Key)
			if !reflect.DeepEqual(actual, tc.Value) {
				t.Fatalf("Expected: %#v\nGot: %#v", tc.Value, actual)
			}
		})
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
				"foo": {Type: TypeNameString},
			},
			map[string]interface{}{
				"foo": "bar baz",
			},
			"foo",
		},
		"name string type, invalid value with special characters at beginning": {
			map[string]*FieldSchema{
				"foo": {Type: TypeNameString},
			},
			map[string]interface{}{
				"foo": ".barbaz",
			},
			"foo",
		},
		"name string type, invalid value with special characters at end": {
			map[string]*FieldSchema{
				"foo": {Type: TypeNameString},
			},
			map[string]interface{}{
				"foo": "barbaz-",
			},
			"foo",
		},
		"name string type, empty string": {
			map[string]*FieldSchema{
				"foo": {Type: TypeNameString},
			},
			map[string]interface{}{
				"foo": "",
			},
			"foo",
		},
		"keypair type, csv version empty key name": {
			map[string]*FieldSchema{
				"foo": {Type: TypeKVPairs},
			},
			map[string]interface{}{
				"foo": []interface{}{"=value1", "key2=value2", "key3=1"},
			},
			"foo",
		},
		"duration type, negative string value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeDurationSecond},
			},
			map[string]interface{}{
				"foo": "-42",
			},
			"foo",
		},
		"duration type, negative string duration value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeDurationSecond},
			},
			map[string]interface{}{
				"foo": "-42m",
			},
			"foo",
		},
		"duration type, negative int value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeDurationSecond},
			},
			map[string]interface{}{
				"foo": -42,
			},
			"foo",
		},
		"duration type, negative float value": {
			map[string]*FieldSchema{
				"foo": {Type: TypeDurationSecond},
			},
			map[string]interface{}{
				"foo": -42.0,
			},
			"foo",
		},
	}

	for name, tc := range cases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			data := &FieldData{
				Raw:    tc.Raw,
				Schema: tc.Schema,
			}

			got, _, err := data.GetOkErr(tc.Key)
			if err == nil {
				t.Fatalf("error expected, none received, got result: %#v", got)
			}
		})
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
