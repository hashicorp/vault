// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package changed

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Config represents the configuration for categorizing changed files into groups.
type Config struct {
	Groups []*GroupConfig `hcl:"group,block" json:"groups,omitempty"`
}

// GroupConfig defines a named group with matchers to include or exclude files.
type GroupConfig struct {
	Name   string   `hcl:"name,label" json:"name,omitempty"`
	Ignore Matchers `hcl:"ignore,block" json:"ignore,omitempty"`
	Match  Matchers `hcl:"match,block" json:"match,omitempty"`
}

// Matchers is a collection of match configurations used to determine if a file matches criteria.
type Matchers []*MatchConfig

// MatchConfig defines the criteria for matching file paths.
type MatchConfig struct {
	// File is the entire file path relative to the shared working directory
	File []string `hcl:"file,optional" json:"file,omitempty"`
	// BaseDir is the file path's base directory
	BaseDir []string `hcl:"base_dir,optional" json:"base_dir,omitempty"`
	// BaseName is the files base name
	BaseName []string `hcl:"base_name,optional" json:"base_name,omitempty"`
	// BaseNamePrefix is the file's base name prefix
	BaseNamePrefix []string `hcl:"base_name_prefix,optional" json:"base_name_prefix,omitempty"`
	// Contains are string matches in the path
	Contains []string `hcl:"contains,optional" json:"contains,omitempty"`
	// Extension is the file's extension
	Extension []string `hcl:"extension,optional" json:"extension,omitempty"`
}

// FileGroups evaluates a file against all configured groups and returns the groups it belongs to.
func (c *Config) FileGroups(ctx context.Context, file *File) FileGroups {
	if c == nil || len(c.Groups) < 1 {
		return nil
	}

	name := file.Filename
	if file.GithubCommitFile != nil && file.GithubCommitFile.GetFilename() != "" {
		name = file.GithubCommitFile.GetFilename()
	}
	if name == "" {
		return nil
	}

	res := FileGroups{}
	for _, group := range c.Groups {
		if group.Ignore.Match(name) {
			continue
		}

		if group.Match.Match(name) {
			res = res.Add(FileGroup(group.Name))
		}
	}

	return res
}

// Match returns true if any of the matchers in the collection match the given path.
func (m Matchers) Match(path string) bool {
	if len(m) < 1 {
		return false
	}

	for _, matcher := range m {
		if matcher.Match(path) {
			return true
		}
	}

	return false
}

// Match returns true if the path matches all configured criteria in the MatchConfig.
func (m *MatchConfig) Match(path string) bool {
	if m == nil {
		return false
	}
	matched := false

	if len(m.File) > 0 {
		if !slices.Contains(m.File, path) {
			return false
		}
		matched = true
	}

	if len(m.BaseDir) > 0 {
		if !slices.ContainsFunc(m.BaseDir, func(dir string) bool {
			return hasBaseDir(path, dir)
		}) {
			return false
		}
		matched = true
	}

	if len(m.BaseName) > 0 {
		if !slices.ContainsFunc(m.BaseName, func(name string) bool {
			return filepath.Base(path) == name
		}) {
			return false
		}
		matched = true
	}

	if len(m.BaseNamePrefix) > 0 {
		if !slices.ContainsFunc(m.BaseNamePrefix, func(prefix string) bool {
			return strings.Contains(filepath.Base(path), prefix)
		}) {
			return false
		}
		matched = true
	}

	if len(m.Contains) > 0 {
		if !slices.ContainsFunc(m.Contains, func(contains string) bool {
			return strings.Contains(path, contains)
		}) {
			return false
		}
		matched = true
	}

	if len(m.Extension) > 0 {
		if !slices.ContainsFunc(m.Extension, func(ext string) bool {
			return filepath.Ext(path) == ext
		}) {
			return false
		}
		matched = true
	}

	return matched
}

func hasBaseDir(name, dir string) bool {
	return strings.HasPrefix(name, dir+string(os.PathSeparator))
}
