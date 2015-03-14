package backend

import (
	"reflect"
	"testing"
)

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
