package framework

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/wrapping"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/version"
)

func TestOpenAPI_Regex(t *testing.T) {
	t.Run("Required", func(t *testing.T) {
		tests := []struct {
			input    string
			captures []string
		}{
			{`/foo/bar/(?P<val>.*)`, []string{"val"}},
			{`/foo/bar/` + GenericNameRegex("val"), []string{"val"}},
			{`/foo/bar/` + GenericNameRegex("first") + "/b/" + GenericNameRegex("second"), []string{"first", "second"}},
			{`/foo/bar`, []string{}},
		}

		for _, test := range tests {
			result := reqdRe.FindAllStringSubmatch(test.input, -1)
			if len(result) != len(test.captures) {
				t.Fatalf("Capture error (%s): expected %d matches, actual: %d", test.input, len(test.captures), len(result))
			}

			for i := 0; i < len(result); i++ {
				if result[i][1] != test.captures[i] {
					t.Fatalf("Capture error (%s): expected %s, actual: %s", test.input, test.captures[i], result[i][1])
				}
			}
		}
	})
	t.Run("Optional", func(t *testing.T) {
		input := "foo/(maybe/)?bar"
		expStart := len("foo/")
		expEnd := len(input) - len("bar")

		match := optRe.FindStringIndex(input)
		if diff := deep.Equal(match, []int{expStart, expEnd}); diff != nil {
			t.Fatal(diff)
		}

		input = "/foo/maybe/bar"
		match = optRe.FindStringIndex(input)
		if match != nil {
			t.Fatalf("Expected nil match (%s), got %+v", input, match)
		}
	})
	t.Run("Alternation", func(t *testing.T) {
		input := `(raw/?$|raw/(?P<path>.+))`

		matches := altRe.FindAllStringSubmatch(input, -1)
		exp1 := "raw/?$"
		exp2 := "raw/(?P<path>.+)"
		if matches[0][1] != exp1 || matches[0][2] != exp2 {
			t.Fatalf("Capture error. Expected %s and %s, got %v", exp1, exp2, matches[0][1:])
		}

		input = `/foo/bar/` + GenericNameRegex("val")

		matches = altRe.FindAllStringSubmatch(input, -1)
		if matches != nil {
			t.Fatalf("Expected nil match (%s), got %+v", input, matches)
		}
	})
	t.Run("Alternation Fields", func(t *testing.T) {
		input := `/foo/bar/(?P<type>auth|database|secret)/(?P<blah>a|b)`

		act := altFieldsGroupRe.ReplaceAllStringFunc(input, func(s string) string {
			return altFieldsRe.ReplaceAllString(s, ".+")
		})

		exp := "/foo/bar/(?P<type>.+)/(?P<blah>.+)"
		if act != exp {
			t.Fatalf("Replace error. Expected %s, got %v", exp, act)
		}
	})
	t.Run("Path fields", func(t *testing.T) {
		input := `/foo/bar/{inner}/baz/{outer}`

		matches := pathFieldsRe.FindAllStringSubmatch(input, -1)

		exp1 := "inner"
		exp2 := "outer"
		if matches[0][1] != exp1 || matches[1][1] != exp2 {
			t.Fatalf("Capture error. Expected %s and %s, got %v", exp1, exp2, matches)
		}

		input = `/foo/bar/inner/baz/outer`
		matches = pathFieldsRe.FindAllStringSubmatch(input, -1)

		if matches != nil {
			t.Fatalf("Expected nil match (%s), got %+v", input, matches)
		}
	})
	t.Run("Filtering", func(t *testing.T) {
		tests := []struct {
			input  string
			regex  *regexp.Regexp
			output string
		}{
			{
				input:  `ab?cde^fg(hi?j$k`,
				regex:  cleanCharsRe,
				output: "abcdefghijk",
			},
			{
				input:  `abcde/?`,
				regex:  cleanSuffixRe,
				output: "abcde",
			},
			{
				input:  `abcde/?$`,
				regex:  cleanSuffixRe,
				output: "abcde",
			},
			{
				input:  `abcde`,
				regex:  wsRe,
				output: "abcde",
			},
			{
				input:  `  a         b    cd   e   `,
				regex:  wsRe,
				output: "abcde",
			},
		}

		for _, test := range tests {
			result := test.regex.ReplaceAllString(test.input, "")
			if result != test.output {
				t.Fatalf("Clean Regex error (%s). Expected %s, got %s", test.input, test.output, result)
			}
		}

	})
}

func TestOpenAPI_ExpandPattern(t *testing.T) {
	tests := []struct {
		in_pattern   string
		out_pathlets []string
	}{
		{"rekey/backup", []string{"rekey/backup"}},
		{"rekey/backup$", []string{"rekey/backup"}},
		{"auth/(?P<path>.+?)/tune$", []string{"auth/{path}/tune"}},
		{"auth/(?P<path>.+?)/tune/(?P<more>.*?)$", []string{"auth/{path}/tune/{more}"}},
		{"tools/hash(/(?P<urlalgorithm>.+))?", []string{
			"tools/hash",
			"tools/hash/{urlalgorithm}",
		}},
		{"(leases/)?renew(/(?P<url_lease_id>.+))?", []string{
			"leases/renew",
			"leases/renew/{url_lease_id}",
			"renew",
			"renew/{url_lease_id}",
		}},
		{`config/ui/headers/` + GenericNameRegex("header"), []string{"config/ui/headers/{header}"}},
		{`leases/lookup/(?P<prefix>.+?)?`, []string{
			"leases/lookup/",
			"leases/lookup/{prefix}",
		}},
		{`(raw/?$|raw/(?P<path>.+))`, []string{
			"raw",
			"raw/{path}",
		}},
		{"lookup" + OptionalParamRegex("urltoken"), []string{
			"lookup",
			"lookup/{urltoken}",
		}},
		{"roles/?$", []string{
			"roles",
		}},
		{"roles/?", []string{
			"roles",
		}},
		{"accessors/$", []string{
			"accessors/",
		}},
		{"verify/" + GenericNameRegex("name") + OptionalParamRegex("urlalgorithm"), []string{
			"verify/{name}",
			"verify/{name}/{urlalgorithm}",
		}},
		{"^plugins/catalog/(?P<type>auth|database|secret)/(?P<name>.+)$", []string{
			"plugins/catalog/{type}/{name}",
		}},
		{"^plugins/catalog/(?P<type>auth|database|secret)/?$", []string{
			"plugins/catalog/{type}",
		}},
	}

	for i, test := range tests {
		out := expandPattern(test.in_pattern)
		sort.Strings(out)
		if !reflect.DeepEqual(out, test.out_pathlets) {
			t.Fatalf("Test %d: Expected %v got %v", i, test.out_pathlets, out)
		}
	}
}

func TestOpenAPI_SplitFields(t *testing.T) {
	fields := map[string]*FieldSchema{
		"a": {Description: "path"},
		"b": {Description: "body"},
		"c": {Description: "body"},
		"d": {Description: "body"},
		"e": {Description: "path"},
	}

	pathFields, bodyFields := splitFields(fields, "some/{a}/path/{e}")

	lp := len(pathFields)
	lb := len(bodyFields)
	l := len(fields)
	if lp+lb != l {
		t.Fatalf("split length error: %d + %d != %d", lp, lb, l)
	}

	for name, field := range pathFields {
		if field.Description != "path" {
			t.Fatalf("expected field %s to be in 'path', found in %s", name, field.Description)
		}
	}
	for name, field := range bodyFields {
		if field.Description != "body" {
			t.Fatalf("expected field %s to be in 'body', found in %s", name, field.Description)
		}
	}
}

func TestOpenAPI_SpecialPaths(t *testing.T) {
	tests := []struct {
		pattern     string
		rootPaths   []string
		root        bool
		unauthPaths []string
		unauth      bool
	}{
		{"foo", []string{}, false, []string{"foo"}, true},
		{"foo", []string{"foo"}, true, []string{"bar"}, false},
		{"foo/bar", []string{"foo"}, false, []string{"foo/*"}, true},
		{"foo/bar", []string{"foo/*"}, true, []string{"foo"}, false},
		{"foo/", []string{"foo/*"}, true, []string{"a", "b", "foo/"}, true},
		{"foo", []string{"foo*"}, true, []string{"a", "fo*"}, true},
		{"foo/bar", []string{"a", "b", "foo/*"}, true, []string{"foo/baz/*"}, false},
	}
	for i, test := range tests {
		doc := NewOASDocument()
		path := Path{
			Pattern: test.pattern,
		}
		sp := &logical.Paths{
			Root:            test.rootPaths,
			Unauthenticated: test.unauthPaths,
		}
		documentPath(&path, sp, logical.TypeLogical, doc)
		result := test.root
		if doc.Paths["/"+test.pattern].Sudo != result {
			t.Fatalf("Test (root) %d: Expected %v got %v", i, test.root, result)
		}
		result = test.unauth
		if doc.Paths["/"+test.pattern].Unauthenticated != result {
			t.Fatalf("Test (unauth) %d: Expected %v got %v", i, test.unauth, result)
		}
	}
}

func TestOpenAPI_Paths(t *testing.T) {
	origDepth := deep.MaxDepth
	defer func() { deep.MaxDepth = origDepth }()
	deep.MaxDepth = 20

	t.Run("Legacy callbacks", func(t *testing.T) {
		p := &Path{
			Pattern: "lookup/" + GenericNameRegex("id"),

			Fields: map[string]*FieldSchema{
				"id": &FieldSchema{
					Type:        TypeString,
					Description: "My id parameter",
				},
				"token": &FieldSchema{
					Type:        TypeString,
					Description: "My token",
				},
			},

			Callbacks: map[logical.Operation]OperationFunc{
				logical.ReadOperation:   nil,
				logical.UpdateOperation: nil,
			},

			HelpSynopsis:    "Synopsis",
			HelpDescription: "Description",
		}

		sp := &logical.Paths{
			Root:            []string{},
			Unauthenticated: []string{},
		}
		testPath(t, p, sp, expected("legacy"))
	})

	t.Run("Operations", func(t *testing.T) {
		p := &Path{
			Pattern: "foo/" + GenericNameRegex("id"),
			Fields: map[string]*FieldSchema{
				"id": {
					Type:        TypeString,
					Description: "id path parameter",
				},
				"flavors": {
					Type:        TypeCommaStringSlice,
					Description: "the flavors",
				},
				"name": {
					Type:        TypeNameString,
					Default:     "Larry",
					Description: "the name",
				},
				"age": {
					Type:             TypeInt,
					Description:      "the age",
					AllowedValues:    []interface{}{1, 2, 3},
					Required:         true,
					DisplayName:      "Age",
					DisplayValue:     7,
					DisplaySensitive: true,
				},
				"x-abc-token": {
					Type:          TypeHeader,
					Description:   "a header value",
					AllowedValues: []interface{}{"a", "b", "c"},
				},
				"format": {
					Type:        TypeString,
					Description: "a query param",
					Query:       true,
				},
			},
			HelpSynopsis:    "Synopsis",
			HelpDescription: "Description",
			Operations: map[logical.Operation]OperationHandler{
				logical.ReadOperation: &PathOperation{
					Summary:     "My Summary",
					Description: "My Description",
				},
				logical.UpdateOperation: &PathOperation{
					Summary:     "Update Summary",
					Description: "Update Description",
				},
				logical.CreateOperation: &PathOperation{
					Summary:     "Create Summary",
					Description: "Create Description",
				},
				logical.ListOperation: &PathOperation{
					Summary:     "List Summary",
					Description: "List Description",
				},
				logical.DeleteOperation: &PathOperation{
					Summary:     "This shouldn't show up",
					Unpublished: true,
				},
			},
		}

		sp := &logical.Paths{
			Root: []string{"foo*"},
		}
		testPath(t, p, sp, expected("operations"))
	})

	t.Run("Responses", func(t *testing.T) {
		p := &Path{
			Pattern:         "foo",
			HelpSynopsis:    "Synopsis",
			HelpDescription: "Description",
			Operations: map[logical.Operation]OperationHandler{
				logical.ReadOperation: &PathOperation{
					Summary:     "My Summary",
					Description: "My Description",
					Responses: map[int][]Response{
						202: {{
							Description: "Amazing",
							Example: &logical.Response{
								Data: map[string]interface{}{
									"amount": 42,
								},
							},
						}},
					},
				},
				logical.DeleteOperation: &PathOperation{
					Summary: "Delete stuff",
				},
			},
		}

		sp := &logical.Paths{
			Unauthenticated: []string{"x", "y", "foo"},
		}

		testPath(t, p, sp, expected("responses"))
	})
}

func TestOpenAPI_OperationID(t *testing.T) {
	path1 := &Path{
		Pattern: "foo/" + GenericNameRegex("id"),
		Fields: map[string]*FieldSchema{
			"id": {Type: TypeString},
		},
		Operations: map[logical.Operation]OperationHandler{
			logical.ReadOperation:   &PathOperation{},
			logical.UpdateOperation: &PathOperation{},
			logical.DeleteOperation: &PathOperation{},
		},
	}

	path2 := &Path{
		Pattern: "Foo/" + GenericNameRegex("id"),
		Fields: map[string]*FieldSchema{
			"id": {Type: TypeString},
		},
		Operations: map[logical.Operation]OperationHandler{
			logical.ReadOperation: &PathOperation{},
		},
	}

	for _, context := range []string{"", "bar"} {
		doc := NewOASDocument()
		documentPath(path1, nil, logical.TypeLogical, doc)
		documentPath(path2, nil, logical.TypeLogical, doc)
		doc.CreateOperationIDs(context)

		tests := []struct {
			path string
			op   string
			opID string
		}{
			{"/Foo/{id}", "get", "getFooId"},
			{"/foo/{id}", "get", "getFooId_2"},
			{"/foo/{id}", "post", "postFooId"},
			{"/foo/{id}", "delete", "deleteFooId"},
		}

		for _, test := range tests {
			actual := getPathOp(doc.Paths[test.path], test.op).OperationID
			expected := test.opID
			if context != "" {
				expected += "_" + context
			}

			if actual != expected {
				t.Fatalf("expected %v, got %v", expected, actual)
			}
		}
	}
}

func TestOpenAPI_CustomDecoder(t *testing.T) {
	p := &Path{
		Pattern:      "foo",
		HelpSynopsis: "Synopsis",
		Operations: map[logical.Operation]OperationHandler{
			logical.ReadOperation: &PathOperation{
				Summary: "My Summary",
				Responses: map[int][]Response{
					100: {{
						Description: "OK",
						Example: &logical.Response{
							Data: map[string]interface{}{
								"foo": 42,
							},
						},
					}},
					200: {{
						Description: "Good",
						Example:     (*logical.Response)(nil),
					}},
					599: {{
						Description: "Bad",
					}},
				},
			},
		},
	}

	docOrig := NewOASDocument()
	documentPath(p, nil, logical.TypeLogical, docOrig)

	docJSON := mustJSONMarshal(t, docOrig)

	var intermediate map[string]interface{}
	if err := jsonutil.DecodeJSON(docJSON, &intermediate); err != nil {
		t.Fatal(err)
	}

	docNew, err := NewOASDocumentFromMap(intermediate)
	if err != nil {
		t.Fatal(err)
	}

	docNewJSON := mustJSONMarshal(t, docNew)

	if diff := deep.Equal(docJSON, docNewJSON); diff != nil {
		t.Fatal(diff)
	}
}

func TestOpenAPI_CleanResponse(t *testing.T) {
	// Verify that an all-null input results in empty JSON
	orig := &logical.Response{}

	cr := cleanResponse(orig)

	newJSON := mustJSONMarshal(t, cr)

	if !bytes.Equal(newJSON, []byte("{}")) {
		t.Fatalf("expected {}, got: %q", newJSON)
	}

	// Verify that all non-null inputs results in JSON that matches the marshalling of
	// logical.Response. This will fail if logical.Response changes without a corresponding
	// change to cleanResponse()
	orig = &logical.Response{
		Secret:   new(logical.Secret),
		Auth:     new(logical.Auth),
		Data:     map[string]interface{}{"foo": 42},
		Redirect: "foo",
		Warnings: []string{"foo"},
		WrapInfo: &wrapping.ResponseWrapInfo{Token: "foo"},
		Headers:  map[string][]string{"foo": {"bar"}},
	}
	origJSON := mustJSONMarshal(t, orig)

	cr = cleanResponse(orig)

	cleanJSON := mustJSONMarshal(t, cr)

	if diff := deep.Equal(origJSON, cleanJSON); diff != nil {
		t.Fatal(diff)
	}
}

func testPath(t *testing.T, path *Path, sp *logical.Paths, expectedJSON string) {
	t.Helper()

	doc := NewOASDocument()
	if err := documentPath(path, sp, logical.TypeLogical, doc); err != nil {
		t.Fatal(err)
	}
	doc.CreateOperationIDs("")

	docJSON, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	// Compare json by first decoding, then comparing with a deep equality check.
	var expected, actual interface{}
	if err := jsonutil.DecodeJSON(docJSON, &actual); err != nil {
		t.Fatal(err)
	}

	if err := jsonutil.DecodeJSON([]byte(expectedJSON), &expected); err != nil {
		t.Fatal(err)
	}

	if diff := deep.Equal(actual, expected); diff != nil {
		//fmt.Println(string(docJSON)) // uncomment to debug generated JSON (very helpful when fixing tests)
		t.Fatal(diff)
	}
}

func getPathOp(pi *OASPathItem, op string) *OASOperation {
	switch op {
	case "get":
		return pi.Get
	case "post":
		return pi.Post
	case "delete":
		return pi.Delete
	default:
		panic("unexpected operation: " + op)
	}
}

func expected(name string) string {
	data, err := ioutil.ReadFile(filepath.Join("testdata", name+".json"))
	if err != nil {
		panic(err)
	}

	content := strings.Replace(string(data), "<vault_version>", version.GetVersion().Version, 1)

	return content
}

func mustJSONMarshal(t *testing.T, data interface{}) []byte {
	j, err := json.MarshalIndent(data, "", "  ")

	if err != nil {
		t.Fatal(err)
	}
	return j
}
