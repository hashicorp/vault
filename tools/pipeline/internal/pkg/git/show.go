// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"context"
	"fmt"
	"strings"
)

type (
	DiffAlgorithm   = string
	DiffMergeFormat = string
)

const (
	DiffAlgorithmPatience  DiffAlgorithm = "patience"
	DiffAlgorithmMinimal   DiffAlgorithm = "minimal"
	DiffAlgorithmHistogram DiffAlgorithm = "histogram"
	DiffAlgorithmMyers     DiffAlgorithm = "myers"

	DiffMergeFormatOff           DiffMergeFormat = "off"
	DiffMergeFormatNone          DiffMergeFormat = "none"
	DiffMergeFormatFirstParent   DiffMergeFormat = "first-parent"
	DiffMergeFormatSeparate      DiffMergeFormat = "separate"
	DiffMergeFormatCombined      DiffMergeFormat = "combined"
	DiffMergeFormatDenseCombined DiffMergeFormat = "dense-combined"
	DiffMergeFormatRemerge       DiffMergeFormat = "remerge"
)

// ShowOpts are the git show flags and arguments
// See: https://git-scm.com/docs/git-show
type ShowOpts struct {
	// Options
	DiffAlgorithm DiffAlgorithm   // --diff-algorithm=<algo>
	DiffMerges    DiffMergeFormat // --diff-merges=<format>
	Format        string          // --format <format>
	NoColor       bool            // --no-color
	NoPatch       bool            // --no-patch
	Patch         bool            // --patch
	Output        string          // --output=<file>
	Raw           bool            // --raw

	// Targets
	Object   string   // <object>
	PathSpec []string // <pathspec>
}

// Show runs the git show command
func (c *Client) Show(ctx context.Context, opts *ShowOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "show", opts)
}

// String returns the options as a string
func (o *ShowOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice
func (o *ShowOpts) Strings() []string {
	if o == nil {
		return nil
	}

	opts := []string{}

	if o.DiffAlgorithm != "" {
		opts = append(opts, fmt.Sprintf("--diff-algorithm=%s", string(o.DiffAlgorithm)))
	}

	if o.DiffMerges != "" {
		opts = append(opts, fmt.Sprintf("--diff-merges=%s", string(o.DiffMerges)))
	}

	if o.Format != "" {
		opts = append(opts, fmt.Sprintf("--format=%s", string(o.Format)))
	}

	if o.NoColor {
		opts = append(opts, "--no-color")
	}

	if o.NoPatch {
		opts = append(opts, "--no-patch")
	}

	if o.Output != "" {
		opts = append(opts, fmt.Sprintf("--output=%s", o.Output))
	}

	if o.Patch {
		opts = append(opts, "--patch")
	}

	if o.Raw {
		opts = append(opts, "--raw")
	}

	opts = append(opts, o.Object)

	// If there's a pathspec, append the paths at the very end
	if len(o.PathSpec) > 0 {
		opts = append(append(opts, "--"), o.PathSpec...)
	}

	return opts
}
