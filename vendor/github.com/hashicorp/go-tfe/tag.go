// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import "fmt"

type TagList struct {
	*Pagination
	Items []*Tag
}

// Tag is owned by an organization and applied to workspaces. Used for grouping and search.
type Tag struct {
	ID   string `jsonapi:"primary,tags"`
	Name string `jsonapi:"attr,name,omitempty"`
}

type TagBinding struct {
	ID    string `jsonapi:"primary,tag-bindings"`
	Key   string `jsonapi:"attr,key"`
	Value string `jsonapi:"attr,value,omitempty"`
}

func encodeTagFiltersAsParams(filters []*TagBinding) map[string][]string {
	if len(filters) == 0 {
		return nil
	}

	var tagFilter = make(map[string][]string, len(filters))
	for index, tag := range filters {
		tagFilter[fmt.Sprintf("filter[tagged][%d][key]", index)] = []string{tag.Key}
		tagFilter[fmt.Sprintf("filter[tagged][%d][value]", index)] = []string{tag.Value}
	}

	return tagFilter
}
