// Copyright IBM Corp. 2026
// SPDX-License-Identifier: MPL-2.0

package tokenutil

import (
	"testing"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/stretchr/testify/require"
)

func tokenPolicyFieldData(raw map[string]interface{}) *framework.FieldData {
	schema := map[string]*framework.FieldSchema{
		"policies": {
			Type: framework.TypeCommaStringSlice,
		},
	}
	for key, value := range TokenFields() {
		schema[key] = value
	}

	return &framework.FieldData{
		Raw:    raw,
		Schema: schema,
	}
}

// TestTokenParams_ParseTokenFields_TokenPoliciesValidation verifies that token
// policy names are canonicalized and validated during token field parsing.
func TestTokenParams_ParseTokenFields_TokenPoliciesValidation(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		rawPolicies   interface{}
		expected      []string
		expectedError string
	}{
		"accepts valid policy names": {
			rawPolicies: []string{" Team-Read ", "default"},
			expected:    []string{"team-read", "default"},
		},
		"accepts slash-containing policy names": {
			rawPolicies: []string{"team/read", "default"},
			expected:    []string{"team/read", "default"},
		},
		"accepts backslash-containing policy names": {
			rawPolicies: []string{`team\read`, "default"},
			expected:    []string{`team\read`, "default"},
		},
		"accepts embedded dot dot policy names": {
			rawPolicies: []string{"policy..backup", "default"},
			expected:    []string{"policy..backup", "default"},
		},
		"rejects malformed traversal policy names": {
			rawPolicies:   []string{"../team/dev"},
			expectedError: "policy name cannot contain '..'",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tokenParams := &TokenParams{}
			data := tokenPolicyFieldData(map[string]interface{}{
				"token_policies": tc.rawPolicies,
			})

			err := tokenParams.ParseTokenFields(nil, data)
			if tc.expectedError != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.expectedError)
				return
			}

			require.NoError(t, err)
			require.ElementsMatch(t, tc.expected, tokenParams.TokenPolicies)
		})
	}
}

// TestUpgradeValue_TokenPoliciesValidation verifies that upgrades from the
// deprecated policies field enforce policy-name validation for token_policies.
func TestUpgradeValue_TokenPoliciesValidation(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		rawPolicies   interface{}
		expected      []string
		expectedError string
	}{
		"accepts valid policies": {
			rawPolicies: "team/read,default",
			expected:    []string{"team/read", "default"},
		},
		"accepts backslash policy": {
			rawPolicies: `team\ops,default`,
			expected:    []string{`team\ops`, "default"},
		},
		"accepts embedded dot dot policy": {
			rawPolicies: "policy..backup,default",
			expected:    []string{"policy..backup", "default"},
		},
		"rejects malformed traversal policies": {
			rawPolicies:   "../team/ops",
			expectedError: "policy name cannot contain '..'",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			data := tokenPolicyFieldData(map[string]interface{}{
				"policies": tc.rawPolicies,
			})

			var oldValue []string
			var newValue []string
			err := UpgradeValue(data, "policies", "token_policies", &oldValue, &newValue)

			if tc.expectedError != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.expectedError)
				return
			}

			require.NoError(t, err)
			require.ElementsMatch(t, tc.expected, oldValue)
			require.ElementsMatch(t, tc.expected, newValue)
		})
	}
}
