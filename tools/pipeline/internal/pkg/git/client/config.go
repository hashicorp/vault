// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package client

import (
	"context"
	"fmt"
	"strings"
)

// ConfigAction is the git config subcommand.
// See: https://git-scm.com/docs/git-config
type ConfigAction string

const (
	// ConfigActionGet reads the value of a single key (git config get).
	ConfigActionGet ConfigAction = "get"
	// ConfigActionSet writes a value for a key (git config set).
	ConfigActionSet ConfigAction = "set"
	// ConfigActionUnset removes a key (git config unset).
	ConfigActionUnset ConfigAction = "unset"
	// ConfigActionList lists all key=value pairs (git config list).
	ConfigActionList ConfigAction = "list"
	// ConfigActionRenameSection renames a config section (git config rename-section).
	ConfigActionRenameSection ConfigAction = "rename-section"
	// ConfigActionRemoveSection removes a config section (git config remove-section).
	ConfigActionRemoveSection ConfigAction = "remove-section"
)

// ConfigType constrains how git config canonicalizes values on read and write.
type ConfigType string

const (
	ConfigTypeBool       ConfigType = "bool"
	ConfigTypeInt        ConfigType = "int"
	ConfigTypeBoolOrInt  ConfigType = "bool-or-int"
	ConfigTypePath       ConfigType = "path"
	ConfigTypeExpiryDate ConfigType = "expiry-date"
	ConfigTypeColor      ConfigType = "color"
)

// ConfigScope limits which config file git config reads from or writes to.
type ConfigScope string

const (
	ConfigScopeSystem   ConfigScope = "system"
	ConfigScopeGlobal   ConfigScope = "global"
	ConfigScopeLocal    ConfigScope = "local"
	ConfigScopeWorktree ConfigScope = "worktree"
)

// ConfigOpts are the flags and arguments for git config.
// See: https://git-scm.com/docs/git-config
type ConfigOpts struct {
	// Action is the subcommand: get, set, unset, list, rename-section,
	// remove-section. Required.
	Action ConfigAction

	// Key is the config key for get/set/unset (e.g. "user.name"). Required
	// for get, set, and unset.
	Key string

	// Value is the value to write for set.
	Value string

	// OldName / NewName are the section names for rename-section.
	OldName string
	NewName string

	// SectionName is the section to remove for remove-section.
	SectionName string

	// --- Scope / file-location flags (mutually exclusive) ---

	// Scope limits reads/writes to the given config scope (system, global,
	// local, worktree). Corresponds to --system / --global / --local /
	// --worktree.
	Scope ConfigScope

	// File reads from or writes to the given path instead of the default
	// config file. Corresponds to -f / --file.
	File string

	// Blob reads config from a blob ref instead of a file (e.g.
	// "master:.gitmodules"). Corresponds to --blob.
	Blob string

	// --- Query modifiers ---

	// All returns all values for a multi-valued key (get) or replaces all of
	// them (set/unset). Corresponds to --all.
	All bool

	// Regexp interprets Key as an extended regular expression when used with
	// get. Corresponds to --regexp.
	Regexp bool

	// URL performs url-matching on the key (e.g. section.URL.key lookup).
	// Corresponds to --url=<URL>.
	URL string

	// ValuePattern filters set/unset to only entries whose value matches this
	// extended regular expression. Corresponds to --value=<pattern>.
	ValuePattern string

	// FixedValue treats ValuePattern as a literal string instead of a regexp.
	// Corresponds to --fixed-value.
	FixedValue bool

	// Default is the fallback value emitted by get when the key is absent.
	// Corresponds to --default=<value>.
	Default string

	// Type canonicalizes the value for get/set. Corresponds to --type=<type>.
	Type ConfigType

	// --- Output modifiers ---

	// Null terminates values with NUL instead of newline. Corresponds to
	// -z / --null.
	Null bool

	// NameOnly outputs only key names, not values. Corresponds to --name-only.
	NameOnly bool

	// ShowNames outputs key names in addition to values for get. Corresponds
	// to --show-names.
	ShowNames bool

	// ShowOrigin augments output with the config file origin. Corresponds to
	// --show-origin.
	ShowOrigin bool

	// ShowScope augments output with the config scope. Corresponds to
	// --show-scope.
	ShowScope bool

	// Includes respects include.* directives. Corresponds to --includes.
	Includes bool

	// Append adds a new line without altering existing values (set).
	// Corresponds to --append.
	Append bool

	// Comment appends a comment to new or modified lines (set).
	// Corresponds to --comment=<message>.
	Comment string

	// ReplaceAll replaces all matching lines (set). Corresponds to
	// --replace-all.
	ReplaceAll bool
}

// String returns the options as a string.
func (o *ConfigOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice.
func (o *ConfigOpts) Strings() []string {
	if o == nil {
		return nil
	}

	opts := []string{}

	// Subcommand comes first.
	if o.Action != "" {
		opts = append(opts, string(o.Action))
	}

	// Scope flags.
	switch o.Scope {
	case ConfigScopeSystem:
		opts = append(opts, "--system")
	case ConfigScopeGlobal:
		opts = append(opts, "--global")
	case ConfigScopeLocal:
		opts = append(opts, "--local")
	case ConfigScopeWorktree:
		opts = append(opts, "--worktree")
	}

	if o.File != "" {
		opts = append(opts, fmt.Sprintf("--file=%s", o.File))
	}

	if o.Blob != "" {
		opts = append(opts, fmt.Sprintf("--blob=%s", o.Blob))
	}

	// Write modifiers.
	if o.ReplaceAll {
		opts = append(opts, "--replace-all")
	}

	if o.Append {
		opts = append(opts, "--append")
	}

	if o.Comment != "" {
		opts = append(opts, fmt.Sprintf("--comment=%s", o.Comment))
	}

	// Query modifiers.
	if o.All {
		opts = append(opts, "--all")
	}

	if o.Regexp {
		opts = append(opts, "--regexp")
	}

	if o.URL != "" {
		opts = append(opts, fmt.Sprintf("--url=%s", o.URL))
	}

	if o.ValuePattern != "" {
		opts = append(opts, fmt.Sprintf("--value=%s", o.ValuePattern))
	}

	if o.FixedValue {
		opts = append(opts, "--fixed-value")
	}

	if o.Default != "" {
		opts = append(opts, fmt.Sprintf("--default=%s", o.Default))
	}

	if o.Type != "" {
		opts = append(opts, fmt.Sprintf("--type=%s", string(o.Type)))
	}

	// Output modifiers.
	if o.Null {
		opts = append(opts, "--null")
	}

	if o.NameOnly {
		opts = append(opts, "--name-only")
	}

	if o.ShowNames {
		opts = append(opts, "--show-names")
	}

	if o.ShowOrigin {
		opts = append(opts, "--show-origin")
	}

	if o.ShowScope {
		opts = append(opts, "--show-scope")
	}

	if o.Includes {
		opts = append(opts, "--includes")
	}

	// Positional arguments — order depends on action.
	switch o.Action {
	case ConfigActionGet, ConfigActionSet, ConfigActionUnset:
		if o.Key != "" {
			opts = append(opts, o.Key)
		}
		if o.Action == ConfigActionSet && o.Value != "" {
			opts = append(opts, o.Value)
		}
	case ConfigActionRenameSection:
		if o.OldName != "" {
			opts = append(opts, o.OldName)
		}
		if o.NewName != "" {
			opts = append(opts, o.NewName)
		}
	case ConfigActionRemoveSection:
		if o.SectionName != "" {
			opts = append(opts, o.SectionName)
		}
	}

	return opts
}

// Config runs the git config command.
func (c *Client) Config(ctx context.Context, opts *ConfigOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "config", opts)
}

// ConfigGet is a convenience wrapper around Config that reads a single key and
// returns its trimmed string value.
func (c *Client) ConfigGet(ctx context.Context, key string) (string, error) {
	res, err := c.Config(ctx, &ConfigOpts{
		Action: ConfigActionGet,
		Key:    key,
	})
	if err != nil {
		return "", fmt.Errorf("git config get %s: %w", key, err)
	}

	return strings.TrimSpace(string(res.Stdout)), nil
}
