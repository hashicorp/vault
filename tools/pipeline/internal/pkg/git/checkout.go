// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"context"
	"fmt"
	"strings"
)

// CheckoutOpts are the git checkout flags and arguments
// See: https://git-scm.com/docs/git-checkout
type CheckoutOpts struct {
	// Options
	NewBranch              string      // -b
	NewBranchForceCheckout string      // -B
	Detach                 bool        // --detach
	Force                  bool        // -f
	Guess                  bool        // --guess
	Progress               bool        // --progress
	NoTrack                bool        // --no-track
	Orphan                 string      // --orphan
	Ours                   bool        // --ours
	Quiet                  bool        // --quiet
	Theirs                 bool        // --theirs
	Track                  BranchTrack // --track

	// Targets
	Branch     string // <new-branch>
	StartPoint string // <start-point>
	Treeish    string // <tree-ish>

	// Paths
	PathSpec []string // -- <pathspec>
}

// Branch runs the git checkout command
func (c *Client) Checkout(ctx context.Context, opts *CheckoutOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "checkout", opts)
}

// String returns the options as a string
func (o *CheckoutOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice
func (o *CheckoutOpts) Strings() []string {
	if o == nil {
		return nil
	}

	opts := []string{}
	if o.NewBranch != "" {
		opts = append(opts, "-b", o.NewBranch)
	}

	if o.NewBranchForceCheckout != "" {
		opts = append(opts, "-B", o.NewBranchForceCheckout)
	}

	if o.Detach {
		opts = append(opts, "--detach")
	}

	if o.Force {
		opts = append(opts, "--force")
	}

	if o.Guess {
		opts = append(opts, "--guess")
	}

	if o.NoTrack {
		opts = append(opts, "--no-track")
	}

	if o.Orphan != "" {
		opts = append(opts, "--orphan", string(o.Orphan))
	}

	if o.Ours {
		opts = append(opts, "--ours")
	}

	if o.Progress {
		opts = append(opts, "--progress")
	}

	if o.Quiet {
		opts = append(opts, "--quiet")
	}

	if o.Theirs {
		opts = append(opts, "--theirs")
	}

	if o.Track != "" {
		opts = append(opts, fmt.Sprintf("--track=%s", string(o.Track)))
	}

	// Do the <branch>, <start-point>, and <tree-ish> before pathspec
	if o.Branch != "" {
		opts = append(opts, o.Branch)
	}

	if o.StartPoint != "" {
		opts = append(opts, o.StartPoint)
	}

	if o.Treeish != "" {
		opts = append(opts, o.Treeish)
	}

	// If there's a pathspec always set it last
	if len(o.PathSpec) > 0 {
		opts = append(append(opts, "--"), o.PathSpec...)
	}

	return opts
}
