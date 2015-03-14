package backend

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/vault"
)

func BenchmarkBackendRoute(b *testing.B) {
	patterns := []string{
		"foo",
		"bar/(?P<name>.+?)",
		"baz/(?P<name>what)",
		`aws/policy/(?P<policy>\w)`,
		`aws/(?P<policy>\w)`,
	}

	backend := &Backend{Paths: make([]*Path, 0, len(patterns))}
	for _, p := range patterns {
		backend.Paths = append(backend.Paths, &Path{Pattern: p})
	}

	// Warm any caches
	backend.Route("aws/policy/foo")

	// Reset the timer since we did a lot above
	b.ResetTimer()

	// Run through and route. We do a sanity check of the return value
	for i := 0; i < b.N; i++ {
		if p := backend.Route("aws/policy/foo"); p == nil {
			b.Fatal("p should not be nil")
		}
	}
}

func TestBackend_impl(t *testing.T) {
	var _ vault.LogicalBackend = new(Backend)
}

func TestBackendHandleRequest(t *testing.T) {
	callback := func(req *vault.Request, data *FieldData) (*vault.Response, error) {
		return &vault.Response{
			Data: map[string]interface{}{
				"value": data.Get("value"),
			},
		}, nil
	}

	b := &Backend{
		Paths: []*Path{
			&Path{
				Pattern: "foo/bar",
				Fields: map[string]*FieldSchema{
					"value": &FieldSchema{Type: TypeInt},
				},
				Callback: callback,
			},
		},
	}

	resp, err := b.HandleRequest(&vault.Request{
		Path: "foo/bar",
		Data: map[string]interface{}{"value": "42"},
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.Data["value"] != 42 {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestBackendHandleRequest_404(t *testing.T) {
	callback := func(req *vault.Request, data *FieldData) (*vault.Response, error) {
		return &vault.Response{
			Data: map[string]interface{}{
				"value": data.Get("value"),
			},
		}, nil
	}

	b := &Backend{
		Paths: []*Path{
			&Path{
				Pattern: `foo/bar`,
				Fields: map[string]*FieldSchema{
					"value": &FieldSchema{Type: TypeInt},
				},
				Callback: callback,
			},
		},
	}

	_, err := b.HandleRequest(&vault.Request{
		Path: "foo/baz",
		Data: map[string]interface{}{"value": "84"},
	})
	if err != vault.ErrUnsupportedPath {
		t.Fatalf("err: %s", err)
	}
}

func TestBackendHandleRequest_urlPriority(t *testing.T) {
	callback := func(req *vault.Request, data *FieldData) (*vault.Response, error) {
		return &vault.Response{
			Data: map[string]interface{}{
				"value": data.Get("value"),
			},
		}, nil
	}

	b := &Backend{
		Paths: []*Path{
			&Path{
				Pattern: `foo/(?P<value>\d+)`,
				Fields: map[string]*FieldSchema{
					"value": &FieldSchema{Type: TypeInt},
				},
				Callback: callback,
			},
		},
	}

	resp, err := b.HandleRequest(&vault.Request{
		Path: "foo/42",
		Data: map[string]interface{}{"value": "84"},
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.Data["value"] != 42 {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestBackendRoute(t *testing.T) {
	cases := map[string]struct {
		Patterns []string
		Path     string
		Match    string
	}{
		"no match": {
			[]string{"foo"},
			"bar",
			"",
		},

		"exact": {
			[]string{"foo"},
			"foo",
			"foo",
		},

		"regexp": {
			[]string{"fo+"},
			"foo",
			"fo+",
		},
	}

	for n, tc := range cases {
		paths := make([]*Path, len(tc.Patterns))
		for i, pattern := range tc.Patterns {
			paths[i] = &Path{Pattern: pattern}
		}

		b := &Backend{Paths: paths}
		result := b.Route(tc.Path)
		match := ""
		if result != nil {
			match = result.Pattern
		}

		if match != tc.Match {
			t.Fatalf("bad: %s\n\nExpected: %s\nGot: %s",
				n, tc.Match, match)
		}
	}
}

func TestFieldSchemaDefaultOrZero(t *testing.T) {
	cases := map[string]struct {
		Schema *FieldSchema
		Value  interface{}
	}{
		"default set": {
			&FieldSchema{Type: TypeString, Default: "foo"},
			"foo",
		},

		"default not set": {
			&FieldSchema{Type: TypeString},
			"",
		},
	}

	for name, tc := range cases {
		actual := tc.Schema.DefaultOrZero()
		if !reflect.DeepEqual(actual, tc.Value) {
			t.Fatalf("bad: %s\n\nExpected: %#v\nGot: %#v",
				name, tc.Value, actual)
		}
	}
}
