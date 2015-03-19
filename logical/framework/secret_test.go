package framework

import (
	"testing"
)

func TestSecretType(t *testing.T) {
	cases := [][3]string{
		{"foo-bar", "foo", "bar"},
		{"foo", "", "foo"},
	}

	for _, tc := range cases {
		actual, actual2 := SecretType(tc[0])
		if actual != tc[1] {
			t.Fatalf("Input: %s, Output: %s, Expected: %s", tc[0], actual, tc[1])
		}
		if actual2 != tc[2] {
			t.Fatalf("Input: %s, Output: %s, Expected: %s", tc[0], actual2, tc[2])
		}
	}
}
