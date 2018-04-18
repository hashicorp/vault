package creds

import (
	"regexp"
	"strings"
	"testing"
)

func TestPathRegexp(t *testing.T) {

	m := &handler{}
	path := m.Path()
	re, err := regexp.Compile(path.Pattern)
	if err != nil {
		t.Fatal(err)
	}

	matches := re.FindStringSubmatch("creds")
	if len(matches) > 0 {
		t.Fatal("creds shouldn't be a path that's hit by the matcher")
	}

	matches = re.FindStringSubmatch("creds/")
	if len(matches) > 0 {
		t.Fatal("creds/ shouldn't be a path that's hit by the matcher")
	}

	matches = re.FindStringSubmatch("credssuper")
	if len(matches) > 0 {
		t.Fatal("credssuper shouldn't be a path that's hit by the matcher")
	}

	matches = re.FindStringSubmatch("creds/candy")
	if len(matches) <= 0 {
		t.Fatal("creds/candy should be a path that's hit by the matcher")
	}

	matches = re.FindStringSubmatch("cats/creds")
	if len(matches) > 0 {
		t.Fatal("cats/creds shouldn't be a path that's hit by the matcher")
	}

	if !strings.HasPrefix(path.Pattern, "^") {
		t.Fatal("pattern needs to start with a ^ or it'll be added outside the package and the regex won't behave as expected")
	}

	if !strings.HasSuffix(path.Pattern, "$") {
		t.Fatal("pattern needs to end with a $ or it'll be added outside the package and the regex won't behave as expected")
	}
}
