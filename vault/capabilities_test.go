package vault

import (
	"reflect"
	"testing"
)

func TestCapabilities_Basic(t *testing.T) {
	c, _, token := TestCoreUnsealed(t)

	actual, err := c.Capabilities(token, "path")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	expected := []string{"root"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}

	// Create a policy
	policy, _ := Parse(aclPolicy)
	err = c.policyStore.SetPolicy(policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Create a token for the policy
	ent := &TokenEntry{
		ID:       "capabilitiestoken",
		Path:     "testpath",
		Policies: []string{"dev"},
	}
	if err := c.tokenStore.create(ent); err != nil {
		t.Fatalf("err: %v", err)
	}

	actual, err = c.Capabilities("capabilitiestoken", "foo/bar")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	expected = []string{"sudo", "read", "create"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}
}
