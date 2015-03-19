package framework

import (
	"testing"
)

func TestSecretType(t *testing.T) {
	cases := [][2]string{
		{"foo-bar", "foo"},
		{"foo", ""},
	}

	for _, tc := range cases {
		actual := SecretType(tc[0])
		if actual != tc[1] {
			t.Fatalf("Input: %s, Output: %s, Expected: %s", tc[0], actual, tc[1])
		}
	}
}
