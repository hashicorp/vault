package framework

import (
	"encoding/json"
	"reflect"
	"sort"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/openapi"
	"github.com/hashicorp/vault/logical"
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
		documentPath(&path, test.rootPaths, doc)
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
					Description: "My id param",
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

		testPath(t, p, expectedJSON["Legacy callbacks"])
	})

	t.Run("Simple (new)", func(t *testing.T) {
		p := &Path{
			Pattern: "simple",

			HelpSynopsis:    "Synopsis",
			HelpDescription: "Description",
			Operations: map[logical.Operation]OperationHandler{
				logical.ReadOperation: &PathOperation{
					Callback:    nil,
					Summary:     "My Summary",
					Description: "My Description",
				},
			},
		}

		expectedJSON := `{
  "openapi": "3.0.2",
  "info": {
    "title": "HashiCorp Vault API",
    "version": "0.11.3"
  },
  "paths": {
    "/simple": {
      "description": "Synopsis",
      "get": {
        "summary": "My Summary",
        "description": "My Description",
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      }
    }
  }
}
`
		testPath(t, p, expectedJSON)
	})
}

func testPath(t *testing.T, path *Path, expectedJSON string) {
	t.Helper()

	doc := openapi.NewDocument()
	documentPath(path, []string{}, doc)

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
		// fmt.Println(string(docJSON)) // uncomment to debug generated JSON, which can be lengthy
		t.Fatal(diff)
	}
}

var expectedJSON = map[string]string{
	"Legacy callbacks": `
{
  "openapi": "3.0.2",
  "info": {
    "title": "HashiCorp Vault API",
    "version": "0.11.3"
  },
  "paths": {
    "/lookup/{id}": {
      "description": "Synopsis",
      "get": {
        "summary": "Synopsis",
        "parameters": [
          {
            "name": "id",
            "description": "My id param",
            "in": "path",
            "schema": {
              "type": "string"
            },
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      },
      "post": {
        "summary": "Synopsis",
        "parameters": [
          {
            "name": "id",
            "description": "My id param",
            "in": "path",
            "schema": {
              "type": "string"
            },
            "required": true
          }
        ],
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "token": {
                    "type": "string",
                    "description": "My token"
                  }
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      }
    }
  }
}
`,
}
