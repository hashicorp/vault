// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"context"
	"strconv"
	"strings"
)

// FetchOpts are the git fetch flags and arguments
// See: https://git-scm.com/docs/git-fetch
type FetchOpts struct {
	// Options
	All         bool // --all
	Atomic      bool // --atomic
	Depth       uint // --depth
	Deepen      uint // --deepen
	Force       bool // --force
	NoTags      bool // --no-tags
	Porcelain   bool // --porcelain
	Progress    bool // --progress
	Prune       bool // --prune
	PruneTags   bool // --prune-tags
	Quiet       bool // --quiet
	SetUpstream bool // --set-upstream
	Unshallow   bool // --unshallow
	Verbose     bool // --verbose

	// Targets
	Repository string   // <repository>
	Refspec    []string // <refspec>
}

// Fetch runs the git fetch command
func (c *Client) Fetch(ctx context.Context, opts *FetchOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "fetch", opts)
}

// String returns the options as a string
func (o *FetchOpts) String() string {
	if o == nil {
		return ""
	}

	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice
func (o *FetchOpts) Strings() []string {
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

	if o.Depth > 0 {
		opts = append(opts, "--depth", strconv.FormatUint(uint64(o.Depth), 10))
	}

	if o.Deepen > 0 {
		opts = append(opts, "--deepen", strconv.FormatUint(uint64(o.Deepen), 10))
	}

	if o.Force {
		opts = append(opts, "--force")
	}

	if o.NoTags {
		opts = append(opts, "--no-tags")
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
