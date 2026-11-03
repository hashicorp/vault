// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/stretchr/testify/require"
)

// TestSyncGoModReq_Validate_EmptySourceBranch verifies that validate returns an
// error when a branch is empty or both are the same, and when paths are given
// together with workspace modules.
func TestSyncGoModReq_Validate_EmptySourceBranch(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		req        *SyncGoModReq
		shouldFail bool
	}{
		"empty source branch": {
			&SyncGoModReq{
				SourceBranch: "",
				DestBranch:   "main",
				Paths:        []string{"go.mod"},
			},
			true,
		},
		"empty dest branch": {
			&SyncGoModReq{
				SourceBranch: "feature",
				DestBranch:   "",
				Paths:        []string{"go.mod"},
			},
			true,
		},
		"same branches": {
			&SyncGoModReq{
				SourceBranch: "main",
				DestBranch:   "main",
				Paths:        []string{"go.mod"},
			},
			true,
		},
		"empty paths default to go.mod": {
			&SyncGoModReq{
				SourceBranch: "feature",
				DestBranch:   "main",
				Paths:        []string{},
			},
			false,
		},
		"workspace modules": {
			&SyncGoModReq{
				SourceBranch:     "feature",
				DestBranch:       "main",
				WorkspaceModules: true,
			},
			false,
		},
		"workspace modules with paths": {
			&SyncGoModReq{
				SourceBranch:     "feature",
				DestBranch:       "main",
				Paths:            []string{"go.mod"},
				WorkspaceModules: true,
			},
			true,
		},
		"valid": {
			&SyncGoModReq{
				SourceBranch: "feature",
				DestBranch:   "main",
				Paths:        []string{"go.mod"},
			},
			false,
		},
		"valid with commit": {
			&SyncGoModReq{
				SourceBranch: "feature",
				DestBranch:   "main",
				Paths:        []string{"go.mod"},
				Commit:       true,
			},
			false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := test.req.validate()
			if test.shouldFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestBranchName verifies that branchName produces the correct intermediate
// branch name for combinations of ticket presence and branch names containing
// slashes.
func TestBranchName(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		source string
		dest   string
		ticket string
		want   string
	}{
		"no ticket": {
			source: "main", dest: "ce/main", ticket: "",
			want: "main-into-ce-main",
		},
		"with ticket": {
			source: "main", dest: "ce/main", ticket: "VAULT-1234",
			want: "VAULT-1234-main-into-ce-main",
		},
		"slashes in source and dest": {
			source: "release/1.18.x", dest: "ce/release/1.18.x", ticket: "VAULT-5678",
			want: "VAULT-5678-release-1.18.x-into-ce-release-1.18.x",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			req := &SyncGoModReq{
				SourceBranch: tc.source,
				DestBranch:   tc.dest,
				Ticket:       tc.ticket,
			}
			require.Equal(t, tc.want, req.branchName())
		})
	}
}

// TestCommitSubject verifies the commit subject depending on request params.
func TestCommitSubject(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		req      *SyncGoModReq
		expected string
	}{
		"no ticket": {
			&SyncGoModReq{
				SourceBranch: "main",
				DestBranch:   "ce/main",
			},
			"go: sync go.mod from main to ce/main",
		},
		"with ticket": {
			&SyncGoModReq{
				SourceBranch: "main",
				DestBranch:   "ce/main",
				Ticket:       "VAULT-1234",
			},
			"[VAULT-1234] go: sync go.mod from main to ce/main",
		},
		"custom": {
			&SyncGoModReq{
				SourceBranch:  "main",
				DestBranch:    "ce/main",
				Ticket:        "VAULT-1234",
				CommitMessage: "custom commit message",
			},
			"custom commit message",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, test.expected, test.req.commitSubject())
		})
	}
}

// TestSyncGoModRes_ToString_WithResults verifies that ToString produces a
// non-empty table containing path, changes, tidy, and error columns.
func TestSyncGoModRes_ToString_WithResults(t *testing.T) {
	t.Parallel()

	res := &SyncGoModRes{
		SourceBranch: "feature",
		DestBranch:   "main",
		NewBranch:    "feature-into-main",
		Results: []*golang.GoModPathSyncResult{
			{
				SyncModRes: &golang.SyncModRes{
					Path: "go.mod",
					Changes: []*golang.SyncChange{
						{
							Directive:  golang.DirectiveRequire,
							Module:     "github.com/foo/bar",
							OldVersion: "v1.0.0",
							NewVersion: "v1.2.0",
							Action:     golang.SyncActionUpdated,
						},
					},
				},
				TidyRan: true,
			},
		},
	}

	out := res.ToString()
	require.NotEmpty(t, out)
	require.Contains(t, out, "go.mod")
	require.Contains(t, out, "1")   // changes count
	require.Contains(t, out, "yes") // tidy ran
}

// TestSyncGoModRes_ToString_WithErrors verifies that ToString includes both the
// sync error and tidy error joined with "; " in the error column.
func TestSyncGoModRes_ToString_WithErrors(t *testing.T) {
	t.Parallel()

	syncRes := &golang.SyncModRes{
		Path:  "go.mod",
		Error: "sync failed",
	}
	syncRes.Err = errors.New(syncRes.Error)

	res := &SyncGoModRes{
		SourceBranch: "feature",
		DestBranch:   "main",
		Results: []*golang.GoModPathSyncResult{
			{
				SyncModRes: syncRes,
				TidyRan:    true,
				TidyError:  "tidy failed",
			},
		},
	}

	out := res.ToString()
	require.NotEmpty(t, out)
	require.Contains(t, out, "go.mod")
	require.Contains(t, out, "sync failed")
	require.Contains(t, out, "tidy failed")
}

// TestSyncGoModRes_ToString_NilRes verifies that calling ToString on a nil
// pointer returns an empty string without panicking.
func TestSyncGoModRes_ToString_NilRes(t *testing.T) {
	t.Parallel()

	var res *SyncGoModRes
	require.Empty(t, res.ToString(), "nil receiver must return empty string")
}

// TestSyncGoModRes_ToMarkdown_WithResults verifies that ToMarkdown produces a
// pipe-delimited markdown table when results are present.
func TestSyncGoModRes_ToMarkdown_WithResults(t *testing.T) {
	t.Parallel()

	res := &SyncGoModRes{
		SourceBranch: "feature",
		DestBranch:   "main",
		Results: []*golang.GoModPathSyncResult{
			{
				SyncModRes: &golang.SyncModRes{
					Path: "go.mod",
				},
				TidyRan: false,
			},
		},
	}

	out := res.ToMarkdown()
	require.NotEmpty(t, out)
	require.Contains(t, out, "|")
	require.Contains(t, out, "go.mod")
	require.Contains(t, out, "no") // tidy not run
}

// TestSyncGoModRes_ToMarkdown_NilRes verifies that calling ToMarkdown on a nil
// pointer returns an empty string without panicking.
func TestSyncGoModRes_ToMarkdown_NilRes(t *testing.T) {
	t.Parallel()

	var res *SyncGoModRes
	require.Empty(t, res.ToMarkdown(), "nil receiver must return empty string")
}

// TestSyncGoModRes_ToJSON_Roundtrip verifies that ToJSON marshals a result and
// the output can be unmarshaled back into the expected fields, including the
// new new_branch and commit_sha fields.
func TestSyncGoModRes_ToJSON_Roundtrip(t *testing.T) {
	t.Parallel()

	res := &SyncGoModRes{
		SourceBranch: "feature",
		DestBranch:   "main",
		NewBranch:    "feature-into-main",
		CommitSHA:    "abc123def456",
		Results: []*golang.GoModPathSyncResult{
			{
				SyncModRes: &golang.SyncModRes{
					Path: "go.mod",
					Changes: []*golang.SyncChange{
						{
							Directive:  golang.DirectiveRequire,
							Module:     "github.com/foo/bar",
							OldVersion: "v1.0.0",
							NewVersion: "v1.2.0",
							Action:     golang.SyncActionUpdated,
						},
					},
				},
				TidyRan: true,
			},
		},
	}

	b, err := res.ToJSON()
	require.NoError(t, err)
	require.NotEmpty(t, b)

	var parsed map[string]any
	require.NoError(t, json.Unmarshal(b, &parsed))
	require.Equal(t, "feature", parsed["source_branch"])
	require.Equal(t, "main", parsed["dest_branch"])
	require.Equal(t, "feature-into-main", parsed["new_branch"])
	require.Equal(t, "abc123def456", parsed["commit_sha"])

	results, ok := parsed["results"].([]any)
	require.True(t, ok, "results must be a JSON array")
	require.Len(t, results, 1)

	pr := results[0].(map[string]any)
	require.Equal(t, "go.mod", pr["path"])
	require.Equal(t, true, pr["tidy_ran"])
}

// TestSyncGoModRes_ToJSON_CommitSHAOmittedWhenEmpty verifies that commit_sha is
// absent from JSON output when the field is empty.
func TestSyncGoModRes_ToJSON_CommitSHAOmittedWhenEmpty(t *testing.T) {
	t.Parallel()

	res := &SyncGoModRes{
		SourceBranch: "feature",
		DestBranch:   "main",
		NewBranch:    "feature-into-main",
	}

	b, err := res.ToJSON()
	require.NoError(t, err)

	var parsed map[string]any
	require.NoError(t, json.Unmarshal(b, &parsed))
	require.Equal(t, "feature-into-main", parsed["new_branch"])
	_, hasCommitSHA := parsed["commit_sha"]
	require.False(t, hasCommitSHA, "commit_sha must be absent when empty (omitempty)")
}

// TestSyncGoModRes_ToJSON_NilRes verifies that calling ToJSON on a nil pointer
// returns an error without panicking.
func TestSyncGoModRes_ToJSON_NilRes(t *testing.T) {
	t.Parallel()

	var res *SyncGoModRes
	b, err := res.ToJSON()
	require.Error(t, err, "nil receiver must return an error")
	require.Nil(t, b)
}

// TestPathSyncResult_JSONMarshal_TidyFields verifies that PathSyncResult
// marshals tidy_ran and tidy_error fields correctly and that tidy_error is
// omitted when empty.
func TestPathSyncResult_JSONMarshal_TidyFields(t *testing.T) {
	t.Parallel()

	t.Run("tidy_ran true with tidy_error", func(t *testing.T) {
		t.Parallel()

		pr := &golang.GoModPathSyncResult{
			SyncModRes: &golang.SyncModRes{Path: "go.mod"},
			TidyRan:    true,
			TidyError:  "go mod tidy failed: exit status 1",
		}

		b, err := json.Marshal(pr)
		require.NoError(t, err)

		var parsed map[string]any
		require.NoError(t, json.Unmarshal(b, &parsed))
		require.Equal(t, true, parsed["tidy_ran"])
		require.Equal(t, "go mod tidy failed: exit status 1", parsed["tidy_error"])
	})

	t.Run("tidy_ran false omits tidy_error", func(t *testing.T) {
		t.Parallel()

		pr := &golang.GoModPathSyncResult{
			SyncModRes: &golang.SyncModRes{Path: "go.mod"},
			TidyRan:    false,
		}

		b, err := json.Marshal(pr)
		require.NoError(t, err)

		var parsed map[string]any
		require.NoError(t, json.Unmarshal(b, &parsed))
		require.Equal(t, false, parsed["tidy_ran"])
		_, hasTidyError := parsed["tidy_error"]
		require.False(t, hasTidyError, "tidy_error must be absent when empty (omitempty)")
	})

	t.Run("TidyErr unexported from JSON", func(t *testing.T) {
		t.Parallel()

		pr := &golang.GoModPathSyncResult{
			SyncModRes: &golang.SyncModRes{Path: "go.mod"},
			TidyRan:    true,
			TidyError:  "some error",
			TidyErr:    errors.New("some error"),
		}

		b, err := json.Marshal(pr)
		require.NoError(t, err)

		var parsed map[string]any
		require.NoError(t, json.Unmarshal(b, &parsed))
		// TidyErr has json:"-" so it must never appear in JSON output.
		_, hasTidyErr := parsed["TidyErr"]
		require.False(t, hasTidyErr, "TidyErr (json:\"-\") must not appear in JSON output")
	})
}

// TestSyncGoModReq_BranchName_LongBaseIsTruncated verifies that branchName
// truncates the result to 240 characters when the generated base exceeds the cap.
func TestSyncGoModReq_BranchName_LongBaseIsTruncated(t *testing.T) {
	t.Parallel()

	req := &SyncGoModReq{
		SourceBranch: strings.Repeat("a", 150),
		DestBranch:   strings.Repeat("b", 150),
	}
	// base = 150 + "-into-" + 150 = 306 chars; must be capped at 240
	got := req.branchName()
	require.LessOrEqual(t, len(got), 240)
}

// TestSyncGoModReq_BranchName_LongTicketIsTruncated verifies that branchName
// truncates the result to 240 characters when the ticket prefix pushes the name
// past the cap.
func TestSyncGoModReq_BranchName_LongTicketIsTruncated(t *testing.T) {
	t.Parallel()

	req := &SyncGoModReq{
		SourceBranch: strings.Repeat("a", 100),
		DestBranch:   strings.Repeat("b", 100),
		Ticket:       strings.Repeat("T", 50),
	}
	// full = 50 + 1 + 100 + "-into-" + 100 = 257 chars; must be capped at 240
	got := req.branchName()
	require.LessOrEqual(t, len(got), 240)
}

// TestSyncGoModReq_Validate_Paths verifies that we sync the go.mod in the
// working directory when we're not given a path, that the paths we're given are
// used as they are, and that in workspace mode we have no paths because we get
// them from the go.work, so they can't be given as well.
func TestSyncGoModReq_Validate_Paths(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		paths            []string
		workspaceModules bool
		want             []string
		errHas           string
	}{
		"nil paths":         {want: []string{"go.mod"}},
		"empty paths":       {paths: []string{}, want: []string{"go.mod"}},
		"given paths":       {paths: []string{"go.mod", "sdk/go.mod"}, want: []string{"go.mod", "sdk/go.mod"}},
		"workspace modules": {workspaceModules: true},
		"workspace modules with paths": {
			paths:            []string{"go.mod"},
			workspaceModules: true,
			errHas:           "cannot be combined with workspace-modules",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := &SyncGoModReq{
				SourceBranch:     "main",
				DestBranch:       "ce/main",
				Paths:            test.paths,
				WorkspaceModules: test.workspaceModules,
			}

			err := req.validate()
			if test.errHas != "" {
				require.ErrorContains(t, err, test.errHas)

				return
			}
			require.NoError(t, err)
			require.Equal(t, test.want, req.Paths)
		})
	}
}

// TestSyncGoModRes_GoWorkSum verifies that syncing the workspace is shown as a
// go.work.sum row in the table, only when we tried to sync it, with a change
// when it was updated and the reason when it failed, and that it's in the JSON.
func TestSyncGoModRes_GoWorkSum(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		goWorkSum *golang.GoWorkSumSyncResult
		wantRow   string // what the row has after the path, empty when there is no row
		wantJSON  string
	}{
		"not synced": {},
		"unchanged": {
			goWorkSum: &golang.GoWorkSumSyncResult{},
			wantRow:   "go.work.sum 0 no",
			wantJSON:  `"go_work_sum":{"updated":false}`,
		},
		"updated": {
			goWorkSum: &golang.GoWorkSumSyncResult{Updated: true},
			wantRow:   "go.work.sum 1 no",
			wantJSON:  `"go_work_sum":{"updated":true}`,
		},
		"failed": {
			goWorkSum: &golang.GoWorkSumSyncResult{Error: "download failed"},
			wantRow:   "go.work.sum 0 no download failed",
			wantJSON:  `"go_work_sum":{"updated":false,"error":"download failed"}`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			res := &SyncGoModRes{
				SourceBranch: "main",
				DestBranch:   "ce/main",
				Results: []*golang.GoModPathSyncResult{
					{SyncModRes: &golang.SyncModRes{Path: "go.mod"}},
				},
				GoWorkSum: test.goWorkSum,
			}

			out := strings.Join(strings.Fields(res.ToString()), " ")
			b, err := res.ToJSON()
			require.NoError(t, err)

			if test.goWorkSum == nil {
				require.NotContains(t, out, "go.work.sum")
				require.NotContains(t, string(b), "go_work_sum")

				return
			}
			require.Contains(t, out, test.wantRow)
			require.Contains(t, string(b), test.wantJSON)
		})
	}
}
