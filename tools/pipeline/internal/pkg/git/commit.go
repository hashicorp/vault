// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"context"
	"fmt"
	"strings"
)

type (
	// Cleanup are cleanup modes
	Cleanup string
	// FixupLog configures how to fix up a commit
	FixupLog string
	// UntrackedFiles are how to show untracked files
	UntrackedFiles string
)

const (
	CommitCleanupModeString     Cleanup        = "strip"
	CommitCleanupModeWhitespace Cleanup        = "whitespace"
	CommitCleanupModeVerbatim   Cleanup        = "verbatim"
	CommitCleanupModeScissors   Cleanup        = "scissors"
	CommitCleanupModeDefault    Cleanup        = "default"
	CommitFixupLogNone          FixupLog       = ""
	CommitFixupLogAmend         FixupLog       = "amend"
	CommitFixupLogReword        FixupLog       = "reword"
	UntrackedFilesNo            UntrackedFiles = "no"
	UntrackedFilesNormal        UntrackedFiles = "normal"
	UntrackedFilesAll           UntrackedFiles = "all"
)

// CommitFixup is how to fixup a commit
type CommitFixup struct {
	FixupLog
	Commit string
}

// CommitOpts are the git commit flags and arguments
// See: https://git-scm.com/docs/git-commit
type CommitOpts struct {
	// Options
	All               bool         // --all
	AllowEmpty        bool         // --allow-empty
	AllowEmptyMessage bool         // --allow-empty-message
	Amend             bool         // --amend
	Author            string       // --author=<author>
	Branch            bool         // --branch
	Cleanup           Cleanup      // --cleanup=<mode>
	Date              string       // --date=<date>
	DryRun            bool         // --dry-run
	File              string       // --file=<file>
	Fixup             *CommitFixup // --fixup=
	GPGSign           bool         // --gpgsign
	GPGSignKeyID      string       // --gpgsign=<key-id>
	Long              bool         // --long
	Patch             bool         // --patch
	Porcelain         bool         // --porcelain
	Message           string       // --message=<message>
	NoEdit            bool         // --no-edit
	NoPostRewrite     bool         // --no-post-rewrite
	NoVerify          bool         // --no-verify
	Null              bool         // --null
	Only              bool         // --only
	Quiet             bool         // --quiet
	ResetAuthor       bool         // --reset-author
	ReuseMessage      string       // --reuse-message=<commit>
	Short             bool         // --short
	Signoff           bool         // --signoff
	Status            bool         // --status
	Verbose           bool         // --verbose

	// Target
	PathSpec []string // <pathspec>
}

// Commit runs the git commit command
func (c *Client) Commit(ctx context.Context, opts *CommitOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "commit", opts)
}

// String returns the options as a string
func (o *CommitOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice
func (o *CommitOpts) Strings() []string {
	if o == nil {
		return nil
	}

	opts := []string{}
	if o.All {
		opts = append(opts, "--all")
	}

	if o.AllowEmpty {
		opts = append(opts, "--allow-empty")
	}

	if o.AllowEmptyMessage {
		opts = append(opts, "--allow-empty-message")
	}

	if o.Amend {
		opts = append(opts, "--amend")
	}

	if o.Author != "" {
		opts = append(opts, fmt.Sprintf("--author=%s", o.Author))
	}

	if o.Branch {
		opts = append(opts, "--branch")
	}

	if o.Cleanup != "" {
		opts = append(opts, fmt.Sprintf("--cleanup=%s", string(o.Cleanup)))
	}

	if o.Date != "" {
		opts = append(opts, fmt.Sprintf("--date=%s", o.Date))
	}

	if o.DryRun {
		opts = append(opts, "--dry-run")
	}

	if o.File != "" {
		opts = append(opts, fmt.Sprintf("--file=%s", o.File))
	}

	if o.Fixup != nil {
		if o.Fixup.FixupLog == CommitFixupLogNone {
			opts = append(opts, fmt.Sprintf("--fixup=%s", string(o.Fixup.Commit)))
		} else {
			opts = append(opts, fmt.Sprintf("--fixup=%s:%s", string(o.Fixup.FixupLog), string(o.Fixup.Commit)))
		}
	}

	if o.GPGSign {
		opts = append(opts, "--gpg-sign")
	}

	if o.GPGSignKeyID != "" {
		opts = append(opts, fmt.Sprintf("--gpg-sign=%s", o.GPGSignKeyID))
	}

	if o.Long {
		opts = append(opts, "--long")
	}

	if o.Patch {
		opts = append(opts, "--patch")
	}

	if o.Porcelain {
		opts = append(opts, "--porcelain")
	}

	if o.Message != "" {
		opts = append(opts, fmt.Sprintf("--message=%s", o.Message))
	}

	if o.NoEdit {
		opts = append(opts, "--no-edit")
	}

	if o.NoPostRewrite {
		opts = append(opts, "--no-post-rewrite")
	}

	if o.NoVerify {
		opts = append(opts, "--no-verify")
	}

	if o.Null {
		opts = append(opts, "--null")
	}

	if o.Only {
		opts = append(opts, "--only")
	}

	if o.Quiet {
		opts = append(opts, "--quiet")
	}

	if o.ResetAuthor {
		opts = append(opts, "--reset-author")
	}

	if o.ReuseMessage != "" {
		opts = append(opts, fmt.Sprintf("--reuse-message=%s", o.ReuseMessage))
	}

	if o.Short {
		opts = append(opts, "--short")
	}

	if o.Status {
		opts = append(opts, "--status")
	}

	if o.Verbose {
		opts = append(opts, "--verbose")
	}

	// If there's a pathspec, append the paths at the very end
	if len(o.PathSpec) > 0 {
		opts = append(append(opts, "--"), o.PathSpec...)
	}

	return opts
}
