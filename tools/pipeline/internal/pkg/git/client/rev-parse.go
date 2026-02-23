// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package client

import (
	"context"
	"fmt"
	"strings"
)

type (
	// PathFormat specifies the format for path output
	PathFormat = string
)

const (
	PathFormatAbsolute PathFormat = "absolute"
	PathFormatRelative PathFormat = "relative"
)

// RevParseOpts are the git rev-parse flags and arguments
// See: https://git-scm.com/docs/git-rev-parse
type RevParseOpts struct {
	// Operation Modes
	ParseOpt     bool // --parseopt
	SQQuote      bool // --sq-quote
	KeepDashDash bool // --keep-dashdash
	StopAtNonOpt bool // --stop-at-non-option
	StuckLong    bool // --stuck-long
	KeepArgv0    bool // --keep-argv0
	NoRevs       bool // --no-revs
	Revs         bool // --revs
	RevsOnly     bool // --revs-only
	NoFlags      bool // --no-flags
	Flags        bool // --flags

	// Options for Filtering
	All           bool     // --all
	Branches      []string // --branches[=<pattern>]
	Tags          []string // --tags[=<pattern>]
	Remotes       []string // --remotes[=<pattern>]
	Glob          []string // --glob=<pattern>
	Exclude       []string // --exclude=<pattern>
	Disambiguate  string   // --disambiguate=<prefix>
	ExcludeHidden []string // --exclude-hidden[=<pattern>]

	// Options for Output
	AbbrevRef        string // --abbrev-ref[=<mode>]
	Symbolic         bool   // --symbolic
	SymbolicFullName bool   // --symbolic-full-name
	Short            uint   // --short[=<length>]
	Verify           bool   // --verify
	Quiet            bool   // --quiet
	SQ               bool   // --sq
	Not              bool   // --not
	End              bool   // --

	// Options for Objects
	Default string // --default <arg>

	// Options for Files
	LocalEnvVars              bool       // --local-env-vars
	PathFormat                PathFormat // --path-format=<format>
	GitDir                    bool       // --git-dir
	GitCommonDir              bool       // --git-common-dir
	ResolveGitDir             string     // --resolve-git-dir <path>
	GitPath                   string     // --git-path <path>
	ShowTopLevel              bool       // --show-toplevel
	ShowSuperprojectWorkTree  bool       // --show-superproject-working-tree
	SharedIndexPath           bool       // --shared-index-path
	AbsoluteGitDir            bool       // --absolute-git-dir
	IsInsideGitDir            bool       // --is-inside-git-dir
	IsInsideWorkTree          bool       // --is-inside-work-tree
	IsBareRepository          bool       // --is-bare-repository
	IsShallowRepository       bool       // --is-shallow-repository
	ShowCDUp                  bool       // --show-cdup
	ShowPrefix                bool       // --show-prefix
	ShowObjectFormat          bool       // --show-object-format[=<hash-algorithm>]
	ShowObjectFormatAlgorithm string     // hash algorithm for --show-object-format

	// Other Options
	Since  string // --since=<datestring>
	Until  string // --until=<datestring>
	Before string // --before=<datestring>
	After  string // --after=<datestring>

	// Targets
	Args []string // <args>
}

// RevParse runs the git rev-parse command
func (c *Client) RevParse(ctx context.Context, opts *RevParseOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "rev-parse", opts)
}

// String returns the options as a string
func (o *RevParseOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice
func (o *RevParseOpts) Strings() []string {
	if o == nil {
		return nil
	}

	opts := []string{}

	// Operation Modes
	if o.ParseOpt {
		opts = append(opts, "--parseopt")
	}

	if o.SQQuote {
		opts = append(opts, "--sq-quote")
	}

	if o.KeepDashDash {
		opts = append(opts, "--keep-dashdash")
	}

	if o.StopAtNonOpt {
		opts = append(opts, "--stop-at-non-option")
	}

	if o.StuckLong {
		opts = append(opts, "--stuck-long")
	}

	if o.KeepArgv0 {
		opts = append(opts, "--keep-argv0")
	}

	if o.NoRevs {
		opts = append(opts, "--no-revs")
	}

	if o.Revs {
		opts = append(opts, "--revs")
	}

	if o.RevsOnly {
		opts = append(opts, "--revs-only")
	}

	if o.NoFlags {
		opts = append(opts, "--no-flags")
	}

	if o.Flags {
		opts = append(opts, "--flags")
	}

	// Options for Filtering
	if o.All {
		opts = append(opts, "--all")
	}

	for _, branch := range o.Branches {
		if branch == "" {
			opts = append(opts, "--branches")
		} else {
			opts = append(opts, fmt.Sprintf("--branches=%s", branch))
		}
	}

	for _, tag := range o.Tags {
		if tag == "" {
			opts = append(opts, "--tags")
		} else {
			opts = append(opts, fmt.Sprintf("--tags=%s", tag))
		}
	}

	for _, remote := range o.Remotes {
		if remote == "" {
			opts = append(opts, "--remotes")
		} else {
			opts = append(opts, fmt.Sprintf("--remotes=%s", remote))
		}
	}

	for _, pattern := range o.Glob {
		opts = append(opts, fmt.Sprintf("--glob=%s", pattern))
	}

	for _, pattern := range o.Exclude {
		opts = append(opts, fmt.Sprintf("--exclude=%s", pattern))
	}

	if o.Disambiguate != "" {
		opts = append(opts, fmt.Sprintf("--disambiguate=%s", o.Disambiguate))
	}

	for _, pattern := range o.ExcludeHidden {
		if pattern == "" {
			opts = append(opts, "--exclude-hidden")
		} else {
			opts = append(opts, fmt.Sprintf("--exclude-hidden=%s", pattern))
		}
	}

	// Options for Output
	if o.AbbrevRef != "" {
		if o.AbbrevRef == "strict" || o.AbbrevRef == "loose" {
			opts = append(opts, fmt.Sprintf("--abbrev-ref=%s", o.AbbrevRef))
		} else {
			opts = append(opts, "--abbrev-ref")
		}
	}

	if o.Symbolic {
		opts = append(opts, "--symbolic")
	}

	if o.SymbolicFullName {
		opts = append(opts, "--symbolic-full-name")
	}

	if o.Short > 0 {
		opts = append(opts, fmt.Sprintf("--short=%d", o.Short))
	} else if o.Short == 0 {
		// Check if we want --short without a length (uses default)
		// This is a bit tricky - we'll assume if Short is explicitly set to 0
		// and other flags suggest we want short output, we add it
		// For now, we'll skip this case and require explicit length
	}

	if o.Verify {
		opts = append(opts, "--verify")
	}

	if o.Quiet {
		opts = append(opts, "--quiet")
	}

	if o.SQ {
		opts = append(opts, "--sq")
	}

	if o.Not {
		opts = append(opts, "--not")
	}

	// Options for Objects
	if o.Default != "" {
		opts = append(opts, "--default", o.Default)
	}

	// Options for Files
	if o.LocalEnvVars {
		opts = append(opts, "--local-env-vars")
	}

	if o.PathFormat != "" {
		opts = append(opts, fmt.Sprintf("--path-format=%s", string(o.PathFormat)))
	}

	if o.GitDir {
		opts = append(opts, "--git-dir")
	}

	if o.GitCommonDir {
		opts = append(opts, "--git-common-dir")
	}

	if o.ResolveGitDir != "" {
		opts = append(opts, "--resolve-git-dir", o.ResolveGitDir)
	}

	if o.GitPath != "" {
		opts = append(opts, "--git-path", o.GitPath)
	}

	if o.ShowTopLevel {
		opts = append(opts, "--show-toplevel")
	}

	if o.ShowSuperprojectWorkTree {
		opts = append(opts, "--show-superproject-working-tree")
	}

	if o.SharedIndexPath {
		opts = append(opts, "--shared-index-path")
	}

	if o.AbsoluteGitDir {
		opts = append(opts, "--absolute-git-dir")
	}

	if o.IsInsideGitDir {
		opts = append(opts, "--is-inside-git-dir")
	}

	if o.IsInsideWorkTree {
		opts = append(opts, "--is-inside-work-tree")
	}

	if o.IsBareRepository {
		opts = append(opts, "--is-bare-repository")
	}

	if o.IsShallowRepository {
		opts = append(opts, "--is-shallow-repository")
	}

	if o.ShowCDUp {
		opts = append(opts, "--show-cdup")
	}

	if o.ShowPrefix {
		opts = append(opts, "--show-prefix")
	}

	if o.ShowObjectFormat {
		if o.ShowObjectFormatAlgorithm != "" {
			opts = append(opts, fmt.Sprintf("--show-object-format=%s", o.ShowObjectFormatAlgorithm))
		} else {
			opts = append(opts, "--show-object-format")
		}
	}

	// Other Options
	if o.Since != "" {
		opts = append(opts, fmt.Sprintf("--since=%s", o.Since))
	}

	if o.Until != "" {
		opts = append(opts, fmt.Sprintf("--until=%s", o.Until))
	}

	if o.Before != "" {
		opts = append(opts, fmt.Sprintf("--before=%s", o.Before))
	}

	if o.After != "" {
		opts = append(opts, fmt.Sprintf("--after=%s", o.After))
	}

	// Targets - Args should be last
	if len(o.Args) > 0 {
		opts = append(opts, o.Args...)
	}

	return opts
}
