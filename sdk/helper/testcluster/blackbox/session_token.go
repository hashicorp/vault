// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

type TokenOptions struct {
	Policies    []string
	TTL         string
	Renewable   bool
	NoParent    bool
	DisplayName string
}

// MustCreateToken generates a new token with specific properties.
func (s *Session) MustCreateToken(opts TokenOptions) string {
	s.t.Helper()

	payload := map[string]any{
		"policies":     opts.Policies,
		"ttl":          opts.TTL,
		"renewable":    opts.Renewable,
		"no_parent":    opts.NoParent,
		"display_name": opts.DisplayName,
	}

	// Use auth/token/create for child tokens, or auth/token/create-orphan
	path := "auth/token/create"
	if opts.NoParent {
		path = "auth/token/create-orphan"
	}

	secret := s.MustWrite(path, payload)
	if secret.Auth == nil {
		s.t.Fatal("Token creation response missing Auth data")
	}

	return secret.Auth.ClientToken
}

// AssertTokenIsValid checks that a token works and (optionally) has specific policies.
func (s *Session) AssertTokenIsValid(token string, expectedPolicies ...string) {
	s.t.Helper()

	if token == "" {
		s.t.Fatal("token is empty")
	}

	clonedConfig := s.Client.CloneConfig()
	tempClient, err := api.NewClient(clonedConfig)
	require.NoError(s.t, err)

	tempClient.SetToken(token)
	tempClient.SetNamespace(s.Namespace)

	secret, err := tempClient.Auth().Token().LookupSelf()
	require.NoError(s.t, err)

	if len(expectedPolicies) == 0 {
		return
	}

	rawPolicies, ok := secret.Data["policies"].([]any)
	if !ok {
		s.t.Fatalf("token does not contain any policies")
	}

	actualPolicies := make(map[string]struct{})
	for _, p := range rawPolicies {
		if val, ok := p.(string); ok {
			actualPolicies[val] = struct{}{}
		}
	}

	var missing []string
	for _, expected := range expectedPolicies {
		if _, ok := actualPolicies[expected]; !ok {
			missing = append(missing, expected)
		}
	}

	if len(missing) > 0 {
		allActual := make([]string, 0, len(actualPolicies))
		for k := range actualPolicies {
			allActual = append(allActual, k)
		}

		s.t.Fatalf("token policy mismatch.\n\tmissing: %v\n\tactual: %v", missing, allActual)
	}
}
