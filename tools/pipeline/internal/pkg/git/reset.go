// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"context"
	"strings"
)

// ResetMode is is mode to use when resetting the repository
type ResetMode string

const (
	ResetModeSoft  ResetMode = "soft"
	ResetModeMixed ResetMode = "mixed"
	ResetModeHard  ResetMode = "hard"
	ResetModeMerge ResetMode = "merge"
	ResetModeKeep  ResetMode = "keep"
)

// ResetOpts are the git reset flags and arguments
// See: https://git-scm.com/docs/git-reset
type ResetOpts struct {
	// Options
	Mode      ResetMode // [--soft, --hard, etc..]
	NoRefresh bool      // --no-refresh
	Patch     bool      // --patch
	Quiet     bool      // --quiet
	Refresh   bool      // --refresh

	// Targets
	Commit   string   // <commit>
	Treeish  string   // <tree-ish>
	PathSpec []string // <pathspec>
}

// Reset runs the git reset command
func (c *Client) Reset(ctx context.Context, opts *ResetOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "reset", opts)
}

// String returns the options as a string
func (o *ResetOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice
func (o *ResetOpts) Strings() []string {
	if o == nil {
		return nil
	}

	opts := []string{}

	// Do mode before flags if it set
	if o.Mode != "" {
		opts = append(opts, "--"+string(o.Mode))
	}

	// Flags
	if o.NoRefresh {
		opts = append(opts, "--no-refresh")
	}

	if o.Quiet {
		opts = append(opts, "--quiet")
	}

	if o.Refresh {
		opts = append(opts, "--refresh")
	}

	// Do Patch after flags but before our targets
	if o.Patch {
		opts = append(opts, "--patch")
	}

	// Do our targets
	if o.Commit != "" {
		opts = append(opts, o.Commit)
	}

	if o.Treeish != "" {
		opts = append(opts, o.Treeish)
	}

	// If there's a pathspec, append the paths at the very end
	if len(o.PathSpec) > 0 {
		opts = append(append(opts, "--"), o.PathSpec...)
	}

	return opts
}
