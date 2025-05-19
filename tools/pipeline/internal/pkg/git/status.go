// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"context"
	"fmt"
	"strings"
)

type (
	// IgnoredMode determines how to handle ignored files
	IgnoredMode = string
	// IgnoredMode determines how to handle changes to submodules
	IgnoreSubmodulesWhen = string
)

const (
	IgnoredModeTraditional IgnoredMode = "traditional"
	IgnoredModeNo          IgnoredMode = "no"
	IgnoredModeMatching    IgnoredMode = "matching"

	IgnoreSubmodulesWhenNone      IgnoreSubmodulesWhen = "none"
	IgnoreSubmodulesWhenUntracked IgnoreSubmodulesWhen = "untracked"
	IgnoreSubmodulesWhenDirty     IgnoreSubmodulesWhen = "dirty"
	IgnoreSubmodulesWhenAll       IgnoreSubmodulesWhen = "all"
)

// StatusOpts are the git status flags and arguments
// See: https://git-scm.com/docs/git-status
type StatusOpts struct {
	// Options
	AheadBehind      bool                 // --ahead-behind
	Branch           bool                 // --branch
	Column           string               // --column=
	FindRenames      uint                 // --find-renames=
	Ignored          IgnoredMode          // --ignored=
	IgnoreSubmodules IgnoreSubmodulesWhen // --ignore-submodules=<when>
	Long             bool                 // --long
	NoAheadBehind    bool                 // --no-ahead-behind
	NoColumn         bool                 // --no-column
	NoRenames        bool                 // --no-renames
	Porcelain        bool                 // --porcelain
	Renames          bool                 // --renames
	Short            bool                 // --short
	ShowStash        bool                 // --show-stash
	UntrackedFiles   UntrackedFiles       // --untracked-files=<mode>
	Verbose          bool                 // --verbose

	// Targets
	PathSpec []string // <pathspec>
}

// Status runs the git status command
func (c *Client) Status(ctx context.Context, opts *StatusOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "status", opts)
}

// String returns the options as a string
func (o *StatusOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice
func (o *StatusOpts) Strings() []string {
	if o == nil {
		return nil
	}

	opts := []string{}
	if o.AheadBehind {
		opts = append(opts, "--ahead-behind")
	}

	if o.Branch {
		opts = append(opts, "--branch")
	}

	if o.Column != "" {
		opts = append(opts, fmt.Sprintf("--column=%s", o.Column))
	}

	if o.FindRenames > 0 {
		opts = append(opts, fmt.Sprintf("--find-renames=%d", o.FindRenames))
	}

	if o.Ignored != "" {
		opts = append(opts, fmt.Sprintf("--ignored=%s", string(o.Ignored)))
	}

	if o.IgnoreSubmodules != "" {
		opts = append(opts, fmt.Sprintf("--ignore-submodules=%s", string(o.IgnoreSubmodules)))
	}

	if o.Long {
		opts = append(opts, "--long")
	}

	if o.NoAheadBehind {
		opts = append(opts, "--no-ahead-behind")
	}

	if o.NoColumn {
		opts = append(opts, "--no-column")
	}

	if o.NoRenames {
		opts = append(opts, "--no-renames")
	}

	if o.Porcelain {
		opts = append(opts, "--porcelain")
	}

	if o.Renames {
		opts = append(opts, "--renames")
	}

	if o.Short {
		opts = append(opts, "--short")
	}

	if o.ShowStash {
		opts = append(opts, "--show-stash")
	}

	if o.UntrackedFiles != "" {
		opts = append(opts, fmt.Sprintf("--untracked-files=%s", string(o.UntrackedFiles)))
	}

	if o.Verbose {
		opts = append(opts, "--verbose")
	}

	// If there's a pathspec, append the paths at the very end
	if len(o.PathSpec) > 0 {
		opts = append(append(opts, "--"), o.PathSpec...)
	}

	return opts
}
