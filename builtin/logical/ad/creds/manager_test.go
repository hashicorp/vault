package creds

import (
	"regexp"
	"strings"
	"testing"
)

func TestGeneratePassword(t *testing.T) {

	for desiredLength := -200; desiredLength < 200; desiredLength++ {

		password1, err := generatePassword(desiredLength)

		if desiredLength < 14 {
			if err == nil {
				t.Fatalf("desiredLength of %d should yield an error", desiredLength)
			} else {
				// password1 won't be populated, nothing more to check
				continue
			}
		}

		// desired length is appropriate
		if err != nil {
			t.Fatalf("desiredLength of %d generated an err: %s", desiredLength, err)
		}
		if len(password1) != desiredLength {
			t.Fatalf("unexpected password1 length of %d for desired length of %d", len(password1), desiredLength)
		}

		// let's generate a second password1 to ensure it's not the same
		password2, err := generatePassword(desiredLength)
		if err != nil {
			t.Fatalf("desiredLength of %d generated an err: %s", desiredLength, err)
		}

		if password1 == password2 {
			t.Fatalf("received identical passwords of %s, random byte generation is broken", password1)
		}
	}
}

func TestPathRegexp(t *testing.T) {

	m := &Manager{}
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
