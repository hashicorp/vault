package vault

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/logical"
)

func TestCapabilities_DerivedPolicies(t *testing.T) {
	var resp *logical.Response
	var err error

	ctx := namespace.RootContext(nil)
	i, _, c := testIdentityStoreWithGithubAuth(ctx, t)

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
	policy, _ := ParseACLPolicy(namespace.RootNamespace, policy1)
	err = c.policyStore.SetPolicy(ctx, policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	policy, _ = ParseACLPolicy(namespace.RootNamespace, policy2)
	err = c.policyStore.SetPolicy(ctx, policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	policy, _ = ParseACLPolicy(namespace.RootNamespace, policy3)
	err = c.policyStore.SetPolicy(ctx, policy)
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
	resp, err = i.HandleRequest(ctx, entityReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %#v\n", resp, err)
	}
	entityID := resp.Data["id"].(string)

	// Create a token for the entity and assign policy2 on the token
	ent := &logical.TokenEntry{
		ID:       "capabilitiestoken",
		Path:     "secret/sample",
		Policies: []string{"policy2"},
		EntityID: entityID,
		TTL:      time.Hour,
	}
	testMakeTokenDirectly(t, c.tokenStore, ent)

	actual, err := c.Capabilities(ctx, "capabilitiestoken", "secret/sample")
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
	resp, err = i.HandleRequest(ctx, groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %#v\n", resp, err)
	}

	actual, err = c.Capabilities(namespace.RootContext(nil), "capabilitiestoken", "secret/sample")
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

func TestCapabilities_TemplatedPolicies(t *testing.T) {
	var resp *logical.Response
	var err error
	i, _, c := testIdentityStoreWithGithubAuth(namespace.RootContext(nil), t)
	// Create an entity and assign policy1 to it
	entityReq := &logical.Request{
		Path:      "entity",
		Operation: logical.UpdateOperation,
	}
	resp, err = i.HandleRequest(namespace.RootContext(nil), entityReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %#v\n", resp, err)
	}
	entityID := resp.Data["id"].(string)

	// Create a token for the entity and assign policy2 on the token
	ent := &logical.TokenEntry{
		ID:       "capabilitiestoken",
		Path:     "auth/token/create",
		Policies: []string{"testpolicy"},
		EntityID: entityID,
		TTL:      time.Hour,
	}
	testMakeTokenDirectly(t, c.tokenStore, ent)

	tCases := []struct {
		policy   string
		path     string
		expected []string
	}{
		{
			`name = "testpolicy"
			path "secret/{{identity.entity.id}}/sample" {
				capabilities = ["update", "create"]
			}
			`,
			fmt.Sprintf("secret/%s/sample", entityID),
			[]string{"update", "create"},
		},
		{
			`{"name": "testpolicy", "path": {"secret/{{identity.entity.id}}/sample": {"capabilities": ["read", "create"]}}}`,
			fmt.Sprintf("secret/%s/sample", entityID),
			[]string{"read", "create"},
		},
		{
			`{"name": "testpolicy", "path": {"secret/sample": {"capabilities": ["read"]}}}`,
			"secret/sample",
			[]string{"read"},
		},
	}
	for _, tCase := range tCases {
		// Create the above policies
		policy, err := ParseACLPolicy(namespace.RootNamespace, tCase.policy)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		err = c.policyStore.SetPolicy(namespace.RootContext(nil), policy)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		actual, err := c.Capabilities(namespace.RootContext(nil), "capabilitiestoken", tCase.path)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		sort.Strings(actual)
		sort.Strings(tCase.expected)
		if !reflect.DeepEqual(actual, tCase.expected) {
			t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, tCase.expected)
		}
	}
}

func TestCapabilities(t *testing.T) {
	c, _, token := TestCoreUnsealed(t)

	actual, err := c.Capabilities(namespace.RootContext(nil), token, "path")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	expected := []string{"root"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}

	// Create a policy
	policy, _ := ParseACLPolicy(namespace.RootNamespace, aclPolicy)
	err = c.policyStore.SetPolicy(namespace.RootContext(nil), policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Create a token for the policy
	ent := &logical.TokenEntry{
		ID:       "capabilitiestoken",
		Path:     "testpath",
		Policies: []string{"dev"},
		TTL:      time.Hour,
	}
	testMakeTokenDirectly(t, c.tokenStore, ent)

	actual, err = c.Capabilities(namespace.RootContext(nil), "capabilitiestoken", "foo/bar")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	expected = []string{"create", "read", "sudo"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}
}
