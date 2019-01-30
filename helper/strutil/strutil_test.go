package strutil

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"testing"
)

func TestStrUtil_StrListDelete(t *testing.T) {
	output := StrListDelete([]string{"item1", "item2", "item3"}, "item1")
	if StrListContains(output, "item1") {
		t.Fatal("bad: 'item1' should not have been present")
	}

	output = StrListDelete([]string{"item1", "item2", "item3"}, "item2")
	if StrListContains(output, "item2") {
		t.Fatal("bad: 'item2' should not have been present")
	}

	output = StrListDelete([]string{"item1", "item2", "item3"}, "item3")
	if StrListContains(output, "item3") {
		t.Fatal("bad: 'item3' should not have been present")
	}

	output = StrListDelete([]string{"item1", "item1", "item3"}, "item1")
	if !StrListContains(output, "item1") {
		t.Fatal("bad: 'item1' should have been present")
	}

	output = StrListDelete(output, "item1")
	if StrListContains(output, "item1") {
		t.Fatal("bad: 'item1' should not have been present")
	}

	output = StrListDelete(output, "random")
	if len(output) != 1 {
		t.Fatalf("bad: expected: 1, actual: %d", len(output))
	}

	output = StrListDelete(output, "item3")
	if StrListContains(output, "item3") {
		t.Fatal("bad: 'item3' should not have been present")
	}
}

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

func TestStrutil_ListContainsGlob(t *testing.T) {
	haystack := []string{
		"dev",
		"ops*",
		"root/*",
		"*-dev",
		"_*_",
	}
	if StrListContainsGlob(haystack, "tubez") {
		t.Fatalf("Value shouldn't exist")
	}
	if !StrListContainsGlob(haystack, "root/test") {
		t.Fatalf("Value should exist")
	}
	if !StrListContainsGlob(haystack, "ops_test") {
		t.Fatalf("Value should exist")
	}
	if !StrListContainsGlob(haystack, "ops") {
		t.Fatalf("Value should exist")
	}
	if !StrListContainsGlob(haystack, "dev") {
		t.Fatalf("Value should exist")
	}
	if !StrListContainsGlob(haystack, "test-dev") {
		t.Fatalf("Value should exist")
	}
	if !StrListContainsGlob(haystack, "_test_") {
		t.Fatalf("Value should exist")
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
		t.Fatalf("expected an error")
	}
	for k, _ := range actual {
		delete(actual, k)
	}

	input = "key1 = value1, 	=  value2 "
	err = ParseKeyValues(input, actual, ",")
	if err == nil {
		t.Fatalf("expected an error")
	}
	for k, _ := range actual {
		delete(actual, k)
	}

	input = "key1"
	err = ParseKeyValues(input, actual, ",")
	if err == nil {
		t.Fatalf("expected an error")
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

func TestGlobbedStringsMatch(t *testing.T) {
	type tCase struct {
		item   string
		val    string
		expect bool
	}

	tCases := []tCase{
		tCase{"", "", true},
		tCase{"*", "*", true},
		tCase{"**", "**", true},
		tCase{"*t", "t", true},
		tCase{"*t", "test", true},
		tCase{"t*", "test", true},
		tCase{"*test", "test", true},
		tCase{"*test", "a test", true},
		tCase{"test", "a test", false},
		tCase{"*test", "tests", false},
		tCase{"test*", "test", true},
		tCase{"test*", "testsss", true},
		tCase{"test**", "testsss", false},
		tCase{"test**", "test*", true},
		tCase{"**test", "*test", true},
		tCase{"TEST", "test", false},
		tCase{"test", "test", true},
	}

	for _, tc := range tCases {
		actual := GlobbedStringsMatch(tc.item, tc.val)

		if actual != tc.expect {
			t.Fatalf("Bad testcase %#v, expected %t, got %t", tc, tc.expect, actual)
		}
	}
}

func TestTrimStrings(t *testing.T) {
	input := []string{"abc", "123", "abcd ", "123  "}
	expected := []string{"abc", "123", "abcd", "123"}
	actual := TrimStrings(input)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Bad TrimStrings: expected:%#v, got:%#v", expected, actual)
	}
}

func TestRemoveEmpty(t *testing.T) {
	input := []string{"abc", "", "abc", ""}
	expected := []string{"abc", "abc"}
	actual := RemoveEmpty(input)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Bad TrimStrings: expected:%#v, got:%#v", expected, actual)
	}

	input = []string{""}
	expected = []string{}
	actual = RemoveEmpty(input)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Bad TrimStrings: expected:%#v, got:%#v", expected, actual)
	}
}

func TestStrutil_AppendIfMissing(t *testing.T) {
	keys := []string{}

	keys = AppendIfMissing(keys, "foo")

	if len(keys) != 1 {
		t.Fatalf("expected slice to be length of 1: %v", keys)
	}
	if keys[0] != "foo" {
		t.Fatalf("expected slice to contain key 'foo': %v", keys)
	}

	keys = AppendIfMissing(keys, "bar")

	if len(keys) != 2 {
		t.Fatalf("expected slice to be length of 2: %v", keys)
	}
	if keys[0] != "foo" {
		t.Fatalf("expected slice to contain key 'foo': %v", keys)
	}
	if keys[1] != "bar" {
		t.Fatalf("expected slice to contain key 'bar': %v", keys)
	}

	keys = AppendIfMissing(keys, "foo")

	if len(keys) != 2 {
		t.Fatalf("expected slice to still be length of 2: %v", keys)
	}
	if keys[0] != "foo" {
		t.Fatalf("expected slice to still contain key 'foo': %v", keys)
	}
	if keys[1] != "bar" {
		t.Fatalf("expected slice to still contain key 'bar': %v", keys)
	}
}

func TestStrUtil_RemoveDuplicates(t *testing.T) {
	type tCase struct {
		input     []string
		expect    []string
		lowercase bool
	}

	tCases := []tCase{
		tCase{[]string{}, []string{}, false},
		tCase{[]string{}, []string{}, true},
		tCase{[]string{"a", "b", "a"}, []string{"a", "b"}, false},
		tCase{[]string{"A", "b", "a"}, []string{"A", "a", "b"}, false},
		tCase{[]string{"A", "b", "a"}, []string{"a", "b"}, true},
	}

	for _, tc := range tCases {
		actual := RemoveDuplicates(tc.input, tc.lowercase)

		if !reflect.DeepEqual(actual, tc.expect) {
			t.Fatalf("Bad testcase %#v, expected %v, got %v", tc, tc.expect, actual)
		}
	}
}

func TestStrUtil_ParseStringSlice(t *testing.T) {
	type tCase struct {
		input  string
		sep    string
		expect []string
	}

	tCases := []tCase{
		tCase{"", "", []string{}},
		tCase{"   ", ",", []string{}},
		tCase{",   ", ",", []string{"", ""}},
		tCase{"a", ",", []string{"a"}},
		tCase{" a, b,   c   ", ",", []string{"a", "b", "c"}},
		tCase{" a; b;   c   ", ";", []string{"a", "b", "c"}},
		tCase{" a :: b  ::   c   ", "::", []string{"a", "b", "c"}},
	}

	for _, tc := range tCases {
		actual := ParseStringSlice(tc.input, tc.sep)

		if !reflect.DeepEqual(actual, tc.expect) {
			t.Fatalf("Bad testcase %#v, expected %v, got %v", tc, tc.expect, actual)
		}
	}
}

func TestStrUtil_MergeSlices(t *testing.T) {
	res := MergeSlices([]string{"a", "c", "d"}, []string{}, []string{"c", "f", "a"}, nil, []string{"foo"})

	expect := []string{"a", "c", "d", "f", "foo"}

	if !reflect.DeepEqual(res, expect) {
		t.Fatalf("expected %v, got %v", expect, res)
	}
}

func TestDifference(t *testing.T) {
	testCases := []struct {
		Name           string
		SetA           []string
		SetB           []string
		Lowercase      bool
		ExpectedResult []string
	}{
		{
			Name:           "case_sensitive",
			SetA:           []string{"a", "b", "c"},
			SetB:           []string{"b", "c"},
			Lowercase:      false,
			ExpectedResult: []string{"a"},
		},
		{
			Name:           "case_insensitive",
			SetA:           []string{"a", "B", "c"},
			SetB:           []string{"b", "C"},
			Lowercase:      true,
			ExpectedResult: []string{"a"},
		},
		{
			Name:           "no_match",
			SetA:           []string{"a", "b", "c"},
			SetB:           []string{"d"},
			Lowercase:      false,
			ExpectedResult: []string{"a", "b", "c"},
		},
		{
			Name:           "empty_set_a",
			SetA:           []string{},
			SetB:           []string{"d", "e"},
			Lowercase:      false,
			ExpectedResult: []string{},
		},
		{
			Name:           "empty_set_b",
			SetA:           []string{"a", "b"},
			SetB:           []string{},
			Lowercase:      false,
			ExpectedResult: []string{"a", "b"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actualResult := Difference(tc.SetA, tc.SetB, tc.Lowercase)

			if !reflect.DeepEqual(actualResult, tc.ExpectedResult) {
				t.Fatalf("expected %v, got %v", tc.ExpectedResult, actualResult)
			}
		})
	}
}

func TestStrUtil_EqualStringMaps(t *testing.T) {
	m1 := map[string]string{
		"foo": "a",
	}
	m2 := map[string]string{
		"foo": "a",
		"bar": "b",
	}
	var m3 map[string]string

	m4 := map[string]string{
		"dog": "",
	}

	m5 := map[string]string{
		"cat": "",
	}

	tests := []struct {
		a      map[string]string
		b      map[string]string
		result bool
	}{
		{m1, m1, true},
		{m2, m2, true},
		{m1, m2, false},
		{m2, m1, false},
		{m2, m2, true},
		{m3, m1, false},
		{m3, m3, true},
		{m4, m5, false},
	}

	for i, test := range tests {
		actual := EqualStringMaps(test.a, test.b)
		if actual != test.result {
			t.Fatalf("case %d, expected %v, got %v", i, test.result, actual)
		}
	}
}
