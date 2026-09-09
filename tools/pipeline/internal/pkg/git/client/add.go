// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package client

import (
	"context"
	"strings"
)

// AddOpts are the git add flags and arguments.
// See: https://git-scm.com/docs/git-add
type AddOpts struct {
	All                bool   // -A / --all / --no-ignore-removal
	Chmod              string // --chmod=(+/-)x
	DryRun             bool   // -n / --dry-run
	Edit               bool   // -e / --edit
	Force              bool   // -f / --force
	IgnoreErrors       bool   // --ignore-errors
	IgnoreMissing      bool   // --ignore-missing (only valid with --dry-run)
	IntentToAdd        bool   // -N / --intent-to-add
	Interactive        bool   // -i / --interactive
	NoAll              bool   // --no-all / --ignore-removal
	NoWarnEmbeddedRepo bool   // --no-warn-embedded-repo
	Patch              bool   // -p / --patch
	// PathspecFromFile reads pathspec from the given file; when set, PathSpec
	// is ignored because the two input modes are mutually exclusive.
	PathspecFromFile string // --pathspec-from-file=<file>
	// PathspecFileNul is only meaningful alongside PathspecFromFile.
	PathspecFileNul bool // --pathspec-file-nul
	Refresh         bool // --refresh
	Renormalize     bool // --renormalize
	Verbose         bool // -v / --verbose

	PathSpec []string // <pathspec> — ignored when PathspecFromFile is set
}

// Add runs the git add command.
func (c *Client) Add(ctx context.Context, opts *AddOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "add", opts)
}

// String returns the options as a string.
func (o *AddOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice.
// Flags are emitted in alphabetical order, matching the convention used by
// other opts in this package. PathSpec is suppressed when PathspecFromFile
// is set, because the two input modes are mutually exclusive.
func (o *AddOpts) Strings() []string {
	if o == nil {
		return nil
	}

	var opts []string

	if o.All {
		opts = append(opts, "--all")
	}

	if o.Chmod != "" {
		opts = append(opts, "--chmod="+o.Chmod)
	}

	if o.DryRun {
		opts = append(opts, "--dry-run")
	}

	if o.Edit {
		opts = append(opts, "--edit")
	}

	if o.Force {
		opts = append(opts, "--force")
	}

	if o.IgnoreErrors {
		opts = append(opts, "--ignore-errors")
	}

	if o.IgnoreMissing {
		opts = append(opts, "--ignore-missing")
	}

	if o.IntentToAdd {
		opts = append(opts, "--intent-to-add")
	}

	if o.Interactive {
		opts = append(opts, "--interactive")
	}

	if o.NoAll {
		opts = append(opts, "--no-all")
	}

	if o.NoWarnEmbeddedRepo {
		opts = append(opts, "--no-warn-embedded-repo")
	}

	if o.Patch {
		opts = append(opts, "--patch")
	}

	if o.PathspecFromFile != "" {
		opts = append(opts, "--pathspec-from-file="+o.PathspecFromFile)
		// --pathspec-file-nul is only valid alongside --pathspec-from-file.
		if o.PathspecFileNul {
			opts = append(opts, "--pathspec-file-nul")
		}
		// PathSpec is ignored: git does not accept positional paths when
		// --pathspec-from-file is in use.
		return opts
	}

	if o.Refresh {
		opts = append(opts, "--refresh")
	}

	if o.Renormalize {
		opts = append(opts, "--renormalize")
	}

	if o.Verbose {
		opts = append(opts, "--verbose")
	}

	if len(o.PathSpec) > 0 {
		opts = append(append(opts, "--"), o.PathSpec...)
	}

	return opts
}
