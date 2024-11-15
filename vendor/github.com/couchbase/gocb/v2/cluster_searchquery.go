package gocb

import (
	"encoding/json"
	"time"

	cbsearch "github.com/couchbase/gocb/v2/search"
)

// SearchQuery executes the search query on the server.
func (c *Cluster) SearchQuery(indexName string, query cbsearch.Query, opts *SearchOptions) (*SearchResult, error) {
	return autoOpControl(c.searchController(), func(provider searchProvider) (*SearchResult, error) {
		if opts == nil {
			opts = &SearchOptions{}
		}

		return provider.SearchQuery(indexName, query, opts)
	})
}

// Search executes the search request on the server.
func (c *Cluster) Search(indexName string, request SearchRequest, opts *SearchOptions) (*SearchResult, error) {
	return autoOpControl(c.searchController(), func(provider searchProvider) (*SearchResult, error) {
		if request.VectorSearch == nil && request.SearchQuery == nil {
			return nil, makeInvalidArgumentsError("the search request cannot be empty")
		}

		if opts == nil {
			opts = &SearchOptions{}
		}

		return provider.Search(nil, indexName, request, opts)
	})
}

func maybeGetSearchOptionQuery(options map[string]interface{}) interface{} {
	if value, ok := options["query"]; ok {
		return value
	}
	return ""
}

// SearchMetrics encapsulates various metrics gathered during a search queries execution.
type SearchMetrics struct {
	Took                  time.Duration
	TotalRows             uint64
	MaxScore              float64
	TotalPartitionCount   uint64
	SuccessPartitionCount uint64
	ErrorPartitionCount   uint64
}

func (metrics *SearchMetrics) fromData(data jsonSearchResponse) error {
	metrics.TotalRows = data.TotalHits
	metrics.MaxScore = data.MaxScore
	metrics.Took = time.Duration(data.Took) / time.Nanosecond
	metrics.TotalPartitionCount = data.Status.Successful + data.Status.Failed
	metrics.SuccessPartitionCount = data.Status.Successful
	metrics.ErrorPartitionCount = data.Status.Failed

	return nil
}

// SearchMetaData provides access to the meta-data properties of a search query result.
type SearchMetaData struct {
	Metrics SearchMetrics
	Errors  map[string]string
}

func (meta *SearchMetaData) fromData(data jsonSearchResponse) error {
	metrics := SearchMetrics{}
	if err := metrics.fromData(data); err != nil {
		return err
	}

	meta.Metrics = metrics
	meta.Errors = data.Status.Errors

	return nil
}

// SearchTermFacetResult holds the results of a term facet in search results.
type SearchTermFacetResult struct {
	Term  string
	Count int
}

// SearchNumericRangeFacetResult holds the results of a numeric facet in search results.
type SearchNumericRangeFacetResult struct {
	Name  string
	Min   float64
	Max   float64
	Count int
}

// SearchDateRangeFacetResult holds the results of a date facet in search results.
type SearchDateRangeFacetResult struct {
	Name  string
	Start string
	End   string
	Count int
}

// SearchFacetResult provides access to the result of a faceted query.
type SearchFacetResult struct {
	Name          string
	Field         string
	Total         uint64
	Missing       uint64
	Other         uint64
	Terms         []SearchTermFacetResult
	NumericRanges []SearchNumericRangeFacetResult
	DateRanges    []SearchDateRangeFacetResult
}

func (fr *SearchFacetResult) fromData(data jsonSearchFacet) error {
	fr.Name = data.Name
	fr.Field = data.Field
	fr.Total = data.Total
	fr.Missing = data.Missing
	fr.Other = data.Other
	for _, term := range data.Terms {
		fr.Terms = append(fr.Terms, SearchTermFacetResult(term))
	}
	for _, nr := range data.NumericRanges {
		fr.NumericRanges = append(fr.NumericRanges, SearchNumericRangeFacetResult(nr))
	}
	for _, nr := range data.DateRanges {
		fr.DateRanges = append(fr.DateRanges, SearchDateRangeFacetResult(nr))
	}

	return nil
}

// SearchRowLocation represents the location of a row match
type SearchRowLocation struct {
	Position       uint32
	Start          uint32
	End            uint32
	ArrayPositions []uint32
}

func (rl *SearchRowLocation) fromData(data jsonRowLocation) error {
	rl.Position = data.Position
	rl.Start = data.Start
	rl.End = data.End
	rl.ArrayPositions = data.ArrayPositions

	return nil
}

// SearchRow represents a single hit returned from a search query.
type SearchRow struct {
	Index       string
	ID          string
	Score       float64
	Explanation interface{}
	Locations   map[string]map[string][]SearchRowLocation
	Fragments   map[string][]string
	fieldsBytes []byte
}

// Fields decodes the fields included in a search hit.
func (sr *SearchRow) Fields(valuePtr interface{}) error {
	return json.Unmarshal(sr.fieldsBytes, valuePtr)
}
