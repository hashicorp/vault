// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"context"
	"fmt"
	"strings"
)

// EmptyCommit are supported empty commit handling options
type EmptyCommit = string

const (
	EmptyCommitDrop EmptyCommit = "drop"
	EmptyCommitKeep EmptyCommit = "keep"
	EmptyCommitStop EmptyCommit = "stop"
)

// CherryPickOpts are the git cherry-pick flags and arguments
// See: https://git-scm.com/docs/git-cherry-pick
type CherryPickOpts struct {
	// Options
	AllowEmpty         bool                  // --allow-empty
	AllowEmptyMessage  bool                  // --allow-empty-message
	Empty              EmptyCommit           // --empty=
	FF                 bool                  // --ff
	GPGSign            bool                  // --gpgsign
	GPGSignKeyID       string                // --gpgsign=<key-id>
	Mainline           string                // --mainline
	NoReReReAutoupdate bool                  // --no-rerere-autoupdate
	Record             bool                  // -x
	ReReReAutoupdate   bool                  // --rerere-autoupdate
	Signoff            bool                  // --signoff
	Strategy           MergeStrategy         // --strategy
	StrategyOptions    []MergeStrategyOption // --strategy-option=<option>

	// Target
	Commit string // <commit>

	// Sequences
	Continue bool // --continue
	Abort    bool // --abort
	Quit     bool // --quit
}

// CherryPick runs the git cherry-pick command
func (c *Client) CherryPick(ctx context.Context, opts *CherryPickOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "cherry-pick", opts)
}

// CherryPickAbort aborts an in-progress cherry-pick
func (c *Client) CherryPickAbort(ctx context.Context) (*ExecResponse, error) {
	return c.CherryPick(ctx, &CherryPickOpts{Abort: true})
}

// CherryPickContinue continues an in-progress cherry-pick
func (c *Client) CherryPickContinue(ctx context.Context) (*ExecResponse, error) {
	return c.CherryPick(ctx, &CherryPickOpts{Continue: true})
}

// CherryPickQuit quits an in-progress cherry-pick
func (c *Client) CherryPickQuit(ctx context.Context) (*ExecResponse, error) {
	return c.CherryPick(ctx, &CherryPickOpts{Quit: true})
}

// String returns the options as a string
func (o *CherryPickOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice
func (o *CherryPickOpts) Strings() []string {
	if o == nil {
		return nil
	}

	switch {
	case o.Abort:
		return []string{"--abort"}
	case o.Continue:
		return []string{"--continue"}
	case o.Quit:
		return []string{"--quit"}
	}

	opts := []string{}

	if o.AllowEmpty {
		opts = append(opts, "--allow-empty")
	}

	if o.AllowEmptyMessage {
		opts = append(opts, "--allow-empty-message")
	}

	if o.Empty != "" {
		opts = append(opts, fmt.Sprintf("--empty=%s", string(o.Empty)))
	}

	if o.FF {
		opts = append(opts, "--ff")
	}

	if o.GPGSign {
		opts = append(opts, "--gpg-sign")
	}

	if o.GPGSignKeyID != "" {
		opts = append(opts, fmt.Sprintf("--gpg-sign=%s", o.GPGSignKeyID))
	}

	if o.Mainline != "" {
		opts = append(opts, fmt.Sprintf("--mainline=%s", o.Mainline))
	}

	if o.NoReReReAutoupdate {
		opts = append(opts, "--no-rerere-autoupdate")
	}

	if o.Record {
		opts = append(opts, "-x")
	}

	if o.ReReReAutoupdate {
		opts = append(opts, "--rerere-autoupdate")
	}

	if o.Signoff {
		opts = append(opts, "--signoff")
	}

	if o.Strategy != "" {
		opts = append(opts, fmt.Sprintf("--strategy=%s", string(o.Strategy)))
	}

	for _, opt := range o.StrategyOptions {
		opts = append(opts, fmt.Sprintf("--strategy-option=%s", string(opt)))
	}

	if o.Commit != "" {
		opts = append(opts, o.Commit)
	}

	return opts
}
