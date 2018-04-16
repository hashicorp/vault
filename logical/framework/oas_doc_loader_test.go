package framework

import (
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/vault/helper/oas"
)

func TestParsePattern(t *testing.T) {
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
			"raw/",
			"raw/{path}",
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

func TestPathFields(t *testing.T) {
	tests := []struct {
		in_pattern string
		out_params []string
	}{
		{"/sys/{foo}/test/{bar}", []string{"foo", "bar"}},
		{"/sys/foo/test/bar", []string{}},
	}
	for i, test := range tests {
		out := oas.PathFields(test.in_pattern)
		if !reflect.DeepEqual(out, test.out_params) {
			t.Fatalf("Test %d: Expected %v got %v", i, test.out_params, out)
		}
	}
}

func TestRootPath(t *testing.T) {
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
		doc := oas.NewOASDoc()
		path := Path{
			Pattern: test.pattern,
		}
		documentPath(&path, test.rootPaths, &doc)
		result := test.root
		if doc.Paths["/"+test.pattern].Root != result {
			t.Fatalf("Test %d: Expected %v got %v", i, test.root, result)
		}
	}
}
