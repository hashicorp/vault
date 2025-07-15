// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"context"
	"fmt"
	"strings"
)

// BranchTrack are supported branch tracking options
type BranchTrack = string

const (
	BranchTrackDirect  BranchTrack = "direct"
	BranchTrackInherit BranchTrack = "inherit"
)

// BranchOpts are the git branch flags and arguments
// See: https://git-scm.com/docs/git-branch
type BranchOpts struct {
	// Options
	Abbrev        uint        // --abbrev=<n>
	All           bool        // --all
	Contains      string      // --contains <commit>
	Copy          bool        // --copy
	CreateReflog  bool        // --create-reflog
	Delete        bool        // --delete
	Force         bool        // -f
	Format        string      // --format <format>
	IgnoreCase    bool        // --ignore-case
	List          bool        // --list
	Merged        string      // --merged <commit>
	Move          bool        // --move
	NoAbbrev      bool        // --no-abbrev
	NoColor       bool        // --no-color
	NoColumn      bool        // --no-column
	NoContains    string      // --no-contains <commit>
	NoMerged      string      // --no-merged <commit>
	NoTrack       bool        // --no-track
	OmitEmpty     bool        // --omit-empty
	PointsAt      string      // --points-at <object>
	Remotes       bool        // --remotes
	Quiet         bool        // --quiet
	SetUpstream   bool        // --set-upstream
	SetUpstreamTo string      // --set-upstream-to=<upstream>
	ShowCurrent   bool        // --show-current
	Sort          string      // --sort=<key>
	Track         BranchTrack // --track
	UnsetUpstream bool        // --unset-upstream

	// Targets. The branch command has several different modules. Set the correct
	// targets depending on which combination of options you're setting.
	BranchName string   // <branchname>
	StartPoint string   // <start-point>
	OldBranch  string   // <oldbranch>
	NewBranch  string   // <newbranch>
	Pattern    []string // <pattern>
}

// Branch runs the git branch command
func (c *Client) Branch(ctx context.Context, opts *BranchOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "branch", opts)
}

// String returns the options as a string
func (o *BranchOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice
func (o *BranchOpts) Strings() []string {
	if o == nil {
		return nil
	}

	opts := []string{}

	if o.Abbrev > 0 {
		opts = append(opts, fmt.Sprintf("--abbrev=%d", o.Abbrev))
	}

	if o.All {
		opts = append(opts, "--all")
	}

	if o.Contains != "" {
		opts = append(opts, fmt.Sprintf("--contains=%s", string(o.Contains)))
	}

	if o.Copy {
		opts = append(opts, "--copy")
	}

	if o.CreateReflog {
		opts = append(opts, "--create-reflog")
	}

	if o.Delete {
		opts = append(opts, "--delete")
	}

	if o.Force {
		opts = append(opts, "--force")
	}

	if o.Format != "" {
		opts = append(opts, fmt.Sprintf("--format=%s", string(o.Format)))
	}

	if o.IgnoreCase {
		opts = append(opts, "--ignore-case")
	}

	if o.List {
		opts = append(opts, "--list")
	}

	if o.Merged != "" {
		opts = append(opts, fmt.Sprintf("--merged=%s", string(o.Merged)))
	}

	if o.Move {
		opts = append(opts, "--move")
	}

	if o.NoAbbrev {
		opts = append(opts, "--no-abbrev")
	}

	if o.NoColor {
		opts = append(opts, "--no-color")
	}

	if o.NoColumn {
		opts = append(opts, "--no-column")
	}

	if o.NoTrack {
		opts = append(opts, "--no-track")
	}

	if o.NoContains != "" {
		opts = append(opts, fmt.Sprintf("--no-contains=%s", string(o.NoContains)))
	}

	if o.NoMerged != "" {
		opts = append(opts, fmt.Sprintf("--no-merged=%s", string(o.NoMerged)))
	}

	if o.OmitEmpty {
		opts = append(opts, "--omit-empty")
	}

	if o.PointsAt != "" {
		opts = append(opts, fmt.Sprintf("--points-at=%s", string(o.PointsAt)))
	}

	if o.Quiet {
		opts = append(opts, "--quiet")
	}

	if o.Remotes {
		opts = append(opts, "--remotes")
	}

	if o.SetUpstream {
		opts = append(opts, "--set-upstream")
	}

	if o.SetUpstreamTo != "" {
		opts = append(opts, fmt.Sprintf("--set-upstream-to=%s", string(o.SetUpstreamTo)))
	}

	if o.ShowCurrent {
		opts = append(opts, "--show-current")
	}

	if o.Sort != "" {
		opts = append(opts, fmt.Sprintf("--sort=%s", string(o.Sort)))
	}

	if o.Track != "" {
		opts = append(opts, fmt.Sprintf("--track=%s", string(o.Track)))
	}

	if o.UnsetUpstream {
		opts = append(opts, "--unset-upstream")
	}

	// Not all of these can be used at once but we try to put them in an order
	// where we won't cause problems if the correct flags and targets are set.

	if o.BranchName != "" {
		opts = append(opts, o.BranchName)
	}

	if o.OldBranch != "" {
		opts = append(opts, o.OldBranch)
	}

	if o.NewBranch != "" {
		opts = append(opts, o.NewBranch)
	}

	if o.StartPoint != "" {
		opts = append(opts, o.StartPoint)
	}

	if len(o.Pattern) > 0 {
		opts = append(opts, o.Pattern...)
	}

	return opts
}
