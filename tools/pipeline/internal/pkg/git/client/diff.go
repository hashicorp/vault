// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package client

import (
	"context"
	"fmt"
	"strings"
)

// DiffOpts are the git diff flags and arguments
// See: https://git-scm.com/docs/git-diff
type DiffOpts struct {
	// Options
	Cached         bool          // --cached
	DiffAlgorithm  DiffAlgorithm // --diff-algorithm=<algo>
	DstPrefix      string        // --dst-prefix=<prefix>
	ExitCode       bool          // --exit-code
	NameOnly       bool          // --name-only
	NameStatus     bool          // --name-status
	NoColor        bool          // --no-color
	NoExtDiff      bool          // --no-ext-diff
	NoPatch        bool          // --no-patch
	NoRenames      bool          // --no-renames
	NullTerminated bool          // -z
	Output         string        // --output=<file>
	Patch          bool          // --patch
	Quiet          bool          // --quiet
	SrcPrefix      string        // --src-prefix=<prefix>
	Stat           bool          // --stat

	// Targets
	Commits  []string // <commit>...
	PathSpec []string // <pathspec>
}

// Diff runs the git diff command
func (c *Client) Diff(ctx context.Context, opts *DiffOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "diff", opts)
}

// String returns the options as a string
func (o *DiffOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice
func (o *DiffOpts) Strings() []string {
	if o == nil {
		return nil
	}

	opts := []string{}
	if o.Cached {
		opts = append(opts, "--cached")
	}

	if o.DiffAlgorithm != "" {
		opts = append(opts, fmt.Sprintf("--diff-algorithm=%s", string(o.DiffAlgorithm)))
	}

	if o.DstPrefix != "" {
		opts = append(opts, fmt.Sprintf("--dst-prefix=%s", o.DstPrefix))
	}

	if o.ExitCode {
		opts = append(opts, "--exit-code")
	}

	if o.NameOnly {
		opts = append(opts, "--name-only")
	}

	if o.NameStatus {
		opts = append(opts, "--name-status")
	}

	if o.NoColor {
		opts = append(opts, "--no-color")
	}

	if o.NoExtDiff {
		opts = append(opts, "--no-ext-diff")
	}

	if o.NoPatch {
		opts = append(opts, "--no-patch")
	}

	if o.NoRenames {
		opts = append(opts, "--no-renames")
	}

	if o.NullTerminated {
		opts = append(opts, "-z")
	}

	if o.Output != "" {
		opts = append(opts, fmt.Sprintf("--output=%s", o.Output))
	}

	if o.Patch {
		opts = append(opts, "--patch")
	}

	if o.Quiet {
		opts = append(opts, "--quiet")
	}

	if o.SrcPrefix != "" {
		opts = append(opts, fmt.Sprintf("--src-prefix=%s", o.SrcPrefix))
	}

	if o.Stat {
		opts = append(opts, "--stat")
	}

	opts = append(opts, o.Commits...)

	// If there's a pathspec, append the paths at the very end
	if len(o.PathSpec) > 0 {
		opts = append(append(opts, "--"), o.PathSpec...)
	}

	return opts
}
