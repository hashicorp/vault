package gocb

import (
	"github.com/couchbase/gocb/v2/search"
	"github.com/couchbase/gocb/v2/vector"
)

// SearchRequest is used for describing a search request used with Search.
type SearchRequest struct {
	SearchQuery  search.Query
	VectorSearch *vector.Search
}
