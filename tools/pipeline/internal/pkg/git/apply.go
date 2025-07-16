// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"context"
	"fmt"
	"strings"
)

// ApplyWhitespaceAction are actions Git can take when encountering whitespace
// conflicts during apply.
type ApplyWhitespaceAction = string

const (
	ApplyWhitespaceActionNoWarn   ApplyWhitespaceAction = "nowarn"
	ApplyWhitespaceActionWarn     ApplyWhitespaceAction = "warn"
	ApplyWhitespaceActionFix      ApplyWhitespaceAction = "fix"
	ApplyWhitespaceActionError    ApplyWhitespaceAction = "error"
	ApplyWhitespaceActionErrorAll ApplyWhitespaceAction = "error-all"
)

// ApplyOpts are the git apply flags and arguments
// See: https://git-scm.com/docs/git-apply
type ApplyOpts struct {
	// Options
	AllowEmpty    bool                  // --allow-empty
	Cached        bool                  // --cached
	Check         bool                  // --check
	Index         bool                  // --index
	Ours          bool                  // --ours
	Recount       bool                  // --recount
	Stat          bool                  // --stat
	Summary       bool                  // --summary
	Theirs        bool                  // --theirs
	ThreeWayMerge bool                  // -3way
	Union         bool                  // --union
	Whitespace    ApplyWhitespaceAction // --whitespace=<action>

	// Targets, depending on which combination of options you're setting
	Patch []string // <patch>
}

// Apply runs the git apply command
func (c *Client) Apply(ctx context.Context, opts *ApplyOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "apply", opts)
}

// String returns the options as a string
func (o *ApplyOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice
func (o *ApplyOpts) Strings() []string {
	if o == nil {
		return nil
	}

	opts := []string{}
	if o.AllowEmpty {
		opts = append(opts, "--allow-empty")
	}

	if o.Cached {
		opts = append(opts, "--cached")
	}

	if o.Check {
		opts = append(opts, "--check")
	}

	if o.Index {
		opts = append(opts, "--index")
	}

	if o.Ours {
		opts = append(opts, "--ours")
	}

	if o.Recount {
		opts = append(opts, "--recount")
	}

	if o.Stat {
		opts = append(opts, "--stat")
	}

	if o.Summary {
		opts = append(opts, "--summary")
	}

	if o.Theirs {
		opts = append(opts, "--theirs")
	}

	if o.ThreeWayMerge {
		opts = append(opts, "--3way")
	}

	if o.Union {
		opts = append(opts, "--union")
	}

	if o.Whitespace != "" {
		opts = append(opts, fmt.Sprintf("--whitespace=%s", string(o.Whitespace)))
	}

	if len(o.Patch) > 0 {
		opts = append(opts, o.Patch...)
	}

	return opts
}
