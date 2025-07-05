// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"context"
	"strings"
)

// RemoteOpts are the git remote sub-commands, flags and arguments
// See: https://git-scm.com/docs/git-remote
type RemoteOpts struct {
	// Sub-command
	Command RemoteCommand

	// Flags
	Verbose    bool     // --verbose
	Fetch      bool     // -f
	Tags       bool     // --tags
	NoTags     bool     // --no-tags
	Master     string   // -m
	Track      []string // -t
	Delete     bool     // --delete
	Auto       bool     // --auto
	Add        bool     // --add
	Push       bool     // --push
	All        bool     // --all
	NoQuery    bool     // -n
	DryRun     bool     // --dry-run
	Prune      bool     // --prune
	Progress   bool     // --progress
	NoProgress bool     // --no-progress

	// Targets
	Old      string   // <old>
	New      string   // <new>
	NewURL   string   // <newurl>
	OldURL   string   // <oldurl>
	Name     string   // <name>
	Names    []string // <name>...
	Branch   string   // <branch>
	Branches []string // <branch>...
	URL      string   // <URL>
}

type RemoteCommand string

const (
	RemoteCommandAdd         RemoteCommand = "add"
	RemoteCommandRename      RemoteCommand = "rename"
	RemoteCommandRemove      RemoteCommand = "remove"
	RemoteCommandSetHead     RemoteCommand = "set-head"
	RemoteCommandSetBranches RemoteCommand = "set-branches"
	RemoteCommandGetURL      RemoteCommand = "get-url"
	RemoteCommandSetURL      RemoteCommand = "set-url"
	RemoteCommandShow        RemoteCommand = "show"
	RemoteCommandPrune       RemoteCommand = "prune"
	RemoteCommandUpdate      RemoteCommand = "update"
)

// Remote runs the git rm command
func (c *Client) Remote(ctx context.Context, opts *RemoteOpts) (*ExecResponse, error) {
	// TODO: Handle exit code 2 and 3?
	// https://git-scm.com/docs/git-remote#_exit_status
	return c.Exec(ctx, "remote", opts)
}

// String returns the options as a string
func (o *RemoteOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice
func (o *RemoteOpts) Strings() []string {
	if o == nil {
		return nil
	}

	opts := []string{}

	switch o.Command {
	case RemoteCommandAdd:
		opts = append(opts, string(o.Command))
		if o.Fetch {
			opts = append(opts, "-f")
		}
		if o.Tags {
			opts = append(opts, "--tags")
		}
		if o.NoTags {
			opts = append(opts, "--no-tags")
		}
		if o.Master != "" {
			opts = append(opts, "-m", o.Master)
		}
		for _, branch := range o.Track {
			opts = append(opts, "-t", branch)
		}
		opts = append(opts, o.Name, o.URL)
	case RemoteCommandRename:
		opts = append(opts, string(o.Command))
		if o.Progress {
			opts = append(opts, "--progress")
		}
		if o.NoProgress {
			opts = append(opts, "--no-progress")
		}
		opts = append(opts, o.Old, o.New)
	case RemoteCommandRemove:
		opts = append(opts, string(o.Command), o.Name)
	case RemoteCommandSetHead:
		opts = append(opts, string(o.Command), o.Name)
		if o.Auto {
			opts = append(opts, "--auto")
		}
		if o.Delete {
			opts = append(opts, "--delete")
		}
		if o.Branch != "" {
			opts = append(opts, o.Branch)
		}
	case RemoteCommandSetBranches:
		opts = append(opts, string(o.Command))
		if o.Add {
			opts = append(opts, "--add")
		}
		opts = append(opts, o.Name)
		if o.Branch != "" {
			opts = append(opts, o.Branch)
		}
		if len(o.Branches) > 0 {
			opts = append(opts, o.Branches...)
		}
	case RemoteCommandGetURL:
		opts = append(opts, string(o.Command))
		if o.Push {
			opts = append(opts, "--push")
		}
		if o.All {
			opts = append(opts, "--all")
		}
		opts = append(opts, o.Name)
	case RemoteCommandSetURL:
		opts = append(opts, string(o.Command))
		if o.Add {
			opts = append(opts, "--add")
		}
		if o.Delete {
			opts = append(opts, "--delete")
		}
		if o.Push {
			opts = append(opts, "--push")
		}
		opts = append(opts, o.Name)
		if o.NewURL != "" {
			opts = append(opts, o.NewURL)
		}
		if o.OldURL != "" {
			opts = append(opts, o.OldURL)
		}
		if o.URL != "" {
			opts = append(opts, o.URL)
		}
	case RemoteCommandShow:
		if o.Verbose {
			opts = append(opts, "-v")
		}
		opts = append(opts, string(o.Command))
		if o.NoQuery {
			opts = append(opts, "-n")
		}
		if o.Name != "" {
			opts = append(opts, o.Name)
		}
		if len(o.Names) > 0 {
			opts = append(opts, o.Names...)
		}
	case RemoteCommandPrune:
		opts = append(opts, string(o.Command))
		if o.NoQuery {
			opts = append(opts, "-n")
		}
		if o.DryRun {
			opts = append(opts, "--dry-run")
		}
		if o.Name != "" {
			opts = append(opts, o.Name)
		}
		if len(o.Names) > 0 {
			opts = append(opts, o.Names...)
		}
	case RemoteCommandUpdate:
		if o.Verbose {
			opts = append(opts, "-v")
		}
		opts = append(opts, string(o.Command))
		if o.Prune {
			opts = append(opts, "--prune")
		}
	default:
		// git remote
		if o.Verbose {
			opts = append(opts, "-v")
		}
	}

	return opts
}
