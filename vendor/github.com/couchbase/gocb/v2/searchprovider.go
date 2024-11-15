package gocb

import cbsearch "github.com/couchbase/gocb/v2/search"

type searchProvider interface {
	SearchQuery(indexName string, query cbsearch.Query, opts *SearchOptions) (*SearchResult, error)
	Search(scope *Scope, indexName string, request SearchRequest, opts *SearchOptions) (*SearchResult, error)
}
