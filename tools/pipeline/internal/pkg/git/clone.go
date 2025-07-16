// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"context"
	"strconv"
	"strings"
)

// CloneOpts are the git clone flags and arguments
// See: https://git-scm.com/docs/git-clone
type CloneOpts struct {
	// Options
	Branch       string // --branch
	Depth        uint   // --depth
	NoCheckout   bool   // --no-checkout
	NoTags       bool   // --no-tags
	Origin       string // --origin
	Progress     bool   // --progress
	Quiet        bool   // --quiet
	Revision     string // --revision
	SingleBranch bool   // --single-branch
	Sparse       bool   // --sparse
	Verbose      bool   // --verbose

	// Targets
	Repository string // <repository>
	Directory  string // <directory>
}

// Clone runs the git clone command
func (c *Client) Clone(ctx context.Context, opts *CloneOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "clone", opts)
}

// String returns the options as a string
func (o *CloneOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice
func (o *CloneOpts) Strings() []string {
	if o == nil {
		return nil
	}

	opts := []string{}
	if o.Branch != "" {
		opts = append(opts, "--branch", o.Branch)
	}

	if o.Depth > 0 {
		opts = append(opts, "--depth", strconv.FormatUint(uint64(o.Depth), 10))
	}

	if o.NoCheckout {
		opts = append(opts, "--no-checkout")
	}

	if o.NoTags {
		opts = append(opts, "--no-tags")
	}

	if o.Origin != "" {
		opts = append(opts, "--origin", o.Origin)
	}

	if o.Progress {
		opts = append(opts, "--progress")
	}

	if o.Quiet {
		opts = append(opts, "--quiet")
	}

	if o.Revision != "" {
		opts = append(opts, "--revision", o.Revision)
	}

	if o.SingleBranch {
		opts = append(opts, "--single-branch")
	}

	if o.Sparse {
		opts = append(opts, "--sparse")
	}

	if o.Verbose {
		opts = append(opts, "--verbose")
	}

	if o.Repository != "" || o.Directory != "" {
		opts = append(opts, "--")

		if o.Repository != "" {
			opts = append(opts, o.Repository)
		}

		if o.Directory != "" {
			opts = append(opts, o.Directory)
		}
	}

	return opts
}
