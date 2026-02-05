// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package client

import (
	"context"
	"fmt"
	"strings"
)

type (
	// LogPrettyFormat is the format for pretty printing
	LogPrettyFormat = string
	// LogDateFormat is the format for dates
	LogDateFormat = string
	// LogDecorateFormat is the format for decoration
	LogDecorateFormat = string
	// LogDiffFilter is the filter for diff types
	LogDiffFilter = string
)

const (
	LogPrettyFormatOneline   LogPrettyFormat = "oneline"
	LogPrettyFormatShort     LogPrettyFormat = "short"
	LogPrettyFormatMedium    LogPrettyFormat = "medium"
	LogPrettyFormatFull      LogPrettyFormat = "full"
	LogPrettyFormatFuller    LogPrettyFormat = "fuller"
	LogPrettyFormatReference LogPrettyFormat = "reference"
	LogPrettyFormatEmail     LogPrettyFormat = "email"
	LogPrettyFormatRaw       LogPrettyFormat = "raw"
	LogPrettyFormatNone      LogPrettyFormat = "none" // NOTE: renders blank value to support --pretty=

	LogDateFormatRelative LogDateFormat = "relative"
	LogDateFormatISO      LogDateFormat = "iso"
	LogDateFormatISO8601  LogDateFormat = "iso8601"
	LogDateFormatRFC      LogDateFormat = "rfc"
	LogDateFormatShort    LogDateFormat = "short"
	LogDateFormatRaw      LogDateFormat = "raw"
	LogDateFormatHuman    LogDateFormat = "human"
	LogDateFormatUnix     LogDateFormat = "unix"

	LogDecorateFormatShort LogDecorateFormat = "short"
	LogDecorateFull        LogDecorateFormat = "full"
	LogDecorateAuto        LogDecorateFormat = "auto"
	LogDecorateNo          LogDecorateFormat = "no"

	LogDiffFilterAdded       LogDiffFilter = "A"
	LogDiffFilterCopied      LogDiffFilter = "C"
	LogDiffFilterDeleted     LogDiffFilter = "D"
	LogDiffFilterModified    LogDiffFilter = "M"
	LogDiffFilterRenamed     LogDiffFilter = "R"
	LogDiffFilterTypeChanged LogDiffFilter = "T"
	LogDiffFilterUnmerged    LogDiffFilter = "U"
	LogDiffFilterUnknown     LogDiffFilter = "X"
	LogDiffFilterBroken      LogDiffFilter = "B"
	LogDiffFilterAll         LogDiffFilter = "*"
)

// LogOpts are the git log flags and arguments
// See: https://git-scm.com/docs/git-log
type LogOpts struct {
	// Commit Limiting
	MaxCount         uint   // -n, --max-count=<number>
	Skip             uint   // --skip=<number>
	Since            string // --since=<date>
	After            string // --after=<date>
	Until            string // --until=<date>
	Before           string // --before=<date>
	Author           string // --author=<pattern>
	Committer        string // --committer=<pattern>
	Grep             string // --grep=<pattern>
	AllMatch         bool   // --all-match
	InvertGrep       bool   // --invert-grep
	RegexpIgnoreCase bool   // -i, --regexp-ignore-case

	// Merge Options
	Merges      bool // --merges
	NoMerges    bool // --no-merges
	FirstParent bool // --first-parent

	// History Traversal
	All      bool     // --all
	Branches []string // --branches[=<pattern>]
	Tags     []string // --tags[=<pattern>]
	Remotes  []string // --remotes[=<pattern>]

	// Formatting
	Oneline        bool              // --oneline
	Pretty         LogPrettyFormat   // --pretty=<format>
	Format         string            // --format=<string>
	AbbrevCommit   bool              // --abbrev-commit
	NoAbbrevCommit bool              // --no-abbrev-commit
	Abbrev         uint              // --abbrev=<n>
	Decorate       LogDecorateFormat // --decorate[=<format>]
	DecorateRefs   []string          // --decorate-refs=<pattern>
	Source         bool              // --source
	Graph          bool              // --graph
	Date           LogDateFormat     // --date=<format>
	RelativeDate   bool              // --relative-date

	// Diff Options
	Patch         bool            // -p, --patch
	NoPatch       bool            // -s, --no-patch
	Stat          bool            // --stat
	Shortstat     bool            // --shortstat
	NameOnly      bool            // --name-only
	NameStatus    bool            // --name-status
	DiffFilter    []LogDiffFilter // --diff-filter=<filter>
	DiffMerges    DiffMergeFormat // --diff-merges=<format>
	CombinedDiff  bool            // -c
	DenseCombined bool            // --cc
	Follow        bool            // --follow
	FullDiff      bool            // --full-diff

	// Ordering
	DateOrder       bool // --date-order
	AuthorDateOrder bool // --author-date-order
	TopoOrder       bool // --topo-order
	Reverse         bool // --reverse

	// History Simplification
	SimplifyByDecoration bool // --simplify-by-decoration
	FullHistory          bool // --full-history
	AncestryPath         bool // --ancestry-path
	ShowPulls            bool // --show-pulls

	// Reflog
	WalkReflogs bool // -g, --walk-reflogs

	// Output Control
	Color   bool // --color
	NoColor bool // --no-color
	NullSep bool // -z

	// Targets
	Target   string   // <revision range> - can be a range (A..B), branch name, or commit
	PathSpec []string // -- <path>
}

// Log runs the git log command
func (c *Client) Log(ctx context.Context, opts *LogOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "log", opts)
}

// String returns the options as a string
func (o *LogOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice
func (o *LogOpts) Strings() []string {
	if o == nil {
		return nil
	}

	opts := []string{}

	// Targets - Target (revision range, branch, or commit)
	if o.Target != "" {
		opts = append(opts, o.Target)
	}

	// Commit Limiting
	if o.MaxCount > 0 {
		opts = append(opts, fmt.Sprintf("--max-count=%d", o.MaxCount))
	}

	if o.Skip > 0 {
		opts = append(opts, fmt.Sprintf("--skip=%d", o.Skip))
	}

	if o.Since != "" {
		opts = append(opts, fmt.Sprintf("--since=%s", o.Since))
	}

	if o.After != "" {
		opts = append(opts, fmt.Sprintf("--after=%s", o.After))
	}

	if o.Until != "" {
		opts = append(opts, fmt.Sprintf("--until=%s", o.Until))
	}

	if o.Before != "" {
		opts = append(opts, fmt.Sprintf("--before=%s", o.Before))
	}

	if o.Author != "" {
		opts = append(opts, fmt.Sprintf("--author=%s", o.Author))
	}

	if o.Committer != "" {
		opts = append(opts, fmt.Sprintf("--committer=%s", o.Committer))
	}

	if o.Grep != "" {
		opts = append(opts, fmt.Sprintf("--grep=%s", o.Grep))
	}

	if o.AllMatch {
		opts = append(opts, "--all-match")
	}

	if o.InvertGrep {
		opts = append(opts, "--invert-grep")
	}

	if o.RegexpIgnoreCase {
		opts = append(opts, "--regexp-ignore-case")
	}

	// Merge Options
	if o.Merges {
		opts = append(opts, "--merges")
	}

	if o.NoMerges {
		opts = append(opts, "--no-merges")
	}

	if o.FirstParent {
		opts = append(opts, "--first-parent")
	}

	// History Traversal
	if o.All {
		opts = append(opts, "--all")
	}

	for _, branch := range o.Branches {
		opts = append(opts, fmt.Sprintf("--branches=%s", branch))
	}

	for _, tag := range o.Tags {
		opts = append(opts, fmt.Sprintf("--tags=%s", tag))
	}

	for _, remote := range o.Remotes {
		opts = append(opts, fmt.Sprintf("--remotes=%s", remote))
	}

	// Formatting
	if o.Oneline {
		opts = append(opts, "--oneline")
	}

	if o.Pretty != "" {
		if o.Pretty == LogPrettyFormatNone {
			opts = append(opts, "--pretty=")
		} else {
			opts = append(opts, fmt.Sprintf("--pretty=%s", string(o.Pretty)))
		}
	}

	if o.Format != "" {
		opts = append(opts, fmt.Sprintf("--format=%s", o.Format))
	}

	if o.AbbrevCommit {
		opts = append(opts, "--abbrev-commit")
	}

	if o.NoAbbrevCommit {
		opts = append(opts, "--no-abbrev-commit")
	}

	if o.Abbrev > 0 {
		opts = append(opts, fmt.Sprintf("--abbrev=%d", o.Abbrev))
	}

	if o.Decorate != "" {
		opts = append(opts, fmt.Sprintf("--decorate=%s", string(o.Decorate)))
	}

	for _, ref := range o.DecorateRefs {
		opts = append(opts, fmt.Sprintf("--decorate-refs=%s", ref))
	}

	if o.Source {
		opts = append(opts, "--source")
	}

	if o.Graph {
		opts = append(opts, "--graph")
	}

	if o.Date != "" {
		opts = append(opts, fmt.Sprintf("--date=%s", string(o.Date)))
	}

	if o.RelativeDate {
		opts = append(opts, "--relative-date")
	}

	// Diff Options
	if o.Patch {
		opts = append(opts, "--patch")
	}

	if o.NoPatch {
		opts = append(opts, "--no-patch")
	}

	if o.Stat {
		opts = append(opts, "--stat")
	}

	if o.Shortstat {
		opts = append(opts, "--shortstat")
	}

	if o.NameOnly {
		opts = append(opts, "--name-only")
	}

	if o.NameStatus {
		opts = append(opts, "--name-status")
	}

	if len(o.DiffFilter) > 0 {
		filters := make([]string, len(o.DiffFilter))
		for i, filter := range o.DiffFilter {
			filters[i] = string(filter)
		}
		opts = append(opts, fmt.Sprintf("--diff-filter=%s", strings.Join(filters, "")))
	}

	if o.DiffMerges != "" {
		opts = append(opts, fmt.Sprintf("--diff-merges=%s", string(o.DiffMerges)))
	}

	if o.CombinedDiff {
		opts = append(opts, "-c")
	}

	if o.DenseCombined {
		opts = append(opts, "--cc")
	}

	if o.Follow {
		opts = append(opts, "--follow")
	}

	if o.FullDiff {
		opts = append(opts, "--full-diff")
	}

	// Ordering
	if o.DateOrder {
		opts = append(opts, "--date-order")
	}

	if o.AuthorDateOrder {
		opts = append(opts, "--author-date-order")
	}

	if o.TopoOrder {
		opts = append(opts, "--topo-order")
	}

	if o.Reverse {
		opts = append(opts, "--reverse")
	}

	// History Simplification
	if o.SimplifyByDecoration {
		opts = append(opts, "--simplify-by-decoration")
	}

	if o.FullHistory {
		opts = append(opts, "--full-history")
	}

	if o.AncestryPath {
		opts = append(opts, "--ancestry-path")
	}

	if o.ShowPulls {
		opts = append(opts, "--show-pulls")
	}

	// Reflog
	if o.WalkReflogs {
		opts = append(opts, "--walk-reflogs")
	}

	// Output Control
	if o.Color {
		opts = append(opts, "--color")
	}

	if o.NoColor {
		opts = append(opts, "--no-color")
	}

	if o.NullSep {
		opts = append(opts, "-z")
	}

	// PathSpec - must be last with -- separator
	if len(o.PathSpec) > 0 {
		opts = append(append(opts, "--"), o.PathSpec...)
	}

	return opts
}
