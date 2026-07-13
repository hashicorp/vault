// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

const testVersionConfig = `
schema = 1

active_versions {
	version "1.19.x" {
		ce_active = true
		lts       = true
	}
	version "1.18.x" {
		ce_active = true
	}
	version "1.17.x" {
		ce_active = false
	}
	version "1.16.x" {
		ce_active = false
		lts       = true
	}
}
`

func TestListActiveVersionReq_unmarshalConfig(t *testing.T) {
	t.Parallel()

	versionsConfig, err := DecodeBytes([]byte(testVersionConfig))
	require.NoError(t, err)
	require.EqualValues(t, &VersionsConfig{
		Schema: 1,
		ActiveVersion: &ActiveVersion{
			Versions: map[string]*Version{
				"1.19.x": {CEActive: true, LTS: true},
				"1.18.x": {CEActive: true, LTS: false},
				"1.17.x": {CEActive: false, LTS: false},
				"1.16.x": {CEActive: false, LTS: true},
			},
		},
	}, versionsConfig)
}

func TestEnterpriseReleaseBranchForVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		version  string
		expected string
	}{
		{
			name:     "standard version",
			version:  "1.19.x",
			expected: "release/1.19.x+ent",
		},
		{
			name:     "version with patch",
			version:  "1.19.5",
			expected: "release/1.19.5+ent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := enterpriseReleaseBranchForVersion(tt.version)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestCEReleaseBranchForVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		version  string
		prefix   string
		expected string
	}{
		{
			name:     "no prefix",
			version:  "1.19.x",
			prefix:   "",
			expected: "release/1.19.x",
		},
		{
			name:     "with ce prefix",
			version:  "1.19.x",
			prefix:   "ce",
			expected: "ce/release/1.19.x",
		},
		{
			name:     "with custom prefix",
			version:  "2.0.x",
			prefix:   "community",
			expected: "community/release/2.0.x",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ceReleaseBranchForVersion(tt.version, tt.prefix)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestListActiveVersionsRes_ToJSON(t *testing.T) {
	t.Parallel()

	versionsConfig, err := DecodeBytes([]byte(testVersionConfig))
	require.NoError(t, err)

	res := &ListActiveVersionsRes{
		VersionsConfig: versionsConfig,
	}

	tests := []struct {
		name     string
		cePrefix string
	}{
		{
			name:     "without prefix",
			cePrefix: "",
		},
		{
			name:     "with ce prefix",
			cePrefix: "ce",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := res.ToJSON(tt.cePrefix)
			require.NoError(t, err)
			require.NotEmpty(t, jsonBytes)

			// Verify it's valid JSON by unmarshaling
			var output ListActiveVersionsJSONOutput
			err = json.Unmarshal(jsonBytes, &output)
			require.NoError(t, err)

			// Verify structure
			require.NotNil(t, output.VersionsConfig)
			require.Len(t, output.Versions, 4)

			// Check that all versions have enterprise branches
			for _, v := range output.Versions {
				require.NotEmpty(t, v.EnterpriseBranch)
				require.Contains(t, v.EnterpriseBranch, "+ent")

				// CE active versions should have CE branch
				if v.CEActive {
					require.NotEmpty(t, v.CEBranch)
					if tt.cePrefix != "" {
						require.Contains(t, v.CEBranch, tt.cePrefix+"/")
					}
				} else {
					require.Empty(t, v.CEBranch)
				}
			}
		})
	}
}

func TestListActiveVersionsRes_ToGithubOutput(t *testing.T) {
	t.Parallel()

	versionsConfig, err := DecodeBytes([]byte(testVersionConfig))
	require.NoError(t, err)

	res := &ListActiveVersionsRes{
		VersionsConfig: versionsConfig,
	}

	tests := []struct {
		name        string
		includeMain bool
		cePrefix    string
	}{
		{
			name:        "no flags",
			includeMain: false,
			cePrefix:    "",
		},
		{
			name:        "with include-main",
			includeMain: true,
			cePrefix:    "",
		},
		{
			name:        "with ce prefix",
			includeMain: false,
			cePrefix:    "ce",
		},
		{
			name:        "with both flags",
			includeMain: true,
			cePrefix:    "ce",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := res.ToGithubOutput(tt.includeMain, tt.cePrefix)
			require.NoError(t, err)
			require.NotEmpty(t, jsonBytes)

			// Verify it's valid JSON by unmarshaling
			var output ListActiveVersionsGithubOutput
			err = json.Unmarshal(jsonBytes, &output)
			require.NoError(t, err)

			// Verify basic structure
			require.NotNil(t, output.VersionsConfig)
			require.Len(t, output.Versions, 4)
			require.Len(t, output.CEActiveVersions, 2) // 1.19.x and 1.18.x
			require.Len(t, output.LTSVersions, 2)      // 1.19.x and 1.16.x

			// Verify active_branches
			require.NotEmpty(t, output.ActiveBranches)

			expectedBranchCount := 4
			if tt.includeMain {
				expectedBranchCount++
			}
			require.Len(t, output.ActiveBranches, expectedBranchCount)
			for _, branch := range output.ActiveBranches {
				if branch != "main" {
					require.Contains(t, branch, "+ent")
				}
			}

			// Verify ce_active_branches
			require.NotEmpty(t, output.CEActiveBranches)
			expectedCECount := 2 // 1.19.x and 1.18.x
			if tt.includeMain {
				expectedCECount++ // + main or ce/main
			}
			require.Len(t, output.CEActiveBranches, expectedCECount)
			for _, branch := range output.CEActiveBranches {
				if tt.cePrefix != "" && branch != "main" {
					require.Contains(t, branch, tt.cePrefix+"/")
				}
				require.NotContains(t, branch, "+ent")
			}

			// Verify lts_active_branches
			require.NotEmpty(t, output.LTSActiveBranches)
			require.Len(t, output.LTSActiveBranches, 2) // 1.19.x and 1.16.x
			for _, branch := range output.LTSActiveBranches {
				require.Contains(t, branch, "+ent")
			}

			// Verify all_active_branches
			require.NotEmpty(t, output.AllActiveBranches)
			expectedAllCount := 6 // 4 ent + 2 ce
			if tt.includeMain {
				expectedAllCount += 2 // main + ce/main (or just main if no prefix)
				if tt.cePrefix == "" {
					expectedAllCount-- // only one main if no prefix
				}
			}
			require.Len(t, output.AllActiveBranches, expectedAllCount)

			// Verify main is included when flag is set
			if tt.includeMain {
				require.Contains(t, output.ActiveBranches, "main")
				if tt.cePrefix != "" {
					require.Contains(t, output.CEActiveBranches, tt.cePrefix+"/main")
					require.Contains(t, output.AllActiveBranches, "main")
					require.Contains(t, output.AllActiveBranches, tt.cePrefix+"/main")
				} else {
					require.Contains(t, output.CEActiveBranches, "main")
					require.Contains(t, output.AllActiveBranches, "main")
				}
			} else {
				require.NotContains(t, output.ActiveBranches, "main")
				require.NotContains(t, output.CEActiveBranches, "main")
				if tt.cePrefix != "" {
					require.NotContains(t, output.CEActiveBranches, tt.cePrefix+"/main")
				}
			}

			// Verify sorting (all arrays should be sorted)
			require.True(t, slices.IsSorted(output.Versions))
			require.True(t, slices.IsSorted(output.CEActiveVersions))
			require.True(t, slices.IsSorted(output.LTSVersions))
			require.True(t, slices.IsSorted(output.ActiveBranches))
			require.True(t, slices.IsSorted(output.CEActiveBranches))
			require.True(t, slices.IsSorted(output.LTSActiveBranches))
			require.True(t, slices.IsSorted(output.AllActiveBranches))
		})
	}
}
