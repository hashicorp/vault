// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// RecurseSubmodules is a sub-module recurse mode
type RecurseSubmodules = string

const (
	RecurseSubmodulesYes      RecurseSubmodules = "yes"
	RecurseSubmodulesOnDemand RecurseSubmodules = "on-demand"
	RecurseSubmodulesNo       RecurseSubmodules = "no"
)

// PullOpts are the git pull flags and arguments
// See: https://git-scm.com/docs/git-pull
type PullOpts struct {
	// Options
	Quiet               bool              // --quiet
	Verbose             bool              // --verbose
	RecurseSubmodules   RecurseSubmodules // --recurse-submodules=
	NoRecurseSubmodules bool              // --no-recurse-submodules

	// Merge options
	Autostash               bool                  // --autostash
	AllowUnrelatedHistories bool                  // --allow-unrelated-histories
	DoCommit                bool                  // --commit
	NoDoCommit              bool                  // --no-commit
	Cleanup                 Cleanup               // --cleanup=
	FF                      bool                  // --ff
	FFOnly                  bool                  // --ff-onnly
	NoFF                    bool                  // --no-ff
	GPGSign                 bool                  // --gpgsign
	GPGSignKeyID            string                // --gpgsign=<key-id>
	Log                     uint                  // --log=
	NoAutostash             bool                  // --no-autostash
	NoLog                   bool                  // --no-log
	NoRebase                bool                  // --no-rebase
	NoStat                  bool                  // --no-stat
	NoSquash                bool                  // --no-squash
	NoVerify                bool                  // --no-verify
	Stat                    bool                  // --stat
	Squash                  bool                  // --squash
	Strategy                MergeStrategy         // --stategy=
	StrategyOptions         []MergeStrategyOption // --strategy-option=
	Rebase                  RebaseStrategy        // --rebase=
	Verify                  bool                  // --verify

	// Fetch options
	All           bool // --all
	Append        bool // --append
	Atomic        bool // --atomic
	Depth         uint // --depth
	Deepen        uint // --deepen
	Force         bool // --force
	NoTags        bool // --no-tags
	Porcelain     bool // --porcelain
	Progress      bool // --progress
	Prune         bool // --prune
	PruneTags     bool // --prune-tags
	SetUpstream   bool // --set-upstream
	Unshallow     bool // --unshallow
	UpdateShallow bool // --update-shallow

	// Targets
	Repository string   // <repository>
	Refspec    []string // <refspec>
}

// Pull runs the git pull command
func (c *Client) Pull(ctx context.Context, opts *PullOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "pull", opts)
}

// String returns the options as a string
func (o *PullOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice
func (o *PullOpts) Strings() []string {
	if o == nil {
		return nil
	}

	opts := []string{}

	if o.All {
		opts = append(opts, "--all")
	}

	if o.Atomic {
		opts = append(opts, "--atomic")
	}

	if o.Autostash {
		opts = append(opts, "--autostash")
	}

	if o.DoCommit {
		opts = append(opts, "--commit")
	}

	if o.Depth > 0 {
		opts = append(opts, "--depth", strconv.FormatUint(uint64(o.Depth), 10))
	}

	if o.Deepen > 0 {
		opts = append(opts, "--deepen", strconv.FormatUint(uint64(o.Deepen), 10))
	}

	if o.FF {
		opts = append(opts, "--ff")
	}

	if o.FFOnly {
		opts = append(opts, "--ff-only")
	}

	if o.Force {
		opts = append(opts, "--force")
	}

	if o.GPGSign {
		opts = append(opts, "--gpg-sign")
	}

	if o.GPGSignKeyID != "" {
		opts = append(opts, fmt.Sprintf("--gpg-sign=%s", o.GPGSignKeyID))
	}

	if o.Log > 0 {
		opts = append(opts, fmt.Sprintf("--log=%d", o.Log))
	}

	if o.Squash {
		opts = append(opts, "--squash")
	}

	if o.Stat {
		opts = append(opts, "--stat")
	}

	for _, opt := range o.StrategyOptions {
		opts = append(opts, "-X", string(opt))
	}

	if o.NoAutostash {
		opts = append(opts, "--no-autostash")
	}

	if o.NoDoCommit {
		opts = append(opts, "--no-commit")
	}

	if o.NoFF {
		opts = append(opts, "--no-ff")
	}

	if o.NoLog {
		opts = append(opts, "--no-log")
	}

	if o.NoRebase {
		opts = append(opts, "--no-rebase")
	}

	if o.NoRecurseSubmodules {
		opts = append(opts, "--no-recurse-submodules")
	}

	if o.NoSquash {
		opts = append(opts, "--no-squash")
	}

	if o.NoStat {
		opts = append(opts, "--no-stat")
	}

	if o.NoStat {
		opts = append(opts, "--no-stat")
	}

	if o.NoTags {
		opts = append(opts, "--no-tags")
	}

	if o.NoVerify {
		opts = append(opts, "--no-verify")
	}

	if o.Porcelain {
		opts = append(opts, "--porcelain")
	}

	if o.Progress {
		opts = append(opts, "--progress")
	}

	if o.Prune {
		opts = append(opts, "--prune")
	}

	if o.PruneTags {
		opts = append(opts, "--prune-tags")
	}

	if o.Quiet {
		opts = append(opts, "--quiet")
	}

	if o.Rebase != "" {
		opts = append(opts, fmt.Sprintf("--rebase=%s", string(o.Rebase)))
	}

	if o.SetUpstream {
		opts = append(opts, "--set-upstream")
	}

	if o.Unshallow {
		opts = append(opts, "--unshallow")
	}

	if o.Verbose {
		opts = append(opts, "--verbose")
	}

	if o.Repository != "" {
		opts = append(opts, o.Repository)
	}

	if len(o.Refspec) > 0 {
		opts = append(opts, o.Refspec...)
	}

	return opts
}
