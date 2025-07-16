// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"context"
	"fmt"
	"strings"
)

// AmOpts are the git am flags and arguments
// See: https://git-scm.com/docs/git-am
type AmOpts struct {
	// Options
	CommitterDateIsAuthorDate bool                  // --committer-date-is-author-date
	Empty                     EmptyCommit           // --empty=<mode>
	Keep                      bool                  // --keep
	KeepNonPatch              bool                  // --keep-non-patch
	MessageID                 bool                  // --message-id
	NoMessageID               bool                  // --no-message-id
	NoReReReAutoupdate        bool                  // --no-rerere-autoupdate
	NoVerify                  bool                  // --no-verify
	Quiet                     bool                  // --quiet
	ReReReAutoupdate          bool                  // --rerere-autoupdate
	Signoff                   bool                  // --signoff
	ThreeWayMerge             bool                  // --3way
	Whitespace                ApplyWhitespaceAction // --whitespace=<action>

	// Targets, depending on which combination of options you're setting
	Mbox []string // <mbox|Maildir>

	// Sequences
	Abort    bool // --abort
	Continue bool // --continue
	Quit     bool // --quit
	Resolved bool // --resolved
	Retry    bool // --retry

	// Options that are allowed on sequences
	AllowEmpty bool // --allow-empty
}

// Am runs the git am command
func (c *Client) Am(ctx context.Context, opts *AmOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "am", opts)
}

// String returns the options as a string
func (o *AmOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice
func (o *AmOpts) Strings() []string {
	if o == nil {
		return nil
	}

	opts := []string{}

	switch {
	case o.Abort:
		return append(opts, "--abort")
	case o.Continue:
		return append(opts, "--continue")
	case o.Quit:
		return append(opts, "--quit")
	case o.Resolved:
		if o.AllowEmpty {
			opts = append(opts, "--allow-empty")
		}
		return append(opts, "--resolved")
	case o.Retry:
		return append(opts, "--retry")
	}

	if o.CommitterDateIsAuthorDate {
		opts = append(opts, "--committer-date-is-author-date")
	}

	if o.Empty != "" {
		opts = append(opts, fmt.Sprintf("--empty=%s", string(o.Empty)))
	}

	if o.Keep {
		opts = append(opts, "--keep")
	}

	if o.KeepNonPatch {
		opts = append(opts, "--keep-non-patch")
	}

	if o.MessageID {
		opts = append(opts, "--message-id")
	}

	if o.NoMessageID {
		opts = append(opts, "--no-message-id")
	}

	if o.NoReReReAutoupdate {
		opts = append(opts, "--no-rerere-autoupdate")
	}

	if o.NoVerify {
		opts = append(opts, "--no-verify")
	}

	if o.Quiet {
		opts = append(opts, "--quiet")
	}

	if o.ReReReAutoupdate {
		opts = append(opts, "--rerere-autoupdate")
	}

	if o.Signoff {
		opts = append(opts, "--signoff")
	}

	if o.ThreeWayMerge {
		opts = append(opts, "--3way")
	}

	if o.Whitespace != "" {
		opts = append(opts, fmt.Sprintf("--whitespace=%s", string(o.Whitespace)))
	}

	if len(o.Mbox) > 0 {
		opts = append(opts, o.Mbox...)
	}

	return opts
}
