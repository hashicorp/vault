// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"errors"
	"testing"

	libgithub "github.com/google/go-github/v68/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/changed"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/stretchr/testify/require"
)

// TestCreateBackportReq_Validate tests validation of the request
func TestCreateBackportReq_Validate(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		req   *CreateBackportReq
		valid bool
	}{
		"empty": {nil, false},
		"valid": {NewCreateBackportReq(WithCreateBrackportReqPullNumber(1234)), true},
		"no owner": {
			NewCreateBackportReq(
				WithCreateBrackportReqPullNumber(1234),
				WithCreateBackportReqOwner(""),
			), false,
		},
		"no repo": {
			NewCreateBackportReq(
				WithCreateBrackportReqPullNumber(1234),
				WithCreateBrackportReqRepo(""),
			), false,
		},
		"no pull number": {NewCreateBackportReq(), false},
		"no ce branch prefix": {
			NewCreateBackportReq(
				WithCreateBrackportReqPullNumber(1234),
				WithCreateBrackportReqCEBranchPrefix(""),
			), false,
		},
		"no base origin": {
			NewCreateBackportReq(
				WithCreateBrackportReqPullNumber(1234),
				WithCreateBrackportReqBaseOrigin(""),
			), false,
		},
		"uninitialized exclude groups": {
			NewCreateBackportReq(
				WithCreateBrackportReqPullNumber(1234),
				WithCreateBrackportReqCEExclude(nil),
			), false,
		},
		"uninitialized inactive groups": {
			NewCreateBackportReq(
				WithCreateBrackportReqPullNumber(1234),
				WithCreateBrackportReqAllowInactiveGroups(nil),
			), false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if test.valid {
				require.NoError(t, test.req.Validate(context.Background()))
			} else {
				require.Error(t, test.req.Validate(context.Background()))
			}
		})
	}
}

// TestCreateBackportReq_backportNameForRef tests generating the backport
// branch name from branch name ref and the original PR branch name.
func TestCreateBackportReq_backportNameForRef(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		ref      string // These should be full branch names
		prBranch string
		expected string
	}{
		// backporting to ent main should never really happen but we'll test the
		// logic anyway
		"ent main": {
			"main",
			"my-pr",
			"backport/main/my-pr",
		},
		"ent release branch": {
			"release/1.19.x+ent",
			"my-pr",
			"backport/release/1.19.x+ent/my-pr",
		},
		"ce main": {
			"ce/main",
			"my-pr",
			"backport/ce/main/my-pr",
		},
		"ce release branch": {
			"ce/release/1.19.x",
			"my-pr",
			"backport/ce/release/1.19.x/my-pr",
		},
		"truncates super long branch name": {
			"main",
			"my-really-really-long-pr-name-that-must-exceed-two-hundred-and-fifty-characters-when-it-is-appended-to-the-backport-and-base-ref-prefixes-ought-to-be-truncated-so-as-to-not-exceed-the-github-pr-branch-requirements-otherwise-bad-things-happen",
			"backport/main/my-really-really-long-pr-name-that-must-exceed-two-hundred-and-fifty-characters-when-it-is-appended-to-the-backport-and-base-ref-prefixes-ought-to-be-truncated-so-as-to-not-exceed-the-github-pr-branch-requirements-otherwise-bad-things-h",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			req := NewCreateBackportReq()
			require.Equal(t, test.expected, req.backportBranchNameForRef(test.ref, test.prBranch))
		})
	}
}

// TestCreateBackportReq_baseRefVersion tests generating the base ref version
// from the backport branch reference. The base ref version matches the schema
// used in .release/versions.hcl.
func TestCreateBackportReq_baseRefVersion(t *testing.T) {
	t.Parallel()

	for ref, test := range map[string]struct {
		req         *CreateBackportReq
		expectedRef string
	}{
		// backporting to ent main should never really happen but we'll test the
		// logic anyway
		"main":               {req: NewCreateBackportReq(), expectedRef: "main"},
		"ce/main":            {req: NewCreateBackportReq(), expectedRef: "main"},
		"ent/main":           {req: NewCreateBackportReq(WithCreateBrackportReqEntBranchPrefix("ent")), expectedRef: "main"},
		"release/1.19.x+ent": {req: NewCreateBackportReq(), expectedRef: "release/1.19.x"},
		"ce/release/1.19.x":  {req: NewCreateBackportReq(), expectedRef: "release/1.19.x"},
		"ent/release/1.19.x": {req: NewCreateBackportReq(WithCreateBrackportReqEntBranchPrefix("ent")), expectedRef: "release/1.19.x"},
	} {
		t.Run(ref, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, test.expectedRef, test.req.baseRefVersion(ref))
		})
	}
}

// TestCreateBackportReq_determineBackportRefs tests generating a list
// of backport refs when considering the base ref of the PR and any labels
// that have been applied to it.
func TestCreateBackportReq_determineBackportRefs(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		req      *CreateBackportReq
		baseRef  string
		labels   Labels
		expected []string
	}{
		"ent main no labels": {
			NewCreateBackportReq(),
			"main",
			nil,
			[]string{"ce/main"},
		},
		"ent main no labels with ent prefix": {
			NewCreateBackportReq(WithCreateBrackportReqEntBranchPrefix("ent")),
			"ent/main",
			nil,
			[]string{"ce/main"},
		},
		"ent main with labels": {
			NewCreateBackportReq(),
			"main",
			Labels{
				&libgithub.Label{Name: libgithub.Ptr("backport/1.19.x")},
				&libgithub.Label{Name: libgithub.Ptr("backport/1.18.x")},
			},
			[]string{"ce/main", "release/1.19.x+ent", "release/1.18.x+ent"},
		},
		"ent main with labels with ent prefix": {
			NewCreateBackportReq(WithCreateBrackportReqEntBranchPrefix("ent")),
			"ent/main",
			Labels{
				&libgithub.Label{Name: libgithub.Ptr("backport/1.19.x")},
				&libgithub.Label{Name: libgithub.Ptr("backport/1.18.x")},
			},
			[]string{"ce/main", "ent/release/1.19.x+ent", "ent/release/1.18.x+ent"},
		},
		"ent release no labels": {
			NewCreateBackportReq(),
			"release/1.19.x+ent",
			nil,
			[]string{"ce/release/1.19.x"},
		},
		"ent release no labels with ent prefix": {
			NewCreateBackportReq(WithCreateBrackportReqEntBranchPrefix("ent")),
			"ent/release/1.19.x+ent",
			nil,
			[]string{"ce/release/1.19.x"},
		},
		"ent release with labels": {
			NewCreateBackportReq(),
			"release/1.19.x+ent",
			Labels{
				&libgithub.Label{Name: libgithub.Ptr("backport/1.18.x")},
				&libgithub.Label{Name: libgithub.Ptr("backport/1.17.x")},
				&libgithub.Label{Name: libgithub.Ptr("backport/1.16.x")},
			},
			[]string{
				"ce/release/1.19.x",
				"release/1.18.x+ent",
				"release/1.17.x+ent",
				"release/1.16.x+ent",
			},
		},
		"ent release with labels with ent prefix": {
			NewCreateBackportReq(WithCreateBrackportReqEntBranchPrefix("ent")),
			"ent/release/1.19.x+ent",
			Labels{
				&libgithub.Label{Name: libgithub.Ptr("backport/1.18.x")},
				&libgithub.Label{Name: libgithub.Ptr("backport/1.17.x")},
				&libgithub.Label{Name: libgithub.Ptr("backport/1.16.x")},
			},
			[]string{
				"ce/release/1.19.x",
				"ent/release/1.18.x+ent",
				"ent/release/1.17.x+ent",
				"ent/release/1.16.x+ent",
			},
		},
		"ce main no labels": {
			NewCreateBackportReq(),
			"ce/main",
			nil,
			nil,
		},
		"ce main with labels": {
			NewCreateBackportReq(),
			"ce/main",
			Labels{
				&libgithub.Label{Name: libgithub.Ptr("backport/1.19.x")},
				&libgithub.Label{Name: libgithub.Ptr("backport/1.18.x")},
			},
			[]string{"ce/release/1.19.x", "ce/release/1.18.x"},
		},
		"ce release no labels": {
			NewCreateBackportReq(),
			"ce/release/1.19.x",
			nil,
			nil,
		},
		"ce release with labels": {
			NewCreateBackportReq(),
			"ce/release/1.19.x",
			Labels{
				&libgithub.Label{Name: libgithub.Ptr("backport/1.18.x")},
				&libgithub.Label{Name: libgithub.Ptr("backport/1.17.x")},
				&libgithub.Label{Name: libgithub.Ptr("backport/1.16.x")},
			},
			[]string{"ce/release/1.18.x", "ce/release/1.17.x", "ce/release/1.16.x"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.EqualValues(t, test.expected, test.req.determineBackportRefs(context.Background(), test.baseRef, test.labels))
		})
	}
}

// TestCreateBackportReq_shouldSkipRef tests whether various combinations of
// base refs, backport refs, changed files, and active CE versions are
// backportable references or should be skipped.
func TestCreateBackportReq_shouldSkipRef(t *testing.T) {
	t.Parallel()

	defaultActiveVersions := map[string]*releases.Version{
		// main is never going to be in here as it's assumed it's always active
		"1.19.x": {CEActive: true, LTS: true},
		"1.18.x": {CEActive: false},
		"1.17.x": {CEActive: false},
		"1.16.x": {CEActive: true, LTS: true},
	}

	defaultChangedFiles := &ListChangedFilesRes{
		Files: changed.Files{
			{
				File: &libgithub.CommitFile{
					SHA:      libgithub.Ptr("84e0b544965861a7c6373e639cb13755512f84f4"),
					Filename: libgithub.Ptr(".github/workflows/build.yml"),
				},
				Groups: changed.FileGroups{"pipeline"},
			},
			{
				File: &libgithub.CommitFile{
					SHA:      libgithub.Ptr("84e0b544965861a7c6373e639cb13755512f84f4"),
					Filename: libgithub.Ptr("go.mod"),
				},
				Groups: changed.FileGroups{"app", "gotoolchain"},
			},
		},
		Groups: changed.FileGroups{
			"app", "gotoolchain", "pipeline",
		},
	}

	for name, test := range map[string]struct {
		baseRefVersion string
		ref            string
		activeVersions map[string]*releases.Version
		changedFiles   *ListChangedFilesRes
		skip           bool
	}{
		"main to ce/main": {
			baseRefVersion: "main",
			ref:            "ce/main",
			activeVersions: defaultActiveVersions,
			changedFiles:   defaultChangedFiles,
			skip:           false,
		},
		"main to ce/main with ent and ce only files": {
			baseRefVersion: "main",
			ref:            "ce/main",
			activeVersions: defaultActiveVersions,
			changedFiles: &ListChangedFilesRes{
				Files: changed.Files{
					{
						File: &libgithub.CommitFile{
							SHA:      libgithub.Ptr("e1c10eae02e13f5a090b9c29b0b1a3003e8ca7f6"),
							Filename: libgithub.Ptr("go.mod"),
						},
						Groups: changed.FileGroups{"app", "gotoolchain"},
					},
					{
						File: &libgithub.CommitFile{
							SHA:      libgithub.Ptr("a6397662ea1d5fdde744ff3e4246377cf369197a"),
							Filename: libgithub.Ptr("vault_ent/go.mod"),
						},
						Groups: changed.FileGroups{"app", "enterprise", "gotoolchain"},
					},
				},
				Groups: changed.FileGroups{
					"app", "enterprise", "gotoolchain",
				},
			},
			skip: false,
		},
		"main to active release/1.19.x+ent": {
			baseRefVersion: "1.19.x",
			ref:            "release/1.19.x+ent",
			activeVersions: defaultActiveVersions,
			changedFiles:   defaultChangedFiles,
			skip:           false,
		},
		"main to release/1.18.x+ent (inactive CE)": {
			baseRefVersion: "1.18.x",
			ref:            "release/1.18.x+ent",
			activeVersions: defaultActiveVersions,
			changedFiles:   defaultChangedFiles,
			skip:           false,
		},
		"active release branch with app changes": {
			baseRefVersion: "1.19.x",
			ref:            "ce/release/1.19.x",
			activeVersions: defaultActiveVersions,
			changedFiles:   defaultChangedFiles,
			skip:           false,
		},
		"active release branch to CE with only ent changes": {
			baseRefVersion: "1.19.x",
			ref:            "ce/release/1.19.x",
			activeVersions: defaultActiveVersions,
			changedFiles: &ListChangedFilesRes{
				Files: changed.Files{
					{
						File: &libgithub.CommitFile{
							SHA:      libgithub.Ptr("84e0b544965861a7c6373e639cb13755512f84f4"),
							Filename: libgithub.Ptr(".github/workflows/build-artifacts-ent.yml"),
						},
						Groups: changed.FileGroups{"enterprise", "pipeline"},
					},
					{
						File: &libgithub.CommitFile{
							SHA:      libgithub.Ptr("84e0b544965861a7c6373e639cb13755512f84f4"),
							Filename: libgithub.Ptr("vault/vault_ent/go.mod"),
						},
						Groups: changed.FileGroups{"app", "enterprise", "gotoolchain"},
					},
				},
				Groups: changed.FileGroups{
					"app", "enterprise", "gotoolchain", "pipeline",
				},
			},
			skip: true,
		},
		"inactive ce branch with no allowed group changes": {
			baseRefVersion: "1.18.x",
			ref:            "ce/release/1.18.x",
			activeVersions: defaultActiveVersions,
			changedFiles: &ListChangedFilesRes{
				Files: changed.Files{
					{
						File: &libgithub.CommitFile{
							SHA:      libgithub.Ptr("84e0b544965861a7c6373e639cb13755512f84f4"),
							Filename: libgithub.Ptr("vault/go.mod"),
						},
						Groups: changed.FileGroups{"app", "gotoolchain"},
					},
					{
						File: &libgithub.CommitFile{
							SHA:      libgithub.Ptr("84e0b544965861a7c6373e639cb13755512f84f4"),
							Filename: libgithub.Ptr("vault/vault_ent/go.mod"),
						},
						Groups: changed.FileGroups{"app", "enterprise", "gotoolchain"},
					},
				},
				Groups: changed.FileGroups{
					"app", "enterprise", "gotoolchain",
				},
			},

			skip: true,
		},
		"inactive ce with with pipeline changes": {
			baseRefVersion: "1.18.x",
			ref:            "ce/release/1.18.x",
			activeVersions: defaultActiveVersions,
			changedFiles: &ListChangedFilesRes{
				Files: changed.Files{
					{
						File: &libgithub.CommitFile{
							SHA:      libgithub.Ptr("84e0b544965861a7c6373e639cb13755512f84f4"),
							Filename: libgithub.Ptr(".github/workflows/build.yml"),
						},
						Groups: changed.FileGroups{"pipeline"},
					},
					{
						File: &libgithub.CommitFile{
							SHA:      libgithub.Ptr("84e0b544965861a7c6373e639cb13755512f84f4"),
							Filename: libgithub.Ptr("vault/go.mod"),
						},
						Groups: changed.FileGroups{"app", "gotoolchain"},
					},
				},
				Groups: changed.FileGroups{
					"app", "gotoolchain", "pipeline",
				},
			},

			skip: false,
		},
		"inactive ce with with docs changes": {
			baseRefVersion: "1.17.x",
			ref:            "ce/release/1.17.x",
			activeVersions: defaultActiveVersions,
			changedFiles: &ListChangedFilesRes{
				Files: changed.Files{
					{
						File: &libgithub.CommitFile{
							SHA:      libgithub.Ptr("84e0b544965861a7c6373e639cb13755512f84f4"),
							Filename: libgithub.Ptr("website/content/docs/index.mdx"),
						},
						Groups: changed.FileGroups{"docs"},
					},
					{
						File: &libgithub.CommitFile{
							SHA:      libgithub.Ptr("84e0b544965861a7c6373e639cb13755512f84f4"),
							Filename: libgithub.Ptr("vault/go.mod"),
						},
						Groups: changed.FileGroups{"app", "gotoolchain"},
					},
				},
				Groups: changed.FileGroups{
					"app", "gotoolchain", "pipeline",
				},
			},

			skip: false,
		},
		"inactive ce with with changelog changes": {
			baseRefVersion: "1.17.x",
			ref:            "ce/release/1.17.x",
			activeVersions: defaultActiveVersions,
			changedFiles: &ListChangedFilesRes{
				Files: changed.Files{
					{
						File: &libgithub.CommitFile{
							SHA:      libgithub.Ptr("84e0b544965861a7c6373e639cb13755512f84f4"),
							Filename: libgithub.Ptr("changelog/1234.txt"),
						},
						Groups: changed.FileGroups{"changelog"},
					},
					{
						File: &libgithub.CommitFile{
							SHA:      libgithub.Ptr("84e0b544965861a7c6373e639cb13755512f84f4"),
							Filename: libgithub.Ptr("vault/go.mod"),
						},
						Groups: changed.FileGroups{"app", "gotoolchain"},
					},
				},
				Groups: changed.FileGroups{
					"app", "gotoolchain", "pipeline",
				},
			},

			skip: false,
		},
		"empty changed files list is skipped": {
			baseRefVersion: "1.19.x",
			ref:            "ce/release/1.19.x",
			activeVersions: defaultActiveVersions,
			changedFiles: &ListChangedFilesRes{
				Files:  changed.Files{},
				Groups: changed.FileGroups{},
			},
			skip: true,
		},
		"nil changed files list is skipped": {
			baseRefVersion: "1.19.x",
			ref:            "ce/release/1.19.x",
			activeVersions: defaultActiveVersions,
			changedFiles:   nil,
			skip:           true,
		},
		"release branch with no active versions": {
			baseRefVersion: "1.19.x",
			ref:            "ce/release/1.19.x",
			activeVersions: map[string]*releases.Version{},
			changedFiles: &ListChangedFilesRes{
				Files: changed.Files{
					{
						File: &libgithub.CommitFile{
							SHA:      libgithub.Ptr("84e0b544965861a7c6373e639cb13755512f84f4"),
							Filename: libgithub.Ptr("vault/go.mod"),
						},
						Groups: changed.FileGroups{"app", "gotoolchain"},
					},
					{
						File: &libgithub.CommitFile{
							SHA:      libgithub.Ptr("84e0b544965861a7c6373e639cb13755512f84f4"),
							Filename: libgithub.Ptr("vault/vault_ent/go.mod"),
						},
						Groups: changed.FileGroups{"app", "enterprise", "gotoolchain"},
					},
				},
				Groups: changed.FileGroups{
					"app", "enterprise", "gotoolchain",
				},
			},
			skip: true,
		},
		"release branch with nil active versions": {
			baseRefVersion: "1.19.x",
			ref:            "ce/release/1.19.x",
			activeVersions: nil,
			changedFiles: &ListChangedFilesRes{
				Files: changed.Files{
					{
						File: &libgithub.CommitFile{
							SHA:      libgithub.Ptr("84e0b544965861a7c6373e639cb13755512f84f4"),
							Filename: libgithub.Ptr("vault/go.mod"),
						},
						Groups: changed.FileGroups{"app", "gotoolchain"},
					},
					{
						File: &libgithub.CommitFile{
							SHA:      libgithub.Ptr("84e0b544965861a7c6373e639cb13755512f84f4"),
							Filename: libgithub.Ptr("vault/vault_ent/go.mod"),
						},
						Groups: changed.FileGroups{"app", "enterprise", "gotoolchain"},
					},
				},
			},
			skip: true,
		},
		"missing base ref version": {
			baseRefVersion: "",
			ref:            "ce/main",
			activeVersions: defaultActiveVersions,
			changedFiles:   defaultChangedFiles,
			skip:           true,
		},
		"missing ref version": {
			baseRefVersion: "main",
			ref:            "",
			activeVersions: defaultActiveVersions,
			changedFiles:   defaultChangedFiles,
			skip:           true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := NewCreateBackportReq()
			msg, skip := req.shouldSkipRef(
				context.Background(),
				test.baseRefVersion,
				test.ref,
				test.activeVersions,
				test.changedFiles,
			)
			require.Equalf(
				t, test.skip, skip, "should have %t but got %t with %s", test.skip, skip, msg)
		})
	}
}

func TestCreateBackportReq_pullRequestBody(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		expectContains    []string
		expectNotContains []string
		origin            *libgithub.PullRequest
		attempt           *CreateBackportAttempt
	}{
		"no error": {
			expectContains:    []string{"original body"},
			expectNotContains: []string{"error body"},
			origin: &libgithub.PullRequest{
				Body:     libgithub.Ptr("original body"),
				Number:   libgithub.Ptr(1234),
				HTMLURL:  libgithub.Ptr("https://github.com/hashicorp/vault-enterprise/pull/1234"),
				MergedBy: &libgithub.User{Login: libgithub.Ptr("my-login")},
			},
			attempt: &CreateBackportAttempt{
				TargetRef: "release/1.19.x",
			},
		},
		"error": {
			expectContains: []string{"original body", "error body"},
			origin: &libgithub.PullRequest{
				Body:     libgithub.Ptr("original body"),
				Number:   libgithub.Ptr(1234),
				HTMLURL:  libgithub.Ptr("https://github.com/hashicorp/vault-enterprise/pull/1234"),
				MergedBy: &libgithub.User{Login: libgithub.Ptr("my-login")},
			},
			attempt: &CreateBackportAttempt{
				TargetRef: "release/1.19.x",
				Error:     errors.New("error body"),
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			req := NewCreateBackportReq()
			got, err := req.pullRequestBody(test.origin, test.attempt)
			require.NoError(t, err)
			for _, c := range test.expectContains {
				require.Containsf(t, got, c, got)
			}
			for _, nc := range test.expectNotContains {
				require.NotContainsf(t, got, nc, got)
			}
		})
	}
}

func TestCreateBackportRes_Err(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		in     *CreateBackportRes
		failed error
	}{
		"nil": {
			nil,
			errors.New("uninitialized"),
		},
		"no errors": {
			&CreateBackportRes{
				Attempts: map[string]*CreateBackportAttempt{
					"ce/main":        {},
					"release/1.18.x": {},
					"release/1.19.x": {},
				},
			},
			nil,
		},
		"top level error no attempt errors": {
			&CreateBackportRes{
				Error: errors.New("top-failed"),
				Attempts: map[string]*CreateBackportAttempt{
					"ce/main":        {},
					"release/1.18.x": {},
					"release/1.19.x": {},
				},
			},
			errors.New("top-failed"),
		},
		"no top level error attempt errors": {
			&CreateBackportRes{
				Attempts: map[string]*CreateBackportAttempt{
					"ce/main": {
						Error: errors.New("child-failed"),
					},
					"release/1.18.x": {},
					"release/1.19.x": {},
				},
			},
			errors.New("child-failed"),
		},
		"top level and attempt errors": {
			&CreateBackportRes{
				Error: errors.New("top-failed"),
				Attempts: map[string]*CreateBackportAttempt{
					"ce/main":        {},
					"release/1.18.x": {},
					"release/1.19.x": {
						Error: errors.New("child-failed"),
					},
				},
			},
			errors.New("top-failed\nchild-failed"),
		},
		"multiple attempt errors": {
			&CreateBackportRes{
				Error: errors.New("top-failed"),
				Attempts: map[string]*CreateBackportAttempt{
					"ce/main": {},
					"release/1.18.x": {
						Error: errors.New("child-2-failed"),
					},
					"release/1.19.x": {
						Error: errors.New("child-3-failed"),
					},
				},
			},
			// When multiple attempts fail the errros should be stable
			errors.New("top-failed\nchild-2-failed\nchild-3-failed"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			if test.failed == nil {
				require.Nil(t, test.in.Err())
			} else {
				require.Equal(t, test.failed.Error(), test.in.Err().Error())
			}
		})
	}
}
