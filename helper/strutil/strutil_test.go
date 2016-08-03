package strutil

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"testing"
)

func TestStrutil_EquivalentSlices(t *testing.T) {
	slice1 := []string{"test2", "test1", "test3"}
	slice2 := []string{"test3", "test2", "test1"}
	if !EquivalentSlices(slice1, slice2) {
		t.Fatalf("bad: expected a match")
	}

	slice2 = append(slice2, "test4")
	if EquivalentSlices(slice1, slice2) {
		t.Fatalf("bad: expected a mismatch")
	}
}

func TestStrutil_ListContains(t *testing.T) {
	haystack := []string{
		"dev",
		"ops",
		"prod",
		"root",
	}
	if StrListContains(haystack, "tubez") {
		t.Fatalf("Bad")
	}
	if !StrListContains(haystack, "root") {
		t.Fatalf("Bad")
	}
}

func TestStrutil_ListSubset(t *testing.T) {
	parent := []string{
		"dev",
		"ops",
		"prod",
		"root",
	}
	child := []string{
		"prod",
		"ops",
	}
	if !StrListSubset(parent, child) {
		t.Fatalf("Bad")
	}
	if !StrListSubset(parent, parent) {
		t.Fatalf("Bad")
	}
	if !StrListSubset(child, child) {
		t.Fatalf("Bad")
	}
	if !StrListSubset(child, nil) {
		t.Fatalf("Bad")
	}
	if StrListSubset(child, parent) {
		t.Fatalf("Bad")
	}
	if StrListSubset(nil, child) {
		t.Fatalf("Bad")
	}
}

func TestStrutil_ParseKeyValues(t *testing.T) {
	actual := make(map[string]string)
	expected := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	var input string
	var err error

	input = "key1=value1,key2=value2"
	err = ParseKeyValues(input, actual, ",")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: expected: %#v\nactual: %#v", expected, actual)
	}
	for k, _ := range actual {
		delete(actual, k)
	}

	input = "key1 = value1, key2	= value2"
	err = ParseKeyValues(input, actual, ",")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: expected: %#v\nactual: %#v", expected, actual)
	}
	for k, _ := range actual {
		delete(actual, k)
	}

	input = "key1 = value1, key2	=   "
	err = ParseKeyValues(input, actual, ",")
	if err == nil {
		t.Fatal("expected an error")
	}
	for k, _ := range actual {
		delete(actual, k)
	}

	input = "key1 = value1, 	=  value2 "
	err = ParseKeyValues(input, actual, ",")
	if err == nil {
		t.Fatal("expected an error")
	}
	for k, _ := range actual {
		delete(actual, k)
	}
}

func TestStrutil_ParseArbitraryKeyValues(t *testing.T) {
	actual := make(map[string]string)
	expected := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	var input string
	var err error

	// Test <key>=<value> as comma separated string
	input = "key1=value1,key2=value2"
	err = ParseArbitraryKeyValues(input, actual, ",")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: expected: %#v\nactual: %#v", expected, actual)
	}
	for k, _ := range actual {
		delete(actual, k)
	}

	// Test <key>=<value> as base64 encoded comma separated string
	input = base64.StdEncoding.EncodeToString([]byte(input))
	err = ParseArbitraryKeyValues(input, actual, ",")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: expected: %#v\nactual: %#v", expected, actual)
	}
	for k, _ := range actual {
		delete(actual, k)
	}

	// Test JSON encoded <key>=<value> tuples
	input = `{"key1":"value1", "key2":"value2"}`
	err = ParseArbitraryKeyValues(input, actual, ",")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: expected: %#v\nactual: %#v", expected, actual)
	}
	for k, _ := range actual {
		delete(actual, k)
	}

	// Test base64 encoded JSON string of <key>=<value> tuples
	input = base64.StdEncoding.EncodeToString([]byte(input))
	err = ParseArbitraryKeyValues(input, actual, ",")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: expected: %#v\nactual: %#v", expected, actual)
	}
	for k, _ := range actual {
		delete(actual, k)
	}
}

func TestStrutil_ParseArbitraryStringSlice(t *testing.T) {
	input := `CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}';GRANT "foo-role" TO "{{name}}";ALTER ROLE "{{name}}" SET search_path = foo;GRANT CONNECT ON DATABASE "postgres" TO "{{name}}";`

	jsonExpected := []string{
		`DO $$
BEGIN
   IF NOT EXISTS (SELECT * FROM pg_catalog.pg_roles WHERE rolname='foo-role') THEN
      CREATE ROLE "foo-role";
      CREATE SCHEMA IF NOT EXISTS foo AUTHORIZATION "foo-role";
      ALTER ROLE "foo-role" SET search_path = foo;
      GRANT TEMPORARY ON DATABASE "postgres" TO "foo-role";
      GRANT ALL PRIVILEGES ON SCHEMA foo TO "foo-role";
      GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA foo TO "foo-role";
      GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA foo TO "foo-role";
      GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA foo TO "foo-role";
   END IF;
END
$$`,
		`CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'`,
		`GRANT "foo-role" TO "{{name}}"`,
		`ALTER ROLE "{{name}}" SET search_path = foo`,
		`GRANT CONNECT ON DATABASE "postgres" TO "{{name}}"`,
		``,
	}

	nonJSONExpected := jsonExpected[1:]

	var actual []string
	var inputB64 string
	var err error

	// Test non-JSON string
	actual = ParseArbitraryStringSlice(input, ";")
	if !reflect.DeepEqual(nonJSONExpected, actual) {
		t.Fatalf("bad: expected:\n%#v\nactual:\n%#v", nonJSONExpected, actual)
	}

	// Test base64-encoded non-JSON string
	inputB64 = base64.StdEncoding.EncodeToString([]byte(input))
	actual = ParseArbitraryStringSlice(inputB64, ";")
	if !reflect.DeepEqual(nonJSONExpected, actual) {
		t.Fatalf("bad: expected:\n%#v\nactual:\n%#v", nonJSONExpected, actual)
	}

	// Test JSON encoded
	inputJSON, err := json.Marshal(jsonExpected)
	if err != nil {
		t.Fatal(err)
	}

	actual = ParseArbitraryStringSlice(string(inputJSON), ";")
	if !reflect.DeepEqual(jsonExpected, actual) {
		t.Fatalf("bad: expected:\n%#v\nactual:\n%#v", string(inputJSON), actual)
	}

	// Test base64 encoded JSON string of <key>=<value> tuples
	inputB64 = base64.StdEncoding.EncodeToString(inputJSON)
	actual = ParseArbitraryStringSlice(inputB64, ";")
	if !reflect.DeepEqual(jsonExpected, actual) {
		t.Fatalf("bad: expected:\n%#v\nactual:\n%#v", jsonExpected, actual)
	}
}
