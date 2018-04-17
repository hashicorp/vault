package creds

import (
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
