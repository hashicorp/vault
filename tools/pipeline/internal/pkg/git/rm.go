// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"context"
	"strings"
)

// RmOpts are the git rm flags and arguments
// See: https://git-scm.com/docs/git-rm
type RmOpts struct {
	Cached          bool // --cached
	DryRun          bool // --dry-run
	Force           bool // --force
	IgnoreUnmatched bool // --ignore-unmatched
	Quiet           bool // --quiet
	Recursive       bool // -r
	Sparse          bool // --sparse

	PathSpec []string // <pathspec>
}

// Rm runs the git rm command
func (c *Client) Rm(ctx context.Context, opts *RmOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "rm", opts)
}

// String returns the options as a string
func (o *RmOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice
func (o *RmOpts) Strings() []string {
	if o == nil {
		return nil
	}

	opts := []string{}
	if o.Cached {
		opts = append(opts, "--cached")
	}

	if o.DryRun {
		opts = append(opts, "--dry-run")
	}

	if o.Force {
		opts = append(opts, "--force")
	}

	if o.IgnoreUnmatched {
		opts = append(opts, "--ignore-unmatched")
	}

	if o.Quiet {
		opts = append(opts, "--quiet")
	}

	if o.Recursive {
		opts = append(opts, "-r")
	}

	if o.Sparse {
		opts = append(opts, "--sparse")
	}

	// If there's a pathspec, append the paths at the very end
	if len(o.PathSpec) > 0 {
		opts = append(append(opts, "--"), o.PathSpec...)
	}

	return opts
}
