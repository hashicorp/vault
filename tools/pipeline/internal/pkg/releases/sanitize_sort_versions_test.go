// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_SanitizeSortVersionsReq_Run verifies that the sanitize-sort operation
// correctly normalizes, categorizes, and orders a mixed set of version strings.
func Test_SanitizeSortVersionsReq_Run(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	for name, tc := range map[string]struct {
		req            *SanitizeSortVersionsReq
		wantSorted     []string
		wantLTS        []string
		wantInvalid    []string
		wantDuplicates []string
		wantErr        bool
	}{
		"nil request returns error": {
			req:     nil,
			wantErr: true,
		},
		"empty versions returns empty result": {
			req:         &SanitizeSortVersionsReq{Descending: true},
			wantSorted:  []string{},
			wantLTS:     []string{},
			wantInvalid: []string{},
		},
		"plain versions sorted descending": {
			req: &SanitizeSortVersionsReq{
				Descending: true,
				Versions:   []string{"1.19.10", "1.21.0", "1.20.4"},
			},
			wantSorted:  []string{"1.21.0", "1.20.4", "1.19.10"},
			wantLTS:     []string{},
			wantInvalid: []string{},
		},
		"plain versions sorted ascending": {
			req: &SanitizeSortVersionsReq{
				Descending: false,
				Versions:   []string{"1.19.10", "1.21.0", "1.20.4"},
			},
			wantSorted:  []string{"1.19.10", "1.20.4", "1.21.0"},
			wantLTS:     []string{},
			wantInvalid: []string{},
		},
		"+ent suffix is preserved": {
			req: &SanitizeSortVersionsReq{
				Descending: true,
				Versions:   []string{"1.19.10+ent", "1.21.0+ent"},
			},
			wantSorted:  []string{"1.21.0+ent", "1.19.10+ent"},
			wantLTS:     []string{},
			wantInvalid: []string{},
		},
		"rc suffix is preserved": {
			req: &SanitizeSortVersionsReq{
				Descending: true,
				Versions:   []string{"1.21.0-rc1", "1.20.4"},
			},
			wantSorted:  []string{"1.21.0-rc1", "1.20.4"},
			wantLTS:     []string{},
			wantInvalid: []string{},
		},
		"rc suffix with +ent is preserved": {
			req: &SanitizeSortVersionsReq{
				Descending: true,
				Versions:   []string{"1.21.0-rc1+ent", "1.20.4+ent"},
			},
			wantSorted:  []string{"1.21.0-rc1+ent", "1.20.4+ent"},
			wantLTS:     []string{},
			wantInvalid: []string{},
		},
		"-lts suffix stripped from sorted list but base tracked in lts": {
			req: &SanitizeSortVersionsReq{
				Descending: true,
				Versions:   []string{"1.20.4-lts", "1.19.10"},
			},
			wantSorted:  []string{"1.20.4", "1.19.10"},
			wantLTS:     []string{"1.20.4"},
			wantInvalid: []string{},
		},
		"-ce suffix is stripped from sorted list": {
			req: &SanitizeSortVersionsReq{
				Descending: true,
				Versions:   []string{"1.20.4-ce", "1.19.10"},
			},
			wantSorted:  []string{"1.20.4", "1.19.10"},
			wantLTS:     []string{},
			wantInvalid: []string{},
		},
		"invalid version is separated out": {
			req: &SanitizeSortVersionsReq{
				Descending: true,
				Versions:   []string{"1.20.4", "not-a-version", "1.19.10"},
			},
			wantSorted:  []string{"1.20.4", "1.19.10"},
			wantLTS:     []string{},
			wantInvalid: []string{"not-a-version"},
		},
		"leading zeros in version components are invalid": {
			req: &SanitizeSortVersionsReq{
				Descending: true,
				Versions:   []string{"1.20.4", "01.20.4"},
			},
			wantSorted:  []string{"1.20.4"},
			wantLTS:     []string{},
			wantInvalid: []string{"01.20.4"},
		},
		"multi-digit rc number is valid": {
			req: &SanitizeSortVersionsReq{
				Descending: true,
				Versions:   []string{"1.21.0-rc10", "1.20.4"},
			},
			wantSorted:  []string{"1.21.0-rc10", "1.20.4"},
			wantLTS:     []string{},
			wantInvalid: []string{},
		},
		"mixed valid and invalid with lts and ent": {
			req: &SanitizeSortVersionsReq{
				Descending: true,
				Versions:   []string{"2.0.3+ent", "1.21.0-rc1+ent", "1.20.4-lts", "bad", "1.19.10+ent"},
			},
			wantSorted:  []string{"2.0.3+ent", "1.21.0-rc1+ent", "1.20.4", "1.19.10+ent"},
			wantLTS:     []string{"1.20.4"},
			wantInvalid: []string{"bad"},
		},
		"duplicate versions are deduplicated and reported": {
			req: &SanitizeSortVersionsReq{
				Descending: true,
				Versions:   []string{"1.21.0", "1.20.4", "1.21.0", "1.19.10", "1.20.4"},
			},
			wantSorted:     []string{"1.21.0", "1.20.4", "1.19.10"},
			wantLTS:        []string{},
			wantInvalid:    []string{},
			wantDuplicates: []string{"1.21.0", "1.20.4"},
		},
		"duplicate +ent versions are deduplicated": {
			req: &SanitizeSortVersionsReq{
				Descending: true,
				Versions:   []string{"1.21.0+ent", "1.20.4+ent", "1.21.0+ent"},
			},
			wantSorted:     []string{"1.21.0+ent", "1.20.4+ent"},
			wantLTS:        []string{},
			wantInvalid:    []string{},
			wantDuplicates: []string{"1.21.0+ent"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			res, err := tc.req.Run(ctx)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.wantSorted, res.SortedVersions)
			require.Equal(t, tc.wantLTS, res.LTSVersions)
			require.Equal(t, tc.wantInvalid, res.InvalidVersions)
			if tc.wantDuplicates != nil {
				require.Equal(t, tc.wantDuplicates, res.DuplicateVersions)
			} else {
				require.Empty(t, res.DuplicateVersions)
			}
		})
	}
}

// Test_SanitizeSortVersionsRes_ToJSON verifies that the response serializes
// correctly to JSON with the expected keys.
func Test_SanitizeSortVersionsRes_ToJSON(t *testing.T) {
	t.Parallel()

	res := &SanitizeSortVersionsRes{
		SortedVersions:  []string{"1.21.0", "1.20.4"},
		LTSVersions:     []string{"1.20.4"},
		InvalidVersions: []string{"bad"},
	}

	b, err := res.ToJSON()
	require.NoError(t, err)

	var parsed map[string]any
	require.NoError(t, json.Unmarshal(b, &parsed))
	require.Contains(t, parsed, "sorted_versions")
	require.Contains(t, parsed, "lts_versions")
	require.Contains(t, parsed, "invalid_versions")
	require.Contains(t, parsed, "duplicate_versions")
}

// Test_SanitizeSortVersionsRes_ToGithubOutput verifies that the output written
// to $GITHUB_OUTPUT contains the expected key=value lines for all three fields.
func Test_SanitizeSortVersionsRes_ToGithubOutput(t *testing.T) {
	t.Parallel()

	res := &SanitizeSortVersionsRes{
		SortedVersions:  []string{"1.21.0+ent", "1.20.4+ent"},
		LTSVersions:     []string{"1.20.4"},
		InvalidVersions: []string{"bad"},
	}

	b, err := res.ToGithubOutput()
	require.NoError(t, err)
	output := string(b)

	require.Contains(t, output, "sorted_versions=1.21.0+ent 1.20.4+ent\n")
	require.Contains(t, output, "lts_versions=1.20.4\n")
	require.Contains(t, output, "invalid_versions=bad\n")
	require.Contains(t, output, "duplicate_versions=\n")
}

// Test_SanitizeSortVersionsRes_ToGithubOutput_Empty verifies that empty slices
// produce empty values rather than omitting the keys entirely, so downstream
// workflow steps can safely reference the outputs without conditional guards.
func Test_SanitizeSortVersionsRes_ToGithubOutput_Empty(t *testing.T) {
	t.Parallel()

	res := &SanitizeSortVersionsRes{}

	b, err := res.ToGithubOutput()
	require.NoError(t, err)
	output := string(b)

	require.Contains(t, output, "sorted_versions=\n")
	require.Contains(t, output, "lts_versions=\n")
	require.Contains(t, output, "invalid_versions=\n")
	require.Contains(t, output, "duplicate_versions=\n")
}
