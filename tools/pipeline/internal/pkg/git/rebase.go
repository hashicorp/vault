// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

type (
	// RebaseMerges is the strategy for handling merge commits in rebases
	RebaseMerges = string
	// RebaseMerges is the strategy for rebasing
	RebaseStrategy = string
	// RebaseMerges is the strategy for handling whitespace during rebasing
	WhitespaceAction = string
)

const (
	RebaseMergesCousins   RebaseMerges = "rebase-cousins"
	RebaseMergesNoCousins RebaseMerges = "no-rebase-cousins"

	RebaseStrategyTrue        RebaseStrategy = "true"
	RebaseStrategyFalse       RebaseStrategy = "false"
	RebaseStrategyMerges      RebaseStrategy = "merges"
	RebaseStrategyInteractive RebaseStrategy = "interactive"

	WhitespaceActionNoWarn   WhitespaceAction = "nowarn"
	WhitespaceActionWarn     WhitespaceAction = "warn"
	WhitespaceActionFix      WhitespaceAction = "fix"
	WhitespaceActionError    WhitespaceAction = "error"
	WhitespaceActionErrorAll WhitespaceAction = "error-all"
)

// RebaseOpts are the git rebase flags and arguments
// See: https://git-scm.com/docs/git-rebase
type RebaseOpts struct {
	// Options
	AllowEmptyMessage         bool                  // --allow-empty-message
	Apply                     bool                  // --apply
	Autosquash                bool                  // --autosquash
	Autostash                 bool                  // --autostash
	CommitterDateIsAuthorDate bool                  // --committer-date-is-author-date
	Context                   uint                  // -C
	Empty                     EmptyCommit           // --empty=
	Exec                      string                // --exec=
	ForceRebase               bool                  // --force-rebase
	ForkPoint                 bool                  // --fork-point
	GPGSign                   bool                  // --gpgsign
	GPGSignKeyID              string                // --gpgsign=<key-id>
	IgnoreDate                bool                  // --ignore-date
	IgnoreWhitespace          bool                  // --ignore-whitespace
	KeepBase                  string                // --keep-base <upstream|branch>
	KeepEmpty                 bool                  // --keep-empty
	Merge                     bool                  // --merge
	NoAutosquash              bool                  // --no-autosquash
	NoAutostash               bool                  // --no-autostash
	NoKeepEmpty               bool                  // --no-keep-empty
	NoReapplyCherryPicks      bool                  // --no-reapply-cherry-picks
	NoRebaseMerges            bool                  // --no-rebase-merges
	NoRescheduleFailedExec    bool                  // --no-reschedule-failed-exec
	NoReReReAutoupdate        bool                  // --no-rerere-autoupdate
	NoStat                    bool                  // --no-stat
	NoUpdateRefs              bool                  // --no-update-refs
	NoVerify                  bool                  // --no-verify
	Onto                      string                // --onto
	Quiet                     bool                  // --quiet
	ReapplyCherryPicks        bool                  // --reapply-cherry-picks
	RebaseMerges              RebaseMerges          // --rebase-merges=<strategy>
	RescheduleFailedExec      bool                  // --reschedule-failed-exec
	ResetAuthorDate           bool                  // --reset-author-date
	ReReReAutoupdate          bool                  // --rerere-autoupdate
	Root                      bool                  // --root
	Stat                      bool                  // --stat
	Strategy                  MergeStrategy         // --strategy
	StrategyOptions           []MergeStrategyOption // --strategy-option=<option>
	UpdateRefs                bool                  // --update-refs
	Verbose                   bool                  // --verbose
	Verify                    bool                  // --verify
	Whitespace                WhitespaceAction      //--whitespace=<handler>

	// Args
	Branch string // <branch>

	// Mode Options
	Continue         bool // --continue
	Skip             bool // --skip
	Abort            bool // --abort
	Quit             bool // --quit
	ShowCurrentPatch bool // --show-current-patch
}

// Rebase runs the git rebase command
func (c *Client) Rebase(ctx context.Context, opts *RebaseOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "rebase", opts)
}

// RebaseAbort aborts an in-progress rebase
func (c *Client) RebaseAbort(ctx context.Context) (*ExecResponse, error) {
	return c.Rebase(ctx, &RebaseOpts{Abort: true})
}

// RebaseContinue continues an in-progress rebase
func (c *Client) RebaseContinue(ctx context.Context) (*ExecResponse, error) {
	return c.Rebase(ctx, &RebaseOpts{Continue: true})
}

// RebaseQuit quits an in-progress rebase
func (c *Client) RebaseQuit(ctx context.Context) (*ExecResponse, error) {
	return c.Rebase(ctx, &RebaseOpts{Quit: true})
}

// RebaseSkip skips an in-progress rebase
func (c *Client) RebaseSkip(ctx context.Context) (*ExecResponse, error) {
	return c.Rebase(ctx, &RebaseOpts{Skip: true})
}

// RebaseShowCurrentPatch shows the current patch an in-progress rebase
func (c *Client) RebaseShowCurrentPatch(ctx context.Context) (*ExecResponse, error) {
	return c.Rebase(ctx, &RebaseOpts{ShowCurrentPatch: true})
}

// String returns the options as a string
func (o *RebaseOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the set options as a string slice
func (o *RebaseOpts) Strings() []string {
	if o == nil {
		return nil
	}

	switch {
	case o.Abort:
		return []string{"--abort"}
	case o.Continue:
		return []string{"--continue"}
	case o.Quit:
		return []string{"--quit"}
	case o.Skip:
		return []string{"--skip"}
	case o.ShowCurrentPatch:
		return []string{"--show-current-patch"}
	}

	opts := []string{}

	if o.AllowEmptyMessage {
		opts = append(opts, "--allow-empty-message")
	}
	if o.Apply {
		opts = append(opts, "--apply")
	}

	if o.Autosquash {
		opts = append(opts, "--autosquash")
	}

	if o.Autostash {
		opts = append(opts, "--autostash")
	}

	if o.CommitterDateIsAuthorDate {
		opts = append(opts, "--committer-date-is-author-date")
	}

	if o.Context > 0 {
		opts = append(opts, "-C", strconv.FormatUint(uint64(o.Context), 10))
	}

	if o.Empty != "" {
		opts = append(opts, fmt.Sprintf("--empty=%s", string(o.Empty)))
	}

	if o.Exec != "" {
		opts = append(opts, fmt.Sprintf("--exec=%s", o.Exec))
	}

	if o.ForceRebase {
		opts = append(opts, "--force-rebase")
	}

	if o.ForkPoint {
		opts = append(opts, "--fork-point")
	}

	if o.GPGSign {
		opts = append(opts, "--gpg-sign")
	}

	if o.GPGSignKeyID != "" {
		opts = append(opts, fmt.Sprintf("--gpg-sign=%s", o.GPGSignKeyID))
	}

	if o.IgnoreDate {
		opts = append(opts, "--ignore-date")
	}

	if o.IgnoreWhitespace {
		opts = append(opts, "--ignore-whitespace")
	}

	if o.KeepBase != "" {
		opts = append(opts, fmt.Sprintf("--keep-base=%s", o.KeepBase))
	}

	if o.KeepEmpty {
		opts = append(opts, "--keep-empty")
	}

	if o.Merge {
		opts = append(opts, "--merge")
	}

	if o.NoAutosquash {
		opts = append(opts, "--no-autosquash")
	}

	if o.NoAutostash {
		opts = append(opts, "--no-autostash")
	}

	if o.NoKeepEmpty {
		opts = append(opts, "--no-keep-empty")
	}

	if o.NoReapplyCherryPicks {
		opts = append(opts, "--no-reapply-cherry-picks")
	}

	if o.NoRebaseMerges {
		opts = append(opts, "--no-rebase-merges")
	}

	if o.NoRescheduleFailedExec {
		opts = append(opts, "--no-reschedule-failed-exec")
	}

	if o.NoReReReAutoupdate {
		opts = append(opts, "--no-rerere-autoupdate")
	}

	if o.NoStat {
		opts = append(opts, "--no-stat")
	}

	if o.NoUpdateRefs {
		opts = append(opts, "--no-update-refs")
	}

	if o.NoVerify {
		opts = append(opts, "--no-verify")
	}

	if o.Onto != "" {
		opts = append(opts, fmt.Sprintf("--onto=%s", o.Onto))
	}

	if o.Quiet {
		opts = append(opts, "--quiet")
	}

	if o.ReapplyCherryPicks {
		opts = append(opts, "--reapply-cherry-picks")
	}

	if o.RebaseMerges != "" {
		opts = append(opts, fmt.Sprintf("--rebase-merges=%s", string(o.RebaseMerges)))
	}

	if o.RescheduleFailedExec {
		opts = append(opts, "--reschedule-failed-exec")
	}

	if o.ResetAuthorDate {
		opts = append(opts, "--reset-author-date")
	}

	if o.ReReReAutoupdate {
		opts = append(opts, "--rerere-autoupdate")
	}

	if o.Root {
		opts = append(opts, "--root")
	}

	if o.Stat {
		opts = append(opts, "--stat")
	}

	if o.Strategy != "" {
		opts = append(opts, fmt.Sprintf("--strategy=%s", string(o.Strategy)))
	}

	for _, opt := range o.StrategyOptions {
		opts = append(opts, fmt.Sprintf("--strategy-option=%s", string(opt)))
	}

	if o.UpdateRefs {
		opts = append(opts, "--update-refs")
	}

	if o.Verbose {
		opts = append(opts, "--verbose")
	}

	if o.Verify {
		opts = append(opts, "--verify")
	}

	if o.Whitespace != "" {
		opts = append(opts, fmt.Sprintf("--whitespace=%s", string(o.Whitespace)))
	}

	if o.Branch != "" {
		opts = append(opts, o.Branch)
	}

	return opts
}
