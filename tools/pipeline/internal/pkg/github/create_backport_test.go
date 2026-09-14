// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"errors"
	"net/http"
	"testing"

	libgithub "github.com/google/go-github/v83/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/changed"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/config"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/stretchr/testify/require"
)

// TestCreateBackportReq_Validate tests validation of the request
func TestCreateBackportReq_Validate(t *testing.T) {
	t.Parallel()

	changedCfg := func() *changed.Config {
		return &changed.Config{
			Groups: []*changed.GroupConfig{
				{
					Name: "go",
					Match: changed.Matchers{
						{Extension: []string{".go"}},
					},
				},
				{
					Name: "python",
					Match: changed.Matchers{
						{Extension: []string{".py"}},
					},
				},
			},
		}
	}

	versionsCfg := func() *releases.VersionsConfig {
		return &releases.VersionsConfig{
			Schema: 1,
			ActiveVersion: &releases.ActiveVersion{
				Versions: map[string]*releases.Version{
					"1.19.x": {CEActive: true, LTS: true},
					"1.18.x": {CEActive: true, LTS: false},
					"1.17.x": {CEActive: false, LTS: false},
					"1.16.x": {CEActive: false, LTS: true},
				},
			},
		}
	}
	configDecodeRes := func() *config.DecodeRes {
		return &config.DecodeRes{
			Config: &config.Config{
				ChangedFiles: changedCfg(),
			},
		}
	}
	versionsDecodeRes := func() *releases.DecodeRes {
		return &releases.DecodeRes{
			Config: versionsCfg(),
		}
	}

	for name, test := range map[string]struct {
		req   *CreateBackportReq
		valid bool
	}{
		"empty": {nil, false},
		"valid": {
			NewCreateBackportReq(
				WithCreateBackportReqConfigDecodeRes(configDecodeRes()),
				WithCreateBackportReqVersionsDecodeRes(versionsDecodeRes()),
				WithCreateBackportReqPullNumber(1234),
			), true,
		},
		"valid - with backport failed label": {
			NewCreateBackportReq(
				WithCreateBackportReqConfigDecodeRes(configDecodeRes()),
				WithCreateBackportReqVersionsDecodeRes(versionsDecodeRes()),
				WithCreateBackportReqPullNumber(1234),
				WithCreateBackportReqBackportFailedLabel("backport-failed"),
			), true,
		},
		"no changed file config": {
			NewCreateBackportReq(
				WithCreateBackportReqVersionsDecodeRes(versionsDecodeRes()),
				WithCreateBackportReqPullNumber(1234),
			), false,
		},
		"no versions config": {
			NewCreateBackportReq(
				WithCreateBackportReqConfigDecodeRes(configDecodeRes()),
				WithCreateBackportReqPullNumber(1234),
			), false,
		},
		"no owner": {
			NewCreateBackportReq(
				WithCreateBackportReqConfigDecodeRes(configDecodeRes()),
				WithCreateBackportReqVersionsDecodeRes(versionsDecodeRes()),
				WithCreateBackportReqPullNumber(1234),
				WithCreateBackportReqOwner(""),
			), false,
		},
		"no repo": {
			NewCreateBackportReq(
				WithCreateBackportReqConfigDecodeRes(configDecodeRes()),
				WithCreateBackportReqVersionsDecodeRes(versionsDecodeRes()),
				WithCreateBackportReqPullNumber(1234),
				WithCreateBackportReqRepo(""),
			), false,
		},
		"no pull number": {
			NewCreateBackportReq(
				WithCreateBackportReqConfigDecodeRes(configDecodeRes()),
				WithCreateBackportReqVersionsDecodeRes(versionsDecodeRes()),
			), false,
		},
		"no ce branch prefix": {
			NewCreateBackportReq(
				WithCreateBackportReqConfigDecodeRes(configDecodeRes()),
				WithCreateBackportReqVersionsDecodeRes(versionsDecodeRes()),
				WithCreateBackportReqPullNumber(1234),
				WithCreateBackportReqCEBranchPrefix(""),
			), false,
		},
		"no base origin": {
			NewCreateBackportReq(
				WithCreateBackportReqConfigDecodeRes(configDecodeRes()),
				WithCreateBackportReqVersionsDecodeRes(versionsDecodeRes()),
				WithCreateBackportReqPullNumber(1234),
				WithCreateBackportReqBaseOrigin(""),
			), false,
		},
		"uninitialized exclude groups": {
			NewCreateBackportReq(
				WithCreateBackportReqConfigDecodeRes(configDecodeRes()),
				WithCreateBackportReqVersionsDecodeRes(versionsDecodeRes()),
				WithCreateBackportReqPullNumber(1234),
				WithCreateBackportReqCEExclude(nil),
			), false,
		},
		"uninitialized inactive groups": {
			NewCreateBackportReq(
				WithCreateBackportReqConfigDecodeRes(configDecodeRes()),
				WithCreateBackportReqVersionsDecodeRes(versionsDecodeRes()),
				WithCreateBackportReqPullNumber(1234),
				WithCreateBackportReqAllowInactiveGroups(nil),
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
		"ent/main":           {req: NewCreateBackportReq(WithCreateBackportReqEntBranchPrefix("ent")), expectedRef: "main"},
		"release/1.19.x+ent": {req: NewCreateBackportReq(), expectedRef: "release/1.19.x"},
		"ce/release/1.19.x":  {req: NewCreateBackportReq(), expectedRef: "release/1.19.x"},
		"ent/release/1.19.x": {req: NewCreateBackportReq(WithCreateBackportReqEntBranchPrefix("ent")), expectedRef: "release/1.19.x"},
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
			NewCreateBackportReq(WithCreateBackportReqEntBranchPrefix("ent")),
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
			NewCreateBackportReq(WithCreateBackportReqEntBranchPrefix("ent")),
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
			NewCreateBackportReq(WithCreateBackportReqEntBranchPrefix("ent")),
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
			NewCreateBackportReq(WithCreateBackportReqEntBranchPrefix("ent")),
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
		"1.20.x": {CEActive: true},
		"1.19.x": {CEActive: false, LTS: true},
		"1.18.x": {CEActive: false},
		"1.16.x": {CEActive: true, LTS: true},
	}

	noChangedFiles := &ListChangedFilesRes{
		Files:  changed.Files{},
		Groups: changed.FileGroups{},
	}

	allowedInactiveCEChangedFiles := &ListChangedFilesRes{
		Files: changed.Files{
			{
				GithubCommitFile: &libgithub.CommitFile{
					SHA:      libgithub.Ptr("84e0b544965861a7c6373e639cb13755512f84f4"),
					Filename: libgithub.Ptr("changelog/_2837.md"),
				},
				Groups: changed.FileGroups{"changelog"},
			},
		},
		Groups: changed.FileGroups{
			"changelog",
		},
	}

	onlyEnterpriseChangedFiles := &ListChangedFilesRes{
		Files: changed.Files{
			{
				GithubCommitFile: &libgithub.CommitFile{
					SHA:      libgithub.Ptr("84e0b544965861a7c6373e639cb13755512f84f4"),
					Filename: libgithub.Ptr(".github/workflows/build-artifacts-ent.yml"),
				},
				Groups: changed.FileGroups{"enterprise", "pipeline"},
			},
			{
				GithubCommitFile: &libgithub.CommitFile{
					SHA:      libgithub.Ptr("84e0b544965861a7c6373e639cb13755512f84f4"),
					Filename: libgithub.Ptr("vault/vault_ent/go.mod"),
				},
				Groups: changed.FileGroups{"app", "enterprise", "gotoolchain"},
			},
		},
		Groups: changed.FileGroups{
			"app", "enterprise", "gotoolchain", "pipeline",
		},
	}

	mixedCEAndEnterpriseChangedFiles := &ListChangedFilesRes{
		Files: changed.Files{
			{
				GithubCommitFile: &libgithub.CommitFile{
					SHA:      libgithub.Ptr("e1c10eae02e13f5a090b9c29b0b1a3003e8ca7f6"),
					Filename: libgithub.Ptr("go.mod"),
				},
				Groups: changed.FileGroups{"app", "gotoolchain"},
			},
			{
				GithubCommitFile: &libgithub.CommitFile{
					SHA:      libgithub.Ptr("a6397662ea1d5fdde744ff3e4246377cf369197a"),
					Filename: libgithub.Ptr("vault_ent/go.mod"),
				},
				Groups: changed.FileGroups{"app", "enterprise", "gotoolchain"},
			},
		},
		Groups: changed.FileGroups{
			"app", "enterprise", "gotoolchain",
		},
	}

	allCEChangedFiles := &ListChangedFilesRes{
		Files: changed.Files{
			{
				GithubCommitFile: &libgithub.CommitFile{
					SHA:      libgithub.Ptr("84e0b544965861a7c6373e639cb13755512f84f4"),
					Filename: libgithub.Ptr(".github/workflows/build.yml"),
				},
				Groups: changed.FileGroups{"pipeline"},
			},
			{
				GithubCommitFile: &libgithub.CommitFile{
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
		// main -> ce/main
		"main to ce/main with no changed files": {
			baseRefVersion: "main",
			ref:            "ce/main",
			activeVersions: defaultActiveVersions,
			changedFiles:   noChangedFiles,
			skip:           true,
		},
		"main to ce/main with mixed changed files": {
			baseRefVersion: "main",
			ref:            "ce/main",
			activeVersions: defaultActiveVersions,
			changedFiles:   mixedCEAndEnterpriseChangedFiles,
			skip:           false,
		},
		"main to ce/main with enterprise only changed files": {
			baseRefVersion: "main",
			ref:            "ce/main",
			activeVersions: defaultActiveVersions,
			changedFiles:   onlyEnterpriseChangedFiles,
			skip:           true,
		},
		"main to ce/main with all CE changed files": {
			baseRefVersion: "main",
			ref:            "ce/main",
			activeVersions: defaultActiveVersions,
			changedFiles:   allCEChangedFiles,
			skip:           false,
		},
		"main to ce/main with allowed inactive changed files": {
			baseRefVersion: "main",
			ref:            "ce/main",
			activeVersions: defaultActiveVersions,
			changedFiles:   allowedInactiveCEChangedFiles,
			skip:           false,
		},
		// main -> release branch
		"main to release with no changed files": {
			baseRefVersion: "main",
			ref:            "release/1.20.x+ent",
			activeVersions: defaultActiveVersions,
			changedFiles:   noChangedFiles,
			skip:           true,
		},
		"main to release with mixed changed files": {
			baseRefVersion: "main",
			ref:            "release/1.20.x+ent",
			activeVersions: defaultActiveVersions,
			changedFiles:   mixedCEAndEnterpriseChangedFiles,
			skip:           false,
		},
		"main to release with enterprise only changed files": {
			baseRefVersion: "main",
			ref:            "release/1.20.x+ent",
			activeVersions: defaultActiveVersions,
			changedFiles:   onlyEnterpriseChangedFiles,
			skip:           false,
		},
		"main to release with all CE changed files": {
			baseRefVersion: "main",
			ref:            "release/1.20.x+ent",
			activeVersions: defaultActiveVersions,
			changedFiles:   allCEChangedFiles,
			skip:           false,
		},
		"main to release with allowed inactive changed files": {
			baseRefVersion: "main",
			ref:            "release/1.20.x+ent",
			activeVersions: defaultActiveVersions,
			changedFiles:   allowedInactiveCEChangedFiles,
			skip:           false,
		},
		// release -> active ce/release
		"release to ce/release with no changed files": {
			baseRefVersion: "release/1.20.x",
			ref:            "ce/release/1.20.x",
			activeVersions: defaultActiveVersions,
			changedFiles:   noChangedFiles,
			skip:           true,
		},
		"release to ce/release with mixed changed files": {
			baseRefVersion: "release/1.20.x",
			ref:            "ce/release/1.20.x",
			activeVersions: defaultActiveVersions,
			changedFiles:   mixedCEAndEnterpriseChangedFiles,
			skip:           false,
		},
		"release to ce/release with enterprise only changed files": {
			baseRefVersion: "release/1.20.x",
			ref:            "ce/release/1.20.x",
			activeVersions: defaultActiveVersions,
			changedFiles:   onlyEnterpriseChangedFiles,
			skip:           true,
		},
		"release to ce/release with all CE changed files": {
			baseRefVersion: "release/1.20.x",
			ref:            "ce/release/1.20.x",
			activeVersions: defaultActiveVersions,
			changedFiles:   allCEChangedFiles,
			skip:           false,
		},
		"release to ce/release with allowed inactive changed files": {
			baseRefVersion: "release/1.20.x",
			ref:            "ce/release/1.20.x",
			activeVersions: defaultActiveVersions,
			changedFiles:   allowedInactiveCEChangedFiles,
			skip:           false,
		},
		// release -> inactive ce/release
		"release to inactive ce/release with no changed files": {
			baseRefVersion: "release/1.19.x",
			ref:            "ce/release/1.19.x",
			activeVersions: defaultActiveVersions,
			changedFiles:   noChangedFiles,
			skip:           true,
		},
		"release to inactive ce/release with mixed changed files": {
			baseRefVersion: "release/1.19.x",
			ref:            "ce/release/1.19.x",
			activeVersions: defaultActiveVersions,
			changedFiles:   mixedCEAndEnterpriseChangedFiles,
			skip:           true,
		},
		"release to inactive ce/release with enterprise only changed files": {
			baseRefVersion: "release/1.19.x",
			ref:            "ce/release/1.19.x",
			activeVersions: defaultActiveVersions,
			changedFiles:   onlyEnterpriseChangedFiles,
			skip:           true,
		},
		"release to inactive ce/release with all CE changed files": {
			baseRefVersion: "release/1.19.x",
			ref:            "ce/release/1.19.x",
			activeVersions: defaultActiveVersions,
			changedFiles:   allCEChangedFiles,
			skip:           true,
		},
		"release to inactive ce/release with allowed inactive changed files": {
			baseRefVersion: "release/1.19.x",
			ref:            "ce/release/1.19.x",
			activeVersions: defaultActiveVersions,
			changedFiles:   allowedInactiveCEChangedFiles,
			skip:           false,
		},
		// Various corner cases
		"empty changed files list is skipped": {
			baseRefVersion: "release/1.19.x",
			ref:            "ce/release/1.19.x",
			activeVersions: defaultActiveVersions,
			changedFiles:   noChangedFiles,
			skip:           true,
		},
		"nil changed files list is skipped": {
			baseRefVersion: "release/1.19.x",
			ref:            "ce/release/1.19.x",
			activeVersions: defaultActiveVersions,
			changedFiles:   nil,
			skip:           true,
		},
		"release branch with no active versions": {
			baseRefVersion: "release/1.19.x",
			ref:            "ce/release/1.19.x",
			activeVersions: map[string]*releases.Version{},
			changedFiles:   mixedCEAndEnterpriseChangedFiles,
			skip:           true,
		},
		"release branch with nil active versions": {
			baseRefVersion: "release/1.19.x",
			ref:            "ce/release/1.19.x",
			activeVersions: nil,
			changedFiles:   mixedCEAndEnterpriseChangedFiles,
			skip:           true,
		},
		"missing base ref version": {
			baseRefVersion: "",
			ref:            "ce/main",
			activeVersions: defaultActiveVersions,
			changedFiles:   allCEChangedFiles,
			skip:           true,
		},
		"missing ref version": {
			baseRefVersion: "main",
			ref:            "",
			activeVersions: defaultActiveVersions,
			changedFiles:   allCEChangedFiles,
			skip:           true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := NewCreateBackportReq(
				WithCreateBackportReqAllowInactiveGroups(changed.FileGroups{changed.FileGroup("changelog")}),
			)
			msg, skip := req.shouldSkipRef(
				context.Background(),
				test.baseRefVersion,
				test.ref,
				test.activeVersions,
				test.changedFiles,
			)
			require.Equalf(
				t, test.skip, skip, "should have %t but got %t with %s", test.skip, skip, msg,
			)
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
			// When multiple attempts fail the errors should be stable
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

// Test_filterNonBackportLabels tests the label filtering functionality
func Test_filterNonBackportLabels(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		backportPrefix string
		sourceLabels   Labels
		expectedLabels []string
	}{
		"no labels": {
			backportPrefix: "backport",
			sourceLabels:   Labels{},
			expectedLabels: nil,
		},
		"only backport labels": {
			backportPrefix: "backport",
			sourceLabels: Labels{
				&libgithub.Label{Name: libgithub.Ptr("backport/1.18.x")},
				&libgithub.Label{Name: libgithub.Ptr("backport/1.19.x")},
			},
			expectedLabels: nil,
		},
		"mixed labels": {
			backportPrefix: "backport",
			sourceLabels: Labels{
				&libgithub.Label{Name: libgithub.Ptr("bug")},
				&libgithub.Label{Name: libgithub.Ptr("backport/1.18.x")},
				&libgithub.Label{Name: libgithub.Ptr("enhancement")},
				&libgithub.Label{Name: libgithub.Ptr("backport/ce/main")},
				&libgithub.Label{Name: libgithub.Ptr("docs")},
			},
			expectedLabels: []string{"bug", "enhancement", "docs"},
		},
		"no backport labels": {
			backportPrefix: "backport",
			sourceLabels: Labels{
				&libgithub.Label{Name: libgithub.Ptr("bug")},
				&libgithub.Label{Name: libgithub.Ptr("enhancement")},
				&libgithub.Label{Name: libgithub.Ptr("docs")},
				&libgithub.Label{Name: libgithub.Ptr("priority/high")},
			},
			expectedLabels: []string{"bug", "enhancement", "docs", "priority/high"},
		},
		"custom backport prefix": {
			backportPrefix: "cherry-pick",
			sourceLabels: Labels{
				&libgithub.Label{Name: libgithub.Ptr("bug")},
				&libgithub.Label{Name: libgithub.Ptr("cherry-pick/1.18.x")},
				&libgithub.Label{Name: libgithub.Ptr("enhancement")},
			},
			expectedLabels: []string{"bug", "enhancement"},
		},
		"backport-like but different prefix": {
			backportPrefix: "backport",
			sourceLabels: Labels{
				&libgithub.Label{Name: libgithub.Ptr("backup/daily")},
				&libgithub.Label{Name: libgithub.Ptr("backport/1.18.x")},
				&libgithub.Label{Name: libgithub.Ptr("enhancement")},
			},
			expectedLabels: []string{"backup/daily", "enhancement"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			filteredLabels := filterNonBackportLabels(test.sourceLabels, test.backportPrefix)

			require.Equal(t, test.expectedLabels, filteredLabels,
				"filtered labels should match expected labels")
		})
	}
}

// Test_syncBackportFailedLabel tests that syncBackportFailedLabel applies the
// label when runErr is non-nil, removes it when runErr is nil, and is a no-op
// when label is empty.
func Test_syncBackportFailedLabel(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		label        string
		runErr       error
		addStatus    int  // response status for POST .../labels (apply path)
		removeStatus int  // response status for DELETE .../labels/backport-failed (remove path)
		expectAdd    bool // expect POST .../labels to be called
		expectRemove bool // expect DELETE .../labels/backport-failed to be called
		expectError  bool
	}{
		"run failed - label applied": {
			label:     "backport-failed",
			runErr:    errors.New("something went wrong"),
			addStatus: http.StatusOK,
			expectAdd: true,
		},
		"run succeeded - label removed": {
			label:        "backport-failed",
			runErr:       nil,
			removeStatus: http.StatusOK,
			expectRemove: true,
		},
		"run succeeded - label not present 404 ignored": {
			label:        "backport-failed",
			runErr:       nil,
			removeStatus: http.StatusNotFound,
			expectRemove: true,
		},
		"run failed - api error applying label": {
			label:       "backport-failed",
			runErr:      errors.New("something went wrong"),
			addStatus:   http.StatusInternalServerError,
			expectAdd:   true,
			expectError: true,
		},
		"empty label - no-op on failure": {
			label:  "",
			runErr: errors.New("something went wrong"),
		},
		"empty label - no-op on success": {
			label:  "",
			runErr: nil,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			addCalled := false
			removeCalled := false
			client, mux, teardown := setupTestClient(t)
			defer teardown()

			if test.expectAdd {
				mux.HandleFunc("/repos/test-owner/test-repo/issues/42/labels", func(w http.ResponseWriter, r *http.Request) {
					require.Equal(t, http.MethodPost, r.Method)
					addCalled = true
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(test.addStatus)
					w.Write([]byte(`[]`))
				})
			}

			if test.expectRemove {
				mux.HandleFunc("/repos/test-owner/test-repo/issues/42/labels/"+test.label, func(w http.ResponseWriter, r *http.Request) {
					require.Equal(t, http.MethodDelete, r.Method)
					removeCalled = true
					w.WriteHeader(test.removeStatus)
				})
			}

			req := &CreateBackportReq{
				Owner:               "test-owner",
				Repo:                "test-repo",
				PullNumber:          42,
				BackportFailedLabel: test.label,
			}
			err := req.syncBackportFailedLabel(context.Background(), client, test.runErr)

			if test.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, test.expectAdd, addCalled, "add label API call expectation mismatch")
			require.Equal(t, test.expectRemove, removeCalled, "remove label API call expectation mismatch")
		})
	}
}
