// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"sync"
	"testing"

	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/internalshared/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPolicy_ResolveTemplatedPathsMatchesReparse verifies that resolving a
// parsed templated policy produces the same policy as parsing its raw HCL
// with templating, which is what ACL did before it reused parsed policies.
func TestPolicy_ResolveTemplatedPathsMatchesReparse(t *testing.T) {
	t.Parallel()

	childNS := &namespace.Namespace{ID: "child", Path: "child/"}
	entity := &identity.Entity{
		ID:          "entity-id",
		Name:        "entity-name",
		NamespaceID: namespace.RootNamespaceID,
		Metadata: map[string]string{
			"team":  "eng",
			"slash": "a/b",
			"plus":  "+",
			"star":  "*",
		},
		Aliases: []*identity.Alias{{MountAccessor: "auth_userpass_1234", Name: "alias-name"}},
	}
	groups := []*identity.Group{{ID: "group-id", Name: "group-name", NamespaceID: namespace.RootNamespaceID}}

	tests := map[string]struct {
		ns     *namespace.Namespace
		policy string
		entity *identity.Entity
		groups []*identity.Group
		opts   parseACLPolicyOptions
		// paths is the number of paths left after templating
		paths int
	}{
		"entity fields and metadata": {
			policy: `
path "secret/{{identity.entity.id}}/*" { capabilities = ["read", "list"] }
path "secret/{{identity.entity.name}}" { capabilities = ["create"] }
path "kv/{{identity.entity.metadata.team}}/data" {
	capabilities        = ["update"]
	required_parameters = ["a", "b"]
	allowed_parameters  = { "a" = [] "b" = ["x"] }
	min_wrapping_ttl    = "1m"
	max_wrapping_ttl    = "1h"
}`,
			paths: 3,
		},
		"aliases and groups": {
			policy: `
path "auth/{{identity.entity.aliases.auth_userpass_1234.name}}" { capabilities = ["read"] }
path "groups/{{identity.groups.names.group-name.id}}/*" { capabilities = ["read"] }
path "groups/{{identity.groups.ids.group-id.name}}" { capabilities = ["deny"] }`,
			paths: 3,
		},
		"mixed templated and plain paths": {
			policy: `
path "sys/health" { capabilities = ["read", "sudo"] }
path "secret/{{identity.entity.id}}" { policy = "write" }
path "secret/+/shared/*" { capabilities = ["read"] }`,
			paths: 3,
		},
		"unresolvable templates are dropped": {
			policy: `
path "secret/{{identity.entity.metadata.missing}}" { capabilities = ["read"] }
path "secret/{{identity.groups.names.unknown.id}}" { capabilities = ["read"] }
path "secret/kept" { capabilities = ["read"] }`,
			paths: 1,
		},
		"templated paths with literal wildcards": {
			policy: `
path "secret/+/{{identity.entity.id}}/*" { capabilities = ["read"] }
path "secret/{{identity.entity.id}}/+" { capabilities = ["read"] }
path "secret/{{identity.entity.id}}*" { capabilities = ["read"] }`,
			paths: 3,
		},
		"wildcards in template values are rejected": {
			policy: `
path "secret/{{identity.entity.metadata.plus}}/data" { capabilities = ["read"] }
path "secret/{{identity.entity.metadata.star}}" { capabilities = ["read"] }`,
		},
		"slash in value allowed": {
			policy: `path "secret/{{identity.entity.metadata.slash}}/data" { capabilities = ["read"] }`,
			paths:  1,
		},
		"slash in value denied": {
			policy: `path "secret/{{identity.entity.metadata.slash}}/data" { capabilities = ["read"] }`,
			opts:   parseACLPolicyOptions{denySlashInTemplatedPaths: true},
		},
		"no entity": {
			policy: `
path "secret/{{identity.entity.id}}" { capabilities = ["read"] }
path "secret/plain" { capabilities = ["read"] }`,
			entity: &identity.Entity{},
			paths:  1,
		},
		"child namespace": {
			ns:     childNS,
			policy: `path "/secret/{{identity.entity.id}}/*" { capabilities = ["read"] }`,
			paths:  1,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ns := tc.ns
			if ns == nil {
				ns = namespace.RootNamespace
			}
			e, g := tc.entity, tc.groups
			if e == nil {
				e, g = entity, groups
			} else if e.ID == "" {
				e = nil
			}

			parsed, err := ParseACLPolicy(ns, tc.policy, WithDenySlashInTemplatedPaths(tc.opts.denySlashInTemplatedPaths))
			require.NoError(t, err)
			require.True(t, parsed.Templated)
			require.True(t, parsed.canResolveTemplatedPaths())

			want, err := parseACLPolicyWithTemplating(ns, tc.policy, true, e, g, tc.opts)
			require.NoError(t, err)

			got, err := parsed.resolveTemplatedPaths(identity.ToSDKEntity(e), identity.ToSDKGroups(g), tc.opts)
			require.NoError(t, err)
			require.Equal(t, want, got)
			require.Len(t, got.Paths, tc.paths)
		})
	}
}

// TestPolicyStore_ACLTemplatedPolicyWithoutParsedPaths verifies that ACL
// still parses the raw HCL of a templated policy whose paths were not
// produced by the parser.
func TestPolicyStore_ACLTemplatedPolicyWithoutParsedPaths(t *testing.T) {
	t.Parallel()

	ps := NewTestCluster(t, &CoreConfig{}, nil).Cores[0].Core.policyStore
	ctx := namespace.RootContext(context.Background())

	p := &Policy{
		Name:      "templated-without-paths",
		Type:      PolicyTypeACL,
		Templated: true,
		Raw:       `path "secret/{{identity.entity.id}}" { capabilities = ["read"] }`,
		namespace: namespace.RootNamespace,
	}
	require.False(t, p.canResolveTemplatedPaths())
	require.NoError(t, ps.SetPolicy(ctx, p))

	acl, err := ps.ACL(ctx, &identity.Entity{ID: "entity-id"}, map[string][]string{namespace.RootNamespaceID: {p.Name}})
	require.NoError(t, err)
	res := acl.AllowOperation(ctx, &logical.Request{Operation: logical.ReadOperation, Path: "secret/entity-id"}, false)
	require.True(t, res.Allowed)
}

// TestPolicyStore_ACLTemplatedConcurrent builds ACLs concurrently from the
// same cached templated policies, which now share their parsed path rules,
// and verifies every ACL grants the merged permissions. Run with -race to
// check that building an ACL never writes to the shared rules.
func TestPolicyStore_ACLTemplatedConcurrent(t *testing.T) {
	t.Parallel()

	ps := NewTestCluster(t, &CoreConfig{}, nil).Cores[0].Core.policyStore
	ctx := namespace.RootContext(context.Background())

	policies := map[string]string{
		"templated-a": `
path "secret/{{identity.entity.id}}/*" {
	capabilities          = ["read"]
	required_parameters   = ["a"]
	mfa_methods           = ["mfa-a"]
	subscribe_event_types = ["kv*"]
}`,
		"templated-b": `
path "secret/{{identity.entity.id}}/*" {
	capabilities          = ["update"]
	required_parameters   = ["b"]
	mfa_methods           = ["mfa-b"]
	subscribe_event_types = ["sys*"]
}`,
	}
	names := make([]string, 0, len(policies))
	for name, raw := range policies {
		p, err := ParseACLPolicy(namespace.RootNamespace, raw)
		require.NoError(t, err)
		p.Name = name
		require.NoError(t, ps.SetPolicy(ctx, p))
		names = append(names, name)
	}

	const (
		workers    = 16
		iterations = 200
	)
	entity := &identity.Entity{ID: "entity-id"}
	policyNames := map[string][]string{namespace.RootNamespaceID: names}

	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range iterations {
				acl, err := ps.ACL(ctx, entity, policyNames)
				if !assert.NoError(t, err) {
					return
				}
				res := acl.AllowOperation(ctx, &logical.Request{
					Operation: logical.UpdateOperation,
					Path:      "secret/entity-id/data",
				}, true)
				if !assert.Equal(t, ReadCapabilityInt|UpdateCapabilityInt, res.CapabilitiesBitmap) {
					return
				}
				res = acl.AllowOperation(ctx, &logical.Request{
					Operation: logical.UpdateOperation,
					Path:      "secret/entity-id/data",
					Data:      map[string]interface{}{"a": "1", "b": "2"},
				}, false)
				if !assert.True(t, res.Allowed) || !assert.ElementsMatch(t, []string{"mfa-a", "mfa-b"}, res.MFAMethods) {
					return
				}
			}
		}()
	}
	wg.Wait()
}
