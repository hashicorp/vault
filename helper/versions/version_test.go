package versions

import "testing"

func TestIsBuiltinVersion(t *testing.T) {
	for _, tc := range []struct {
		version string
		builtin bool
	}{
		{"v1.0.0+builtin", true},
		{"v2.3.4+builtin.vault", true},
		{"1.0.0+builtin.anythingelse", true},
		{"v1.0.0+other.builtin", true},
		{"v1.0.0+builtinbutnot", false},
		{"v1.0.0", false},
		{"not-a-semver", false},
	} {
		builtin := IsBuiltinVersion(tc.version)
		if builtin != tc.builtin {
			t.Fatalf("%s should give: %v, but got %v", tc.version, tc.builtin, builtin)
		}
	}
}
