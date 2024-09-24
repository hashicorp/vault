// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/wrapping"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestOpenAPI_Regex(t *testing.T) {
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
		inPattern   string
		outPathlets []string
	}{
		// A simple string without regexp metacharacters passes through as is
		{"rekey/backup", []string{"rekey/backup"}},
		// A trailing regexp anchor metacharacter is removed
		{"rekey/backup$", []string{"rekey/backup"}},
		// As is a leading one
		{"^rekey/backup", []string{"rekey/backup"}},
		// Named capture groups become OpenAPI parameters
		{"auth/(?P<path>.+?)/tune$", []string{"auth/{path}/tune"}},
		{"auth/(?P<path>.+?)/tune/(?P<more>.*?)$", []string{"auth/{path}/tune/{more}"}},
		// Even if the capture group contains very complex regexp structure inside it
		{"something/(?P<something>(a|b(c|d))|e+|f{1,3}[ghi-k]?.*)", []string{"something/{something}"}},
		// A question-mark results in a result without and with the optional path part
		{"tools/hash(/(?P<urlalgorithm>.+))?", []string{
			"tools/hash",
			"tools/hash/{urlalgorithm}",
		}},
		// Multiple question-marks evaluate each possible combination
		{"(leases/)?renew(/(?P<url_lease_id>.+))?", []string{
			"leases/renew",
			"leases/renew/{url_lease_id}",
			"renew",
			"renew/{url_lease_id}",
		}},
		// GenericNameRegex is one particular way of writing a named capture group, so behaves the same
		{`config/ui/headers/` + GenericNameRegex("header"), []string{"config/ui/headers/{header}"}},
		// The question-mark behaviour is still works when the question-mark is directly applied to a named capture group
		{`leases/lookup/(?P<prefix>.+?)?`, []string{
			"leases/lookup/",
			"leases/lookup/{prefix}",
		}},
		// Optional trailing slashes at the end of the path get stripped - even if appearing deep inside an alternation
		{`(raw/?$|raw/(?P<path>.+))`, []string{
			"raw",
			"raw/{path}",
		}},
		// OptionalParamRegex is also another way of writing a named capture group, that is optional
		{"lookup" + OptionalParamRegex("urltoken"), []string{
			"lookup",
			"lookup/{urltoken}",
		}},
		// Optional trailign slashes at the end of the path get stripped in simpler cases too
		{"roles/?$", []string{
			"roles",
		}},
		{"roles/?", []string{
			"roles",
		}},
		// Non-optional trailing slashes remain... although don't do this, it breaks HelpOperation!
		// (Existing real examples of this pattern being fixed via https://github.com/hashicorp/vault/pull/18571)
		{"accessors/$", []string{
			"accessors/",
		}},
		// GenericNameRegex and OptionalParamRegex still work when concatenated
		{"verify/" + GenericNameRegex("name") + OptionalParamRegex("urlalgorithm"), []string{
			"verify/{name}",
			"verify/{name}/{urlalgorithm}",
		}},
		// Named capture groups that specify enum-like parameters work as expected
		{"^plugins/catalog/(?P<type>auth|database|secret)/(?P<name>.+)$", []string{
			"plugins/catalog/{type}/{name}",
		}},
		{"^plugins/catalog/(?P<type>auth|database|secret)/?$", []string{
			"plugins/catalog/{type}",
		}},
		// Alternations between various literal path segments work
		{"(pathOne|pathTwo)/", []string{"pathOne/", "pathTwo/"}},
		{"(pathOne|pathTwo)/" + GenericNameRegex("name"), []string{"pathOne/{name}", "pathTwo/{name}"}},
		{
			"(pathOne|path-2|Path_3)/" + GenericNameRegex("name"),
			[]string{"Path_3/{name}", "path-2/{name}", "pathOne/{name}"},
		},
		// They still work when combined with GenericNameWithAtRegex
		{"(creds|sts)/" + GenericNameWithAtRegex("name"), []string{
			"creds/{name}",
			"sts/{name}",
		}},
		// And when they're somewhere other than the start of the pattern
		{"keys/generate/(internal|exported|kms)", []string{
			"keys/generate/exported",
			"keys/generate/internal",
			"keys/generate/kms",
		}},
		// If a plugin author makes their list operation support both singular and plural forms, the OpenAPI notices
		{"rolesets?/?", []string{"roleset", "rolesets"}},
		// Complex nested alternation and question-marks are correctly interpreted
		{"crl(/pem|/delta(/pem)?)?", []string{"crl", "crl/delta", "crl/delta/pem", "crl/pem"}},
	}

	for i, test := range tests {
		paths, _, err := expandPattern(test.inPattern)
		if err != nil {
			t.Fatal(err)
		}
		sort.Strings(paths)
		if !reflect.DeepEqual(paths, test.outPathlets) {
			t.Fatalf("Test %d: Expected %v got %v", i, test.outPathlets, paths)
		}
	}
}

func TestOpenAPI_ExpandPattern_ReturnsError(t *testing.T) {
	tests := []struct {
		inPattern string
		outError  error
	}{
		// None of these regexp constructs are allowed outside of named capture groups
		{"[a-z]", errUnsupportableRegexpOperationForOpenAPI},
		{".", errUnsupportableRegexpOperationForOpenAPI},
		{"a+", errUnsupportableRegexpOperationForOpenAPI},
		{"a*", errUnsupportableRegexpOperationForOpenAPI},
		// So this pattern, which is a combination of two of the above isn't either - this pattern occurs in the KV
		// secrets engine for its catch-all error handler, which provides a helpful hint to people treating a KV v2 as
		// a KV v1.
		{".*", errUnsupportableRegexpOperationForOpenAPI},
	}

	for i, test := range tests {
		_, _, err := expandPattern(test.inPattern)
		if err != test.outError {
			t.Fatalf("Test %d: Expected %q got %q", i, test.outError, err)
		}
	}
}

func TestOpenAPI_SplitFields(t *testing.T) {
	paths, captures, err := expandPattern("some/" + GenericNameRegex("a") + "/path" + OptionalParamRegex("e"))
	if err != nil {
		t.Fatal(err)
	}

	fields := map[string]*FieldSchema{
		"a": {Description: "path"},
		"b": {Description: "body"},
		"c": {Description: "body"},
		"d": {Description: "body"},
		"e": {Description: "path"},
		"f": {Description: "query", Query: true},
	}

	for index, path := range paths {
		pathFields, queryFields, bodyFields := splitFields(fields, path, captures)

		numPath := len(pathFields)
		numQuery := len(queryFields)
		numBody := len(bodyFields)
		numExpectedDiscarded := 0
		// The first path generated is expected to be the one omitting the optional parameter field "e"
		if index == 0 {
			numExpectedDiscarded = 1
		}
		l := len(fields)
		if numPath+numQuery+numBody+numExpectedDiscarded != l {
			t.Fatalf("split length error: %d + %d + %d + %d != %d", numPath, numQuery, numBody, numExpectedDiscarded, l)
		}

		for name, field := range pathFields {
			if field.Description != "path" {
				t.Fatalf("expected field %s to be in 'path', found in %s", name, field.Description)
			}
		}
		for name, field := range queryFields {
			if field.Description != "query" {
				t.Fatalf("expected field %s to be in 'query', found in %s", name, field.Description)
			}
		}
		for name, field := range bodyFields {
			if field.Description != "body" {
				t.Fatalf("expected field %s to be in 'body', found in %s", name, field.Description)
			}
		}
	}
}

func TestOpenAPI_SpecialPaths(t *testing.T) {
	tests := map[string]struct {
		pattern                 string
		rootPaths               []string
		rootExpected            bool
		unauthenticatedPaths    []string
		unauthenticatedExpected bool
	}{
		"empty": {
			pattern:                 "foo",
			rootPaths:               []string{},
			rootExpected:            false,
			unauthenticatedPaths:    []string{},
			unauthenticatedExpected: false,
		},
		"exact-match-unauthenticated": {
			pattern:                 "foo",
			rootPaths:               []string{},
			rootExpected:            false,
			unauthenticatedPaths:    []string{"foo"},
			unauthenticatedExpected: true,
		},
		"exact-match-root": {
			pattern:                 "foo",
			rootPaths:               []string{"foo"},
			rootExpected:            true,
			unauthenticatedPaths:    []string{"bar"},
			unauthenticatedExpected: false,
		},
		"asterisk-match-unauthenticated": {
			pattern:                 "foo/bar",
			rootPaths:               []string{"foo"},
			rootExpected:            false,
			unauthenticatedPaths:    []string{"foo/*"},
			unauthenticatedExpected: true,
		},
		"asterisk-match-root": {
			pattern:                 "foo/bar",
			rootPaths:               []string{"foo/*"},
			rootExpected:            true,
			unauthenticatedPaths:    []string{"foo"},
			unauthenticatedExpected: false,
		},
		"path-ends-with-slash": {
			pattern:                 "foo/",
			rootPaths:               []string{"foo/*"},
			rootExpected:            true,
			unauthenticatedPaths:    []string{"a", "b", "foo*"},
			unauthenticatedExpected: true,
		},
		"asterisk-match-no-slash": {
			pattern:                 "foo",
			rootPaths:               []string{"foo*"},
			rootExpected:            true,
			unauthenticatedPaths:    []string{"a", "fo*"},
			unauthenticatedExpected: true,
		},
		"multiple-root-paths": {
			pattern:                 "foo/bar",
			rootPaths:               []string{"a", "b", "foo/*"},
			rootExpected:            true,
			unauthenticatedPaths:    []string{"foo/baz/*"},
			unauthenticatedExpected: false,
		},
		"plus-match-unauthenticated": {
			pattern:                 "foo/bar/baz",
			rootPaths:               []string{"foo/bar"},
			rootExpected:            false,
			unauthenticatedPaths:    []string{"foo/+/baz"},
			unauthenticatedExpected: true,
		},
		"plus-match-root": {
			pattern:                 "foo/bar/baz",
			rootPaths:               []string{"foo/+/baz"},
			rootExpected:            true,
			unauthenticatedPaths:    []string{"foo/bar"},
			unauthenticatedExpected: false,
		},
		"plus-and-asterisk": {
			pattern:                 "foo/bar/baz/something",
			rootPaths:               []string{"foo/+/baz/*"},
			rootExpected:            true,
			unauthenticatedPaths:    []string{"foo/+/baz*"},
			unauthenticatedExpected: true,
		},
		"double-plus-good": {
			pattern:                 "foo/bar/baz",
			rootPaths:               []string{"foo/+/+"},
			rootExpected:            true,
			unauthenticatedPaths:    []string{"foo/bar"},
			unauthenticatedExpected: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			doc := NewOASDocument("version")
			path := Path{
				Pattern: test.pattern,
			}
			backend := &Backend{
				PathsSpecial: &logical.Paths{
					Root:            test.rootPaths,
					Unauthenticated: test.unauthenticatedPaths,
				},
				BackendType: logical.TypeLogical,
			}

			if err := documentPath(&path, backend, "kv", doc); err != nil {
				t.Fatal(err)
			}

			actual := doc.Paths["/"+test.pattern].Sudo
			if actual != test.rootExpected {
				t.Fatalf("Test (root): expected: %v; got: %v", test.rootExpected, actual)
			}

			actual = doc.Paths["/"+test.pattern].Unauthenticated
			if actual != test.unauthenticatedExpected {
				t.Fatalf("Test (unauth): expected: %v; got: %v", test.unauthenticatedExpected, actual)
			}
		})
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
				"id": {
					Type:        TypeString,
					Description: "My id parameter",
				},
				"token": {
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

	t.Run("Operations - All Operations", func(t *testing.T) {
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
					Type:          TypeInt,
					Description:   "the age",
					AllowedValues: []interface{}{1, 2, 3},
					Required:      true,
					DisplayAttrs: &DisplayAttributes{
						Name:      "Age",
						Sensitive: true,
						Group:     "Some Group",
						Value:     7,
					},
				},
				"x-abc-token": {
					Type:          TypeHeader,
					Description:   "a header value",
					AllowedValues: []interface{}{"a", "b", "c"},
				},
				"maximum": {
					Type:        TypeInt64,
					Description: "a maximum value",
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
			DisplayAttrs: &DisplayAttributes{
				Navigation: true,
			},
		}

		sp := &logical.Paths{
			Root: []string{"foo*"},
		}
		testPath(t, p, sp, expected("operations"))
	})

	t.Run("Operations - List Only", func(t *testing.T) {
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
					Type:          TypeInt,
					Description:   "the age",
					AllowedValues: []interface{}{1, 2, 3},
					Required:      true,
					DisplayAttrs: &DisplayAttributes{
						Name:      "Age",
						Sensitive: true,
						Group:     "Some Group",
						Value:     7,
					},
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
				logical.ListOperation: &PathOperation{
					Summary:     "List Summary",
					Description: "List Description",
				},
			},
			DisplayAttrs: &DisplayAttributes{
				Navigation: true,
			},
		}

		sp := &logical.Paths{
			Root: []string{"foo*"},
		}
		testPath(t, p, sp, expected("operations_list"))
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
							Fields: map[string]*FieldSchema{
								"field_a": {
									Type:        TypeString,
									Description: "field_a description",
								},
								"field_b": {
									Type:        TypeBool,
									Description: "field_b description",
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

	docOrig := NewOASDocument("version")
	err := documentPath(p, &Backend{BackendType: logical.TypeLogical}, "kv", docOrig)
	if err != nil {
		t.Fatal(err)
	}

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
		Secret:    new(logical.Secret),
		Auth:      new(logical.Auth),
		Data:      map[string]interface{}{"foo": 42},
		Redirect:  "foo",
		Warnings:  []string{"foo"},
		WrapInfo:  &wrapping.ResponseWrapInfo{Token: "foo"},
		Headers:   map[string][]string{"foo": {"bar"}},
		MountType: "mount",
	}
	origJSON := mustJSONMarshal(t, orig)

	cr = cleanResponse(orig)

	cleanJSON := mustJSONMarshal(t, cr)

	if diff := deep.Equal(origJSON, cleanJSON); diff != nil {
		t.Fatal(diff)
	}
}

func TestOpenAPI_constructOperationID(t *testing.T) {
	tests := map[string]struct {
		path                string
		pathIndex           int
		pathAttributes      *DisplayAttributes
		operation           logical.Operation
		operationAttributes *DisplayAttributes
		defaultPrefix       string
		expected            string
	}{
		"empty": {
			path:                "",
			pathIndex:           0,
			pathAttributes:      nil,
			operation:           logical.Operation(""),
			operationAttributes: nil,
			defaultPrefix:       "",
			expected:            "",
		},
		"simple-read": {
			path:                "path/to/thing",
			pathIndex:           0,
			pathAttributes:      nil,
			operation:           logical.ReadOperation,
			operationAttributes: nil,
			defaultPrefix:       "test",
			expected:            "test-read-path-to-thing",
		},
		"simple-write": {
			path:                "path/to/thing",
			pathIndex:           0,
			pathAttributes:      nil,
			operation:           logical.UpdateOperation,
			operationAttributes: nil,
			defaultPrefix:       "test",
			expected:            "test-write-path-to-thing",
		},
		"operation-verb": {
			path:                "path/to/thing",
			pathIndex:           0,
			pathAttributes:      &DisplayAttributes{OperationVerb: "do-something"},
			operation:           logical.UpdateOperation,
			operationAttributes: nil,
			defaultPrefix:       "test",
			expected:            "do-something",
		},
		"operation-verb-override": {
			path:                "path/to/thing",
			pathIndex:           0,
			pathAttributes:      &DisplayAttributes{OperationVerb: "do-something"},
			operation:           logical.UpdateOperation,
			operationAttributes: &DisplayAttributes{OperationVerb: "do-something-else"},
			defaultPrefix:       "test",
			expected:            "do-something-else",
		},
		"operation-prefix": {
			path:                "path/to/thing",
			pathIndex:           0,
			pathAttributes:      &DisplayAttributes{OperationPrefix: "my-prefix"},
			operation:           logical.UpdateOperation,
			operationAttributes: nil,
			defaultPrefix:       "test",
			expected:            "my-prefix-write-path-to-thing",
		},
		"operation-prefix-override": {
			path:                "path/to/thing",
			pathIndex:           0,
			pathAttributes:      &DisplayAttributes{OperationPrefix: "my-prefix"},
			operation:           logical.UpdateOperation,
			operationAttributes: &DisplayAttributes{OperationPrefix: "better-prefix"},
			defaultPrefix:       "test",
			expected:            "better-prefix-write-path-to-thing",
		},
		"operation-prefix-and-suffix": {
			path:                "path/to/thing",
			pathIndex:           0,
			pathAttributes:      &DisplayAttributes{OperationPrefix: "my-prefix", OperationSuffix: "my-suffix"},
			operation:           logical.UpdateOperation,
			operationAttributes: nil,
			defaultPrefix:       "test",
			expected:            "my-prefix-write-my-suffix",
		},
		"operation-prefix-and-suffix-override": {
			path:                "path/to/thing",
			pathIndex:           0,
			pathAttributes:      &DisplayAttributes{OperationPrefix: "my-prefix", OperationSuffix: "my-suffix"},
			operation:           logical.UpdateOperation,
			operationAttributes: &DisplayAttributes{OperationPrefix: "better-prefix", OperationSuffix: "better-suffix"},
			defaultPrefix:       "test",
			expected:            "better-prefix-write-better-suffix",
		},
		"operation-prefix-verb-suffix": {
			path:                "path/to/thing",
			pathIndex:           0,
			pathAttributes:      &DisplayAttributes{OperationPrefix: "my-prefix", OperationSuffix: "my-suffix", OperationVerb: "Create"},
			operation:           logical.UpdateOperation,
			operationAttributes: &DisplayAttributes{OperationPrefix: "better-prefix", OperationSuffix: "better-suffix"},
			defaultPrefix:       "test",
			expected:            "better-prefix-create-better-suffix",
		},
		"operation-prefix-verb-suffix-override": {
			path:                "path/to/thing",
			pathIndex:           0,
			pathAttributes:      &DisplayAttributes{OperationPrefix: "my-prefix", OperationSuffix: "my-suffix", OperationVerb: "Create"},
			operation:           logical.UpdateOperation,
			operationAttributes: &DisplayAttributes{OperationPrefix: "better-prefix", OperationSuffix: "better-suffix", OperationVerb: "Login"},
			defaultPrefix:       "test",
			expected:            "better-prefix-login-better-suffix",
		},
		"operation-prefix-verb": {
			path:                "path/to/thing",
			pathIndex:           0,
			pathAttributes:      nil,
			operation:           logical.UpdateOperation,
			operationAttributes: &DisplayAttributes{OperationPrefix: "better-prefix", OperationVerb: "Login"},
			defaultPrefix:       "test",
			expected:            "better-prefix-login",
		},
		"operation-verb-suffix": {
			path:                "path/to/thing",
			pathIndex:           0,
			pathAttributes:      nil,
			operation:           logical.UpdateOperation,
			operationAttributes: &DisplayAttributes{OperationVerb: "Login", OperationSuffix: "better-suffix"},
			defaultPrefix:       "test",
			expected:            "login-better-suffix",
		},
		"pipe-delimited-suffix-0": {
			path:                "path/to/thing",
			pathIndex:           0,
			pathAttributes:      nil,
			operation:           logical.UpdateOperation,
			operationAttributes: &DisplayAttributes{OperationPrefix: "better-prefix", OperationSuffix: "suffix0|suffix1"},
			defaultPrefix:       "test",
			expected:            "better-prefix-write-suffix0",
		},
		"pipe-delimited-suffix-1": {
			path:                "path/to/thing",
			pathIndex:           1,
			pathAttributes:      nil,
			operation:           logical.UpdateOperation,
			operationAttributes: &DisplayAttributes{OperationPrefix: "better-prefix", OperationSuffix: "suffix0|suffix1"},
			defaultPrefix:       "test",
			expected:            "better-prefix-write-suffix1",
		},
		"pipe-delimited-suffix-2-fallback": {
			path:                "path/to/thing",
			pathIndex:           2,
			pathAttributes:      nil,
			operation:           logical.UpdateOperation,
			operationAttributes: &DisplayAttributes{OperationPrefix: "better-prefix", OperationSuffix: "suffix0|suffix1"},
			defaultPrefix:       "test",
			expected:            "better-prefix-write-path-to-thing",
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actual := constructOperationID(
				test.path,
				test.pathIndex,
				test.pathAttributes,
				test.operation,
				test.operationAttributes,
				test.defaultPrefix,
			)
			if actual != test.expected {
				t.Fatalf("expected: %s; got: %s", test.expected, actual)
			}
		})
	}
}

func TestOpenAPI_hyphenatedToTitleCase(t *testing.T) {
	tests := map[string]struct {
		in       string
		expected string
	}{
		"simple": {
			in:       "test",
			expected: "Test",
		},
		"two-words": {
			in:       "two-words",
			expected: "TwoWords",
		},
		"three-words": {
			in:       "one-two-three",
			expected: "OneTwoThree",
		},
		"not-hyphenated": {
			in:       "something_like_this",
			expected: "Something_like_this",
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actual := hyphenatedToTitleCase(test.in)
			if actual != test.expected {
				t.Fatalf("expected: %s; got: %s", test.expected, actual)
			}
		})
	}
}

func testPath(t *testing.T, path *Path, sp *logical.Paths, expectedJSON string) {
	t.Helper()

	doc := NewOASDocument("dummyversion")
	if err := documentPath(path, &Backend{
		PathsSpecial: sp,
		BackendType:  logical.TypeLogical,
	}, "kv", doc); err != nil {
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
		// fmt.Println(string(docJSON)) // uncomment to debug generated JSON (very helpful when fixing tests)
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

	content := strings.Replace(string(data), "<vault_version>", "dummyversion", 1)

	return content
}

func mustJSONMarshal(t *testing.T, data interface{}) []byte {
	j, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return j
}
