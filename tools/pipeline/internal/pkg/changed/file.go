// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

// Package changed is a package for inspecting and categorizing changed files into groups
package changed

import (
	"slices"
	"strings"

	gh "github.com/google/go-github/v81/github"
)

type (
	// File is a changed file in a PR or commit
	File struct {
		File   *gh.CommitFile `json:"file,omitempty"`
		Groups FileGroups     `json:"groups,omitempty"`
	}
	// Files is a slice of changed files in a PR or commit
	Files []*File
	// FileGroup is group name describing a class of file the changed file belongs to
	FileGroup string
	// FileGroups is a set of groups a changed file belongs to. Use FileGroups.Add() instead of append()
	// to ensure uniqueness and ordering
	FileGroups []FileGroup
)

const (
	FileGroupAutopilot   FileGroup = "autopilot"
	FileGroupChangelog   FileGroup = "changelog"
	FileGroupCommunity   FileGroup = "community"
	FileGroupDocs        FileGroup = "docs"
	FileGroupEnos        FileGroup = "enos"
	FileGroupEnterprise  FileGroup = "enterprise"
	FileGroupGithub      FileGroup = "github"
	FileGroupGoApp       FileGroup = "app"
	FileGroupGoToolchain FileGroup = "gotoolchain"
	FileGroupPipeline    FileGroup = "pipeline"
	FileGroupProto       FileGroup = "proto"
	FileGroupTools       FileGroup = "tools"
	FileGroupWebUI       FileGroup = "ui"
)

// Name is the file name of the changed file
func (f *File) Name() string {
	if f == nil || f.File == nil {
		return ""
	}

	return f.File.GetFilename()
}

// Add takes a variadic set of groups and adds them to the ordered set of groups
func (g FileGroups) Add(groups ...FileGroup) FileGroups {
	for _, group := range groups {
		idx, in := g.In(group)
		if in {
			continue
		}

		g = slices.Insert(g, idx, group)
	}

	return g
}

// In takes a group and determines the index and presence of the group in the group set
func (g FileGroups) In(group FileGroup) (int, bool) {
	return slices.BinarySearch(g, group)
}

// All takes another FileGroups and determines whether or not all of the groups in the
// in group are included in FileGroups.
func (g FileGroups) All(groups FileGroups) bool {
	for _, group := range groups {
		if _, in := g.In(group); !in {
			return false
		}
	}

	return true
}

// Any takes another FileGroups and determines whether or not any of the groups in the
// in group are included in FileGroups.
func (g FileGroups) Any(groups FileGroups) bool {
	for _, group := range groups {
		if _, in := g.In(group); in {
			return true
		}
	}

	return false
}

// Groups returns the FileGroups as a slice of strings
func (g FileGroups) Groups() []string {
	groups := []string{}
	for _, g := range g {
		groups = append(groups, string(g))
	}

	return groups
}

// String is a string representation of all groups a file is in
func (g FileGroups) String() string {
	return strings.Join(g.Groups(), ", ")
}

// Names returns a list of file names
func (f Files) Names() []string {
	if len(f) < 1 {
		return nil
	}
	files := []string{}
	for _, file := range f {
		files = append(files, file.Name())
	}

	return files
}

// EachHasAnyGroup determines whether each file contains the any of the given groups
func (f Files) EachHasAnyGroup(groups FileGroups) bool {
	if f == nil {
		return false
	}

	if len(groups) == 0 {
		return true
	}

	for _, file := range f {
		if file.Groups == nil {
			return false
		}

		if !file.Groups.Any(groups) {
			return false
		}
	}

	return true
}
