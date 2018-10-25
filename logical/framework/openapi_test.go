package framework

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/openapi"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/version"
)

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
		{`config/ui/headers/(?P<header>\w(([\w-.]+)?\w)?)`, []string{"config/ui/headers/{header}"}},
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

func TestOpenAPI_RootPath(t *testing.T) {
	tests := []struct {
		pattern   string
		rootPaths []string
		root      bool
	}{
		{"foo", []string{}, false},
		{"foo", []string{"foo"}, true},
		{"foo/bar", []string{"foo"}, false},
		{"foo/bar", []string{"foo/*"}, true},
		{"foo/", []string{"foo/*"}, true},
		{"foo", []string{"foo*"}, true},
		{"foo/bar", []string{"a", "b", "foo/*"}, true},
	}
	for i, test := range tests {
		doc := openapi.NewDocument()
		path := Path{
			Pattern: test.pattern,
		}
		documentPath(&path, test.rootPaths, logical.TypeLogical, doc)
		result := test.root
		if doc.Paths["/"+test.pattern].Sudo != result {
			t.Fatalf("Test %d: Expected %v got %v", i, test.root, result)
		}
	}
}

func TestOpenAPIPaths(t *testing.T) {
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

		testPath(t, p, expected("legacy"))
	})

	t.Run("Simple", func(t *testing.T) {
		p := &Path{
			Pattern: "foo/" + GenericNameRegex("id"),
			Fields: map[string]*FieldSchema{
				"id": {
					Type:        TypeString,
					Description: "id path parameter",
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

		testPath(t, p, expected("simple"))
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
					Responses: map[string][]Response{
						"202": {{
							Description: "Amazing",
							Example: &logical.Response{
								Data: map[string]interface{}{
									"amount": 42,
								},
							},
						}},
					},
				},
			},
		}

		testPath(t, p, expected("responses"))
	})
}

func testPath(t *testing.T, path *Path, expectedJSON string) {
	t.Helper()

	doc := openapi.NewDocument()
	documentPath(path, []string{}, logical.TypeLogical, doc)

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
		//fmt.Println(string(docJSON)) // uncomment to debug generated JSON
		t.Fatal(diff)
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
