// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package client

import (
	"context"
	"fmt"
	"strings"
)

// LsFilesOpts are the git ls-files flags and arguments
// See: https://git-scm.com/docs/git-ls-files
type LsFilesOpts struct {
	// Options
	Cached            bool   // --cached
	Deleted           bool   // --deleted
	Modified          bool   // --modified
	Others            bool   // --others
	Ignored           bool   // --ignored
	Stage             bool   // --stage
	Directory         bool   // --directory
	NoEmptyDirectory  bool   // --no-empty-directory
	Unmerged          bool   // --unmerged
	Killed            bool   // --killed
	NullTerminated    bool   // -z
	ExcludeStandard   bool   // --exclude-standard
	ErrorUnmatch      bool   // --error-unmatch
	FullName          bool   // --full-name
	RecurseSubmodules bool   // --recurse-submodules
	Format            string // --format=<format>

	// Targets
	PathSpec []string // <pathspec>
}

// LsFiles runs the git ls-files command
func (c *Client) LsFiles(ctx context.Context, opts *LsFilesOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "ls-files", opts)
}

// String returns the options as a string
func (o *LsFilesOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice
func (o *LsFilesOpts) Strings() []string {
	if o == nil {
		return nil
	}

	opts := []string{}
	if o.Cached {
		opts = append(opts, "--cached")
	}

	if o.Deleted {
		opts = append(opts, "--deleted")
	}

	if o.Modified {
		opts = append(opts, "--modified")
	}

	if o.Others {
		opts = append(opts, "--others")
	}

	if o.Ignored {
		opts = append(opts, "--ignored")
	}

	if o.Stage {
		opts = append(opts, "--stage")
	}

	if o.Directory {
		opts = append(opts, "--directory")
	}

	if o.NoEmptyDirectory {
		opts = append(opts, "--no-empty-directory")
	}

	if o.Unmerged {
		opts = append(opts, "--unmerged")
	}

	if o.Killed {
		opts = append(opts, "--killed")
	}

	if o.NullTerminated {
		opts = append(opts, "-z")
	}

	if o.ExcludeStandard {
		opts = append(opts, "--exclude-standard")
	}

	if o.ErrorUnmatch {
		opts = append(opts, "--error-unmatch")
	}

	if o.FullName {
		opts = append(opts, "--full-name")
	}

	if o.RecurseSubmodules {
		opts = append(opts, "--recurse-submodules")
	}

	if o.Format != "" {
		opts = append(opts, fmt.Sprintf("--format=%s", o.Format))
	}

	// If there's a pathspec, append the paths at the very end
	if len(o.PathSpec) > 0 {
		opts = append(append(opts, "--"), o.PathSpec...)
	}

	return opts
}
