// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/internalshared/namespace"
)

// BenchmarkPolicyStore_ACLTemplated measures building an ACL for a token
// whose policies are cached and use identity templating, which happens on
// every request made with such a token.
func BenchmarkPolicyStore_ACLTemplated(b *testing.B) {
	for _, tc := range []struct{ policies, pathsPerPolicy int }{
		{policies: 1, pathsPerPolicy: 1},
		{policies: 1, pathsPerPolicy: 10},
		{policies: 3, pathsPerPolicy: 10},
	} {
		b.Run(fmt.Sprintf("policies=%d/paths=%d", tc.policies, tc.pathsPerPolicy), func(b *testing.B) {
			ps := NewTestCluster(b, &CoreConfig{}, nil).Cores[0].Core.policyStore
			ctx := namespace.RootContext(context.Background())

			names := make([]string, tc.policies)
			for i := range names {
				var raw strings.Builder
				for j := range tc.pathsPerPolicy {
					fmt.Fprintf(&raw, "path \"secret/%d/%d/{{identity.entity.id}}/{{identity.entity.metadata.team}}/*\" {\n\tcapabilities = [\"read\", \"list\"]\n}\n", i, j)
				}
				p, err := ParseACLPolicy(namespace.RootNamespace, raw.String())
				if err != nil {
					b.Fatal(err)
				}
				names[i] = fmt.Sprintf("templated-%d", i)
				p.Name = names[i]
				if err := ps.SetPolicy(ctx, p); err != nil {
					b.Fatal(err)
				}
			}

			entity := &identity.Entity{
				ID:          "entity-id",
				NamespaceID: namespace.RootNamespaceID,
				Metadata:    map[string]string{"team": "eng"},
			}
			policyNames := map[string][]string{namespace.RootNamespaceID: names}

			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				if _, err := ps.ACL(ctx, entity, policyNames); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
