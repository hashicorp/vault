// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"context"
	"fmt"
	"strings"
)

type (
	// MergeStrategy are the merge strategy to use
	MergeStrategy = string
	// MergeStrategy are merge strategy options
	MergeStrategyOption = string
)

const (
	MergeStrategyORT       MergeStrategy = "ort"
	MergeStrategyRecursive MergeStrategy = "recursive"
	MergeStrategyResolve   MergeStrategy = "resolve"
	MergeStrategyOctopus   MergeStrategy = "octopus"
	MergeStrategyOurs      MergeStrategy = "ours"
	MergeStrategySubtree   MergeStrategy = "subtree"

	// Ort
	MergeStrategyOptionOurs              MergeStrategy = "ours"
	MergeStrategyOptionTheirs            MergeStrategy = "theirs"
	MergeStrategyOptionIgnoreSpaceChange MergeStrategy = "ignore-space-change"
	MergeStrategyOptionIgnoreAllSpace    MergeStrategy = "ignore-all-space"
	MergeStrategyOptionIgnoreSpaceAtEOL  MergeStrategy = "ignore-space-at-eol"
	MergeStrategyOptionIgnoreCRAtEOL     MergeStrategy = "ignore-cr-at-eol"
	MergeStrategyOptionRenormalize       MergeStrategy = "renormalize"
	MergeStrategyOptionNoRenormalize     MergeStrategy = "no-renormalize"
	MergeStrategyOptionFindRenames       MergeStrategy = "find-renames"

	// Recursive
	MergeStrategyOptionDiffAlgorithmPatience  MergeStrategy = "diff-algorithm=patience"
	MergeStrategyOptionDiffAlgorithmMinimal   MergeStrategy = "diff-algorithm=minimal"
	MergeStrategyOptionDiffAlgorithmHistogram MergeStrategy = "diff-algorithm=histogram"
	MergeStrategyOptionDiffAlgorithmMyers     MergeStrategy = "diff-algorithm=myers"
)

// MergeOpts are the git merge flags and arguments
// See: https://git-scm.com/docs/git-merge
type MergeOpts struct {
	// Options
	Autostash          bool                  // --autostash
	DoCommit           bool                  // --commit
	File               string                // --file=<file>
	FF                 bool                  // --ff
	FFOnly             bool                  // --ff-onnly
	IntoName           string                // --into-name
	Log                uint                  // --log=<n>
	Message            string                // -m
	NoAutostash        bool                  // --no-autostash
	NoDoCommit         bool                  // --no-commit
	NoFF               bool                  // --no-ff
	NoLog              bool                  // --no-log
	NoOverwrite        bool                  // --no-overwrite
	NoProgress         bool                  // --no-progress
	NoRebase           bool                  // --no-rebase
	NoReReReAutoupdate bool                  // --no-rerere-autoupdate
	NoSquash           bool                  // --no-squash
	NoStat             bool                  // --no-stat
	NoVerify           bool                  // --no-verify
	Progress           bool                  // --progress
	Quiet              bool                  // --quiet
	ReReReAutoupdate   bool                  // --rerere-autoupdate
	Squash             bool                  // --squash
	Stat               bool                  // --stat
	Strategy           MergeStrategy         // --stategy=<strategy>
	StrategyOptions    []MergeStrategyOption // --strategy-option=<option>
	Verbose            bool                  // --verbose

	// Targets
	Commit string // <commit>

	// Sequences
	Continue bool // --continue
	Abort    bool // --abort
	Quit     bool // --quit
}

// Merge runs the git merge command
func (c *Client) Merge(ctx context.Context, opts *MergeOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "merge", opts)
}

// MergeAbort aborts an in-progress merge
func (c *Client) MergeAbort(ctx context.Context) (*ExecResponse, error) {
	return c.Merge(ctx, &MergeOpts{Abort: true})
}

// MergeContinue continues an in-progress merge
func (c *Client) MergeContinue(ctx context.Context) (*ExecResponse, error) {
	return c.Merge(ctx, &MergeOpts{Continue: true})
}

// MergeQuit quits an in-progress merge
func (c *Client) MergeQuit(ctx context.Context) (*ExecResponse, error) {
	return c.Merge(ctx, &MergeOpts{Quit: true})
}

// String returns the options as a string
func (m *MergeOpts) String() string {
	return strings.Join(m.Strings(), " ")
}

// Strings returns the options as a string slice
func (m *MergeOpts) Strings() []string {
	if m == nil {
		return nil
	}

	switch {
	case m.Abort:
		return []string{"--abort"}
	case m.Continue:
		return []string{"--continue"}
	case m.Quit:
		return []string{"--quit"}
	}

	opts := []string{}

	if m.Autostash {
		opts = append(opts, "--autostash")
	}

	if m.DoCommit {
		opts = append(opts, "--commit")
	}

	if m.File != "" {
		opts = append(opts, fmt.Sprintf("--file=%s", m.File))
	}

	if m.FF {
		opts = append(opts, "--ff")
	}

	if m.FFOnly {
		opts = append(opts, "--ff-only")
	}

	if m.IntoName != "" {
		opts = append(opts, "--into-name", m.IntoName)
	}

	if m.Log > 0 {
		opts = append(opts, fmt.Sprintf("--log=%d", m.Log))
	}

	if m.ReReReAutoupdate {
		opts = append(opts, "--rerere-autoupdate")
	}

	if m.Squash {
		opts = append(opts, "--squash")
	}

	if m.Stat {
		opts = append(opts, "--stat")
	}

	if m.Strategy != "" {
		opts = append(opts, fmt.Sprintf("--strategy=%s", string(m.Strategy)))
	}

	for _, opt := range m.StrategyOptions {
		opts = append(opts, fmt.Sprintf("--strategy-option=%s", string(opt)))
	}

	if m.NoAutostash {
		opts = append(opts, "--no-autostash")
	}

	if m.NoDoCommit {
		opts = append(opts, "--no-commit")
	}

	if m.NoFF {
		opts = append(opts, "--no-ff")
	}

	if m.NoLog {
		opts = append(opts, "--no-log")
	}

	if m.NoProgress {
		opts = append(opts, "--no-progress")
	}

	if m.NoRebase {
		opts = append(opts, "--no-rebase")
	}

	if m.NoReReReAutoupdate {
		opts = append(opts, "--no-rerere-autoupdate")
	}

	if m.NoSquash {
		opts = append(opts, "--no-squash")
	}

	if m.NoStat {
		opts = append(opts, "--no-stat")
	}

	if m.NoStat {
		opts = append(opts, "--no-stat")
	}

	if m.NoVerify {
		opts = append(opts, "--no-verify")
	}

	if m.Commit != "" {
		opts = append(opts, m.Commit)
	}

	return opts
}
