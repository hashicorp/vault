// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package client

import (
	"context"
	"strings"
)

// InitOpts are the git init flags and arguments
// See: https://git-scm.com/docs/git-init
type InitOpts struct {
	// Options
	Bare          bool   // --bare
	InitialBranch string // --initial-branch
	Quiet         bool   // --quiet

	// Targets
	Directory string // <directory>
}

// Init runs the git init command
func (c *Client) Init(ctx context.Context, opts *InitOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "init", opts)
}

// String returns the options as a string
func (o *InitOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice
func (o *InitOpts) Strings() []string {
	if o == nil {
		return nil
	}

	opts := []string{}
	if o.Bare {
		opts = append(opts, "--bare")
	}

	if o.InitialBranch != "" {
		opts = append(opts, "--initial-branch", o.InitialBranch)
	}

	if o.Quiet {
		opts = append(opts, "--quiet")
	}

	if o.Directory != "" {
		opts = append(opts, "--", o.Directory)
	}

	return opts
}
