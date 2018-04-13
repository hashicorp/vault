package roles

import (
	"regexp"
	"strings"
	"testing"
)

func TestPathRegexpWorks(t *testing.T) {

	m := &Manager{}
	path := m.Path()
	re, err := regexp.Compile(path.Pattern)
	if err != nil {
		t.Fatal(err)
	}

	matches := re.FindStringSubmatch("roles")
	if len(matches) <= 0 {
		t.Fatal("roles should be a path that's hit by the matcher")
	}

	matches = re.FindStringSubmatch("roles/")
	if len(matches) <= 0 {
		t.Fatal("roles/ should be a path that's hit by the matcher")
	}

	matches = re.FindStringSubmatch("rolessuper")
	if len(matches) > 0 {
		t.Fatal("rolessuper shouldn't be a path that's hit by the matcher")
	}

	matches = re.FindStringSubmatch("roles/candy")
	if len(matches) <= 0 {
		t.Fatal("roles/candy should be a path that's hit by the matcher")
	}

	matches = re.FindStringSubmatch("cats/roles")
	if len(matches) > 0 {
		t.Fatal("cats/roles shouldn't be a path that's hit by the matcher")
	}

	if !strings.HasPrefix(path.Pattern, "^") {
		t.Fatal("pattern needs to start with a ^ or it'll be added outside the package and the regex won't behave as expected")
	}

	if !strings.HasSuffix(path.Pattern, "$") {
		t.Fatal("pattern needs to send with a $ or it'll be added outside the package and the regex won't behave as expected")
	}
}
