package vault

import (
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestCapabilities_DerivedPolicies(t *testing.T) {
	var resp *logical.Response
	var err error

	i, _, c := testIdentityStoreWithGithubAuth(t)

	policy1 := `
name = "policy1"
path "secret/sample" {
	capabilities = ["update", "create", "sudo"]
}
`
	policy2 := `
name = "policy2"
path "secret/sample" {
	capabilities = ["read", "delete"]
}
`

	policy3 := `
name = "policy3"
path "secret/sample" {
	capabilities = ["list", "list"]
}
`
	// Create the above policies
	policy, _ := ParseACLPolicy(policy1)
	err = c.policyStore.SetPolicy(policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	policy, _ = ParseACLPolicy(policy2)
	err = c.policyStore.SetPolicy(policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	policy, _ = ParseACLPolicy(policy3)
	err = c.policyStore.SetPolicy(policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Create an entity and assign policy1 to it
	entityReq := &logical.Request{
		Path:      "entity",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"policies": "policy1",
		},
	}
	resp, err = i.HandleRequest(entityReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %#v\n", resp, err)
	}
	entityID := resp.Data["id"].(string)

	// Create a token for the entity and assign policy2 on the token
	ent := &TokenEntry{
		ID:       "capabilitiestoken",
		Path:     "secret/sample",
		Policies: []string{"policy2"},
		EntityID: entityID,
	}
	if err := c.tokenStore.create(ent); err != nil {
		t.Fatalf("err: %v", err)
	}

	actual, err := c.Capabilities("capabilitiestoken", "secret/sample")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	expected := []string{"create", "read", "sudo", "delete", "update"}
	sort.Strings(actual)
	sort.Strings(expected)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}

	// Create a group and add the above created entity to it
	groupReq := &logical.Request{
		Path:      "group",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"member_entity_ids": []string{entityID},
			"policies":          "policy3",
		},
	}
	resp, err = i.HandleRequest(groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %#v\n", resp, err)
	}

	actual, err = c.Capabilities("capabilitiestoken", "secret/sample")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	expected = []string{"create", "read", "sudo", "delete", "update", "list"}
	sort.Strings(actual)
	sort.Strings(expected)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}
}

func TestCapabilities(t *testing.T) {
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
	policy, _ := ParseACLPolicy(aclPolicy)
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
	expected = []string{"create", "read", "sudo"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}
}
