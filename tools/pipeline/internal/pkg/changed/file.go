// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package changed

import (
	"slices"
	"strings"

	gh "github.com/google/go-github/v68/github"
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
	// FileGroup is a set of groups a changed file belongs to. Use FileGroups.Add() instead of append()
	// to ensure uniqueness and ordering
	FileGroups []FileGroup
)

const (
	FileGroupAutopilot  FileGroup = "autopilot"
	FileGroupChangelog  FileGroup = "changelog"
	FileGroupCommunity  FileGroup = "community"
	FileGroupDocs       FileGroup = "docs"
	FileGroupEnos       FileGroup = "enos"
	FileGroupEnterprise FileGroup = "enterprise"
	FileGroupGoApp      FileGroup = "app"
	FileGroupGoModules  FileGroup = "gomod"
	FileGroupPipeline   FileGroup = "pipeline"
	FileGroupProto      FileGroup = "proto"
	FileGroupTools      FileGroup = "tools"
	FileGroupWebUI      FileGroup = "ui"
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

// String is a string representation of all groups a file is in
func (g FileGroups) String() string {
	groups := []string{}
	for _, g := range g {
		groups = append(groups, string(g))
	}

	return strings.Join(groups, ", ")
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
