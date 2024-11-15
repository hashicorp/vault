// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"github.com/hashicorp/nomad/api/contexts"
)

type Search struct {
	client *Client
}

// Search returns a handle on the Search endpoints
func (c *Client) Search() *Search {
	return &Search{client: c}
}

// PrefixSearch returns a set of matches for a particular context and prefix.
func (s *Search) PrefixSearch(prefix string, context contexts.Context, q *QueryOptions) (*SearchResponse, *QueryMeta, error) {
	var resp SearchResponse
	req := &SearchRequest{Prefix: prefix, Context: context}

	qm, err := s.client.putQuery("/v1/search", req, &resp, q)
	if err != nil {
		return nil, nil, err
	}

	return &resp, qm, nil
}

type SearchResponse struct {
	Matches     map[contexts.Context][]string
	Truncations map[contexts.Context]bool
	QueryMeta
}

type SearchRequest struct {
	Prefix  string
	Context contexts.Context
	QueryOptions
}

// FuzzySearch returns a set of matches for a given context and string.
func (s *Search) FuzzySearch(text string, context contexts.Context, q *QueryOptions) (*FuzzySearchResponse, *QueryMeta, error) {
	var resp FuzzySearchResponse

	req := &FuzzySearchRequest{
		Context: context,
		Text:    text,
	}

	qm, err := s.client.putQuery("/v1/search/fuzzy", req, &resp, q)
	if err != nil {
		return nil, nil, err
	}

	return &resp, qm, nil
}

// FuzzyMatch is used to describe the ID of an object which may be a machine
// readable UUID or a human readable Name. If the object is a component of a Job,
// the Scope is a list of IDs starting from Namespace down to the parent object of
// ID.
//
// e.g. A Task-level service would have scope like,
//
//	["<namespace>", "<job>", "<group>", "<task>"]
type FuzzyMatch struct {
	ID    string   // ID is UUID or Name of object
	Scope []string `json:",omitempty"` // IDs of parent objects
}

// FuzzySearchResponse is used to return fuzzy matches and information about
// whether the match list is truncated specific to each type of searchable Context.
type FuzzySearchResponse struct {
	// Matches is a map of Context types to IDs which fuzzy match a specified query.
	Matches map[contexts.Context][]FuzzyMatch

	// Truncations indicates whether the matches for a particular Context have
	// been truncated.
	Truncations map[contexts.Context]bool

	QueryMeta
}

// FuzzySearchRequest is used to parameterize a fuzzy search request, and returns
// a list of matches made up of jobs, allocations, evaluations, and/or nodes,
// along with whether or not the information returned is truncated.
type FuzzySearchRequest struct {
	// Text is what names are fuzzy-matched to. E.g. if the given text were
	// "py", potential matches might be "python", "mypy", etc. of jobs, nodes,
	// allocs, groups, services, commands, images, classes.
	Text string

	// Context is the type that can be matched against. A Context of "all" indicates
	// all Contexts types are queried for matching.
	Context contexts.Context

	QueryOptions
}
