// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_CheckVersionViolationsReq_Run verifies that each of the seven rules is
// correctly detected, and that a clean input produces no violations. The bugs
// in the original shell script (double-counted GA versions, wrong variable
// names for major.minor collisions and RC-patch-level violations) are each
// covered by a dedicated case.
func Test_CheckVersionViolationsReq_Run(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	for name, tc := range map[string]struct {
		versions      []string
		wantViolation bool
		check         func(t *testing.T, res *CheckVersionViolationsRes)
	}{
		"nil request returns error": {
			versions: nil,
			check: func(t *testing.T, res *CheckVersionViolationsRes) {
				t.Helper()
				// handled separately below
			},
		},
		"clean release set has no violations": {
			versions:      []string{"1.21.0-rc1", "1.20.4", "1.19.10", "1.18.15"},
			wantViolation: false,
			check: func(t *testing.T, res *CheckVersionViolationsRes) {
				t.Helper()
				require.Empty(t, res.InvalidVersionFormat)
				require.Empty(t, res.MajorMinorDuplicates)
				require.Empty(t, res.GaOrRcNotTopVersion)
				require.Empty(t, res.GaAndRcMixedInSet)
				require.Empty(t, res.RCWithInvalidPatchLevel)
			},
		},
		"rule 1: two versions with same major.minor detected": {
			versions:      []string{"1.21.1", "1.21.0", "1.20.4"},
			wantViolation: true,
			check: func(t *testing.T, res *CheckVersionViolationsRes) {
				t.Helper()
				require.Contains(t, res.MajorMinorDuplicates, "1.21")
			},
		},
		"rule 2: two versions with patch level 0 detected": {
			versions:      []string{"2.0.0", "1.20.4", "1.19.0"},
			wantViolation: true,
			check: func(t *testing.T, res *CheckVersionViolationsRes) {
				t.Helper()
				require.Len(t, res.PatchZeroHits, 2)
			},
		},
		"rule 3: GA version not at top of list detected": {
			versions:      []string{"1.21.1", "1.20.0", "1.19.10"},
			wantViolation: true,
			check: func(t *testing.T, res *CheckVersionViolationsRes) {
				t.Helper()
				require.Contains(t, res.GaOrRcNotTopVersion, "1.20.0")
			},
		},
		"rule 3: RC version not at top of list detected": {
			versions:      []string{"1.21.1", "1.20.0-rc1", "1.19.10"},
			wantViolation: true,
			check: func(t *testing.T, res *CheckVersionViolationsRes) {
				t.Helper()
				require.Contains(t, res.GaOrRcNotTopVersion, "1.20.0-rc1")
			},
		},
		"rule 4: GA and RC coexisting detected": {
			versions:      []string{"1.21.0", "1.20.0-rc1", "1.19.10"},
			wantViolation: true,
			check: func(t *testing.T, res *CheckVersionViolationsRes) {
				t.Helper()
				require.NotEmpty(t, res.GaAndRcMixedInSet)
				require.Contains(t, res.GaAndRcMixedInSet, "1.21.0")
				require.Contains(t, res.GaAndRcMixedInSet, "1.20.0-rc1")
			},
		},
		"rule 5: RC with non-zero patch level detected": {
			versions:      []string{"1.21.1-rc1", "1.20.4", "1.19.10"},
			wantViolation: true,
			check: func(t *testing.T, res *CheckVersionViolationsRes) {
				t.Helper()
				require.Contains(t, res.RCWithInvalidPatchLevel, "1.21.1-rc1")
			},
		},
		"rule 6: invalid format detected": {
			versions:      []string{"1.21.0", "not-a-version", "1.19.10"},
			wantViolation: true,
			check: func(t *testing.T, res *CheckVersionViolationsRes) {
				t.Helper()
				require.Contains(t, res.InvalidVersionFormat, "not-a-version")
			},
		},
		"rule 7: two GA versions detected": {
			versions:      []string{"2.0.0", "1.20.0", "1.19.10"},
			wantViolation: true,
			check: func(t *testing.T, res *CheckVersionViolationsRes) {
				t.Helper()
				require.Len(t, res.GAVersions, 2)
			},
		},
		// Regression: original script double-counted GA versions (blocks 4 and 7
		// both appended). A single GA version must not trigger the >1 violation.
		"single GA version is not a violation (regression: no double-count)": {
			versions:      []string{"1.21.0", "1.20.4", "1.19.10"},
			wantViolation: false,
			check: func(t *testing.T, res *CheckVersionViolationsRes) {
				t.Helper()
				require.Len(t, res.GAVersions, 1)
			},
		},
		// Regression: original script checked wrong variable name for major.minor
		// collisions (major_minor_collisions vs major_minor_with_duplicates).
		"major.minor collision is detected (regression: correct variable name)": {
			versions:      []string{"1.21.2", "1.21.1", "1.20.4"},
			wantViolation: true,
			check: func(t *testing.T, res *CheckVersionViolationsRes) {
				t.Helper()
				require.Contains(t, res.MajorMinorDuplicates, "1.21")
			},
		},
		// Regression: original script checked wrong variable name for RC patch
		// level (invalid_rc_patch_level vs release_candidate_with_invalid_patch_level).
		"RC with invalid patch level is detected (regression: correct variable name)": {
			versions:      []string{"1.21.2-rc1", "1.20.4"},
			wantViolation: true,
			check: func(t *testing.T, res *CheckVersionViolationsRes) {
				t.Helper()
				require.Contains(t, res.RCWithInvalidPatchLevel, "1.21.2-rc1")
			},
		},
		// Regression: original script RC regex only matched -rc1 through -rc9.
		"multi-digit RC number is valid format": {
			versions:      []string{"1.21.0-rc10", "1.20.4"},
			wantViolation: false,
			check: func(t *testing.T, res *CheckVersionViolationsRes) {
				t.Helper()
				require.Empty(t, res.InvalidVersionFormat)
				require.Contains(t, res.RCVersions, "1.21.0-rc10")
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if name == "nil request returns error" {
				var req *CheckVersionViolationsReq
				_, err := req.Run(ctx)
				require.Error(t, err)
				return
			}

			req := &CheckVersionViolationsReq{Versions: tc.versions}
			res, err := req.Run(ctx)
			require.NoError(t, err)
			require.Equal(t, tc.wantViolation, res.HasViolations)
			tc.check(t, res)
		})
	}
}

// Test_CheckVersionViolationsRes_ToJSON verifies that the response serializes
// with the expected key names that the workflow's compile step reads.
func Test_CheckVersionViolationsRes_ToJSON(t *testing.T) {
	t.Parallel()

	res := &CheckVersionViolationsRes{
		HasViolations:        true,
		InvalidVersionFormat: []string{"bad"},
		MajorMinorDuplicates: []string{"1.21"},
	}

	b, err := res.ToJSON()
	require.NoError(t, err)

	var parsed map[string]any
	require.NoError(t, json.Unmarshal(b, &parsed))
	require.Contains(t, parsed, "has_violations")
	require.Contains(t, parsed, "invalid_version_format_violations")
	require.Contains(t, parsed, "major_minor_with_duplicates_violations")
}

// Test_CheckVersionViolationsRes_ToGithubOutput verifies that the output keys
// exactly match what the compile-version-rule-violations workflow step expects.
func Test_CheckVersionViolationsRes_ToGithubOutput(t *testing.T) {
	t.Parallel()

	res := &CheckVersionViolationsRes{
		HasViolations:        true,
		InvalidVersionFormat: []string{"bad"},
		MajorMinorDuplicates: []string{"1.21"},
		RCVersions:           []string{"1.21.0-rc1", "1.21.0-rc2"},
		GAVersions:           []string{"1.20.0"},
	}

	b, err := res.ToGithubOutput()
	require.NoError(t, err)
	output := string(b)

	require.Contains(t, output, "has_violations=true\n")
	require.Contains(t, output, "invalid_version_format_violations=bad\n")
	require.Contains(t, output, "major_minor_with_duplicates_violations=1.21\n")
	require.Contains(t, output, "rc_version_violations=1.21.0-rc1,1.21.0-rc2\n")
	require.Contains(t, output, "ga_version_violations=1.20.0\n")
}

// Test_CheckVersionViolationsRes_ToGithubOutput_NoViolations verifies that all
// keys are still present with empty values when there are no violations, so
// downstream workflow steps can reference them unconditionally.
func Test_CheckVersionViolationsRes_ToGithubOutput_NoViolations(t *testing.T) {
	t.Parallel()

	res := &CheckVersionViolationsRes{
		HasViolations: false,
		RCVersions:    []string{},
		GAVersions:    []string{},
	}

	b, err := res.ToGithubOutput()
	require.NoError(t, err)
	output := string(b)

	require.Contains(t, output, "has_violations=false\n")
	require.Contains(t, output, "invalid_version_format_violations=\n")
	require.Contains(t, output, "rc_version_violations=\n")
	require.Contains(t, output, "ga_version_violations=\n")
}
