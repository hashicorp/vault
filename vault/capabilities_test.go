// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/internalshared/namespace"
	"github.com/hashicorp/vault/sdk/logical"
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
	policy, _ := ParseACLPolicy(namespace.RootNamespace, policy1, WithDenySlashInTemplatedPaths(c.denySlashInTemplatedPolicyPaths))
	err = c.policyStore.SetPolicy(ctx, policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	policy, _ = ParseACLPolicy(namespace.RootNamespace, policy2, WithDenySlashInTemplatedPaths(c.denySlashInTemplatedPolicyPaths))
	err = c.policyStore.SetPolicy(ctx, policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	policy, _ = ParseACLPolicy(namespace.RootNamespace, policy3, WithDenySlashInTemplatedPaths(c.denySlashInTemplatedPolicyPaths))
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
		policy, err := ParseACLPolicy(namespace.RootNamespace, tCase.policy, WithDenySlashInTemplatedPaths(c.denySlashInTemplatedPolicyPaths))
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
	policy, _ := ParseACLPolicy(namespace.RootNamespace, aclPolicy, WithDenySlashInTemplatedPaths(c.denySlashInTemplatedPolicyPaths))
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

// TestSystemBackend_Capabilities_ManyPaths checks that asking about several
// paths in one request returns the same capabilities as asking about each path
// on its own.
func TestSystemBackend_Capabilities_ManyPaths(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)

	policy, err := ParseACLPolicy(namespace.RootNamespace, `
path "sys/capabilities-self" { capabilities = ["update"] }
path "secret/data/app/*" { capabilities = ["read", "list"] }
path "secret/data/app/locked" { capabilities = ["deny"] }
path "secret/+/team/*" { capabilities = ["create"] }
path "kv/exact" { capabilities = ["read", "update"] }
`)
	if err != nil {
		t.Fatal(err)
	}
	policy.Name = "many-paths"
	if err := c.policyStore.SetPolicy(ctx, policy); err != nil {
		t.Fatal(err)
	}
	testMakeServiceTokenViaBackend(t, c.tokenStore, root, "manypathstoken", "", []string{"many-paths"})
	testMakeTokenDirectly(t, c.tokenStore, &logical.TokenEntry{
		ID:   "nopolicytoken",
		Path: "testpath",
		TTL:  time.Hour,
	})

	paths := []string{
		"secret/data/app/one",
		"secret/data/app/locked",
		"secret/metadata/team/two",
		"kv/exact",
		"kv/other",
		"sys/capabilities-self",
		"kv/exact",
	}

	for _, token := range []string{root, "manypathstoken", "nopolicytoken"} {
		expected := map[string]interface{}{}
		for _, path := range paths {
			caps, err := c.Capabilities(ctx, token, path)
			if err != nil {
				t.Fatal(err)
			}
			expected[path] = caps
		}

		for _, endpoint := range []string{"sys/capabilities", "sys/capabilities-self"} {
			req := &logical.Request{
				ClientToken: root,
				Path:        endpoint,
				Operation:   logical.UpdateOperation,
				Data:        map[string]interface{}{"paths": paths, "token": token},
			}
			if endpoint == "sys/capabilities-self" {
				if token == "nopolicytoken" {
					// This token is not allowed to call sys/capabilities-self.
					continue
				}
				req.ClientToken = token
				delete(req.Data, "token")
			}
			resp, err := c.HandleRequest(ctx, req)
			if err != nil || resp == nil || resp.IsError() {
				t.Fatalf("%s with token %q: resp: %#v, err: %v", endpoint, token, resp, err)
			}
			if !reflect.DeepEqual(resp.Data, expected) {
				t.Fatalf("%s with token %q:\ngot      %#v\nexpected %#v", endpoint, token, resp.Data, expected)
			}
		}
	}

	// An empty path anywhere in the list is still rejected.
	req := &logical.Request{
		ClientToken: "manypathstoken",
		Path:        "sys/capabilities-self",
		Operation:   logical.UpdateOperation,
		Data:        map[string]interface{}{"paths": []string{"kv/exact", ""}},
	}
	if _, err := c.HandleRequest(ctx, req); err == nil || !strings.Contains(err.Error(), "missing path") {
		t.Fatalf("expected a missing path error, got: %v", err)
	}
}

// BenchmarkSystemBackend_CapabilitiesSelf measures sys/capabilities-self for a
// token with many policies when one request asks about many paths.
func BenchmarkSystemBackend_CapabilitiesSelf(b *testing.B) {
	c, _, root := TestCoreUnsealed(b)
	ctx := namespace.RootContext(nil)

	const numPolicies = 40
	var policyNames []string
	for i := 0; i < numPolicies; i++ {
		name := fmt.Sprintf("bench-%d", i)
		policyNames = append(policyNames, name)
		raw := `path "sys/capabilities-self" { capabilities = ["update"] }` + "\n"
		for j := 0; j < 20; j++ {
			raw += fmt.Sprintf("path \"engine%d/+/team%d/*\" { capabilities = [\"read\", \"list\"] }\n", i, j)
			raw += fmt.Sprintf("path \"engine%d/data/team%d/app/secret\" { capabilities = [\"read\", \"update\"] }\n", i, j)
		}
		p, err := ParseACLPolicy(namespace.RootNamespace, raw)
		if err != nil {
			b.Fatal(err)
		}
		p.Name = name
		if err := c.policyStore.SetPolicy(ctx, p); err != nil {
			b.Fatal(err)
		}
	}
	testMakeServiceTokenViaBackend(b, c.tokenStore, root, "benchtoken", "", policyNames)

	for _, numPaths := range []int{1, 10, 100} {
		paths := make([]string, numPaths)
		for i := range paths {
			if i%2 == 0 {
				// Matches an exact path rule.
				paths[i] = fmt.Sprintf("engine%d/data/team%d/app/secret", i%numPolicies, i%20)
			} else {
				// Matches a "+" and "*" wildcard rule.
				paths[i] = fmt.Sprintf("engine%d/metadata/team%d/app%d", i%numPolicies, i%20, i)
			}
		}
		b.Run(fmt.Sprintf("paths=%d", numPaths), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				req := &logical.Request{
					ClientToken: "benchtoken",
					Path:        "sys/capabilities-self",
					Operation:   logical.UpdateOperation,
					Data:        map[string]interface{}{"paths": paths},
				}
				resp, err := c.HandleRequest(ctx, req)
				if err != nil || resp == nil || resp.IsError() {
					b.Fatalf("bad: resp: %#v err: %v", resp, err)
				}
				if got := resp.Data[paths[0]]; !reflect.DeepEqual(got, []string{"read", "update"}) {
					b.Fatalf("unexpected capabilities: %#v", got)
				}
			}
		})
	}
}
