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
func (c *File) Name() string {
	if c == nil || c.File == nil {
		return ""
	}

	return c.File.GetFilename()
}

// Add takes a variadic set of groups and adds them to the ordered set of groups
func (c FileGroups) Add(groups ...FileGroup) FileGroups {
	for _, group := range groups {
		idx, in := c.In(group)
		if in {
			continue
		}

		c = slices.Insert(c, idx, group)
	}

	return c
}

// In takes a group and determines the index and presence of the group in the group set
func (c FileGroups) In(group FileGroup) (int, bool) {
	return slices.BinarySearch(c, group)
}

// String is a string representation of all groups a file is in
func (c FileGroups) String() string {
	groups := []string{}
	for _, g := range c {
		groups = append(groups, string(g))
	}

	return strings.Join(groups, ", ")
}
