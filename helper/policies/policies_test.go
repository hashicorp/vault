package policies

import "testing"

func TestEquivalentPolicies(t *testing.T) {
	a := []string{"foo", "bar"}
	var b []string
	if EquivalentPolicies(a, b) {
		t.Fatal("bad")
	}

	b = []string{"foo"}
	if EquivalentPolicies(a, b) {
		t.Fatal("bad")
	}

	b = []string{"bar", "foo"}
	if !EquivalentPolicies(a, b) {
		t.Fatal("bad")
	}

	b = []string{"foo", "default", "bar"}
	if !EquivalentPolicies(a, b) {
		t.Fatal("bad")
	}
}
