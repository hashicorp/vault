package gocb

import (
	"encoding/json"
	"time"

	cbsearch "github.com/couchbase/gocb/v2/search"
	gocbcore "github.com/couchbase/gocbcore/v9"
)

type jsonRowLocation struct {
	Field          string   `json:"field"`
	Term           string   `json:"term"`
	Position       uint32   `json:"position"`
	Start          uint32   `json:"start"`
	End            uint32   `json:"end"`
	ArrayPositions []uint32 `json:"array_positions"`
}

type jsonSearchFacet struct {
	Name    string `json:"name"`
	Field   string `json:"field"`
	Total   uint64 `json:"total"`
	Missing uint64 `json:"missing"`
	Other   uint64 `json:"other"`
}

type jsonSearchRowLocations map[string]map[string][]jsonRowLocation

type jsonSearchRow struct {
	Index       string                 `json:"index"`
	ID          string                 `json:"id"`
	Score       float64                `json:"score"`
	Explanation interface{}            `json:"explanation"`
	Locations   jsonSearchRowLocations `json:"locations"`
	Fragments   map[string][]string    `json:"fragments"`
	Fields      json.RawMessage        `json:"fields"`
}

type jsonSearchResponse struct {
	Errors    map[string]string          `json:"errors"`
	TotalHits uint64                     `json:"total_hits"`
	MaxScore  float64                    `json:"max_score"`
	Took      uint64                     `json:"took"`
	Facets    map[string]jsonSearchFacet `json:"facets"`
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
	metrics.Took = time.Duration(data.Took) * time.Microsecond

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
	meta.Errors = data.Errors

	return nil
}

// SearchFacetResult provides access to the result of a faceted query.
type SearchFacetResult struct {
	Name    string
	Field   string
	Total   uint64
	Missing uint64
	Other   uint64
}

func (fr *SearchFacetResult) fromData(data jsonSearchFacet) error {
	fr.Name = data.Name
	fr.Field = data.Field
	fr.Total = data.Total
	fr.Missing = data.Missing
	fr.Other = data.Other

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

type searchRowReader interface {
	NextRow() []byte
	Err() error
	MetaData() ([]byte, error)
	Close() error
}

// SearchResult allows access to the results of a search query.
type SearchResult struct {
	reader searchRowReader

	currentRow SearchRow
}

func newSearchResult(reader searchRowReader) *SearchResult {
	return &SearchResult{
		reader: reader,
	}
}

// Next assigns the next result from the results into the value pointer, returning whether the read was successful.
func (r *SearchResult) Next() bool {
	rowBytes := r.reader.NextRow()
	if rowBytes == nil {
		return false
	}

	r.currentRow = SearchRow{}

	var rowData jsonSearchRow
	if err := json.Unmarshal(rowBytes, &rowData); err == nil {
		r.currentRow.Index = rowData.Index
		r.currentRow.ID = rowData.ID
		r.currentRow.Score = rowData.Score
		r.currentRow.Explanation = rowData.Explanation
		r.currentRow.Fragments = rowData.Fragments
		r.currentRow.fieldsBytes = rowData.Fields

		locations := make(map[string]map[string][]SearchRowLocation)
		for fieldName, fieldData := range rowData.Locations {
			terms := make(map[string][]SearchRowLocation)
			for termName, termData := range fieldData {
				locations := make([]SearchRowLocation, len(termData))
				for locIdx, locData := range termData {
					err := locations[locIdx].fromData(locData)
					if err != nil {
						logWarnf("failed to parse search query location data: %s", err)
					}
				}
				terms[termName] = locations
			}
			locations[fieldName] = terms
		}
		r.currentRow.Locations = locations
	}

	return true
}

// Row returns the contents of the current row.
func (r *SearchResult) Row() SearchRow {
	return r.currentRow
}

// Err returns any errors that have occurred on the stream
func (r *SearchResult) Err() error {
	return r.reader.Err()
}

// Close marks the results as closed, returning any errors that occurred during reading the results.
func (r *SearchResult) Close() error {
	return r.reader.Close()
}

func (r *SearchResult) getJSONResp() (jsonSearchResponse, error) {
	metaDataBytes, err := r.reader.MetaData()
	if err != nil {
		return jsonSearchResponse{}, err
	}

	var jsonResp jsonSearchResponse
	err = json.Unmarshal(metaDataBytes, &jsonResp)
	if err != nil {
		return jsonSearchResponse{}, err
	}

	return jsonResp, nil
}

// MetaData returns any meta-data that was available from this query.  Note that
// the meta-data will only be available once the object has been closed (either
// implicitly or explicitly).
func (r *SearchResult) MetaData() (*SearchMetaData, error) {
	jsonResp, err := r.getJSONResp()
	if err != nil {
		return nil, err
	}

	var metaData SearchMetaData
	err = metaData.fromData(jsonResp)
	if err != nil {
		return nil, err
	}

	return &metaData, nil
}

// Facets returns any facets that were returned with this query.  Note that the
// facets will only be available once the object has been closed (either
// implicitly or explicitly).
func (r *SearchResult) Facets() (map[string]SearchFacetResult, error) {
	jsonResp, err := r.getJSONResp()
	if err != nil {
		return nil, err
	}

	facets := make(map[string]SearchFacetResult)
	for facetName, facetData := range jsonResp.Facets {
		var facet SearchFacetResult
		err := facet.fromData(facetData)
		if err != nil {
			return nil, err
		}

		facets[facetName] = facet
	}

	return facets, nil
}

// SearchQuery executes the analytics query statement on the server.
func (c *Cluster) SearchQuery(indexName string, query cbsearch.Query, opts *SearchOptions) (*SearchResult, error) {
	if opts == nil {
		opts = &SearchOptions{}
	}

	span := c.tracer.StartSpan("SearchQuery", opts.parentSpan).
		SetTag("couchbase.service", "search")
	defer span.Finish()

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = c.timeoutsConfig.SearchTimeout
	}
	deadline := time.Now().Add(timeout)

	retryStrategy := c.retryStrategyWrapper
	if opts.RetryStrategy != nil {
		retryStrategy = newRetryStrategyWrapper(opts.RetryStrategy)
	}

	searchOpts, err := opts.toMap()
	if err != nil {
		return nil, SearchError{
			InnerError: wrapError(err, "failed to generate query options"),
			Query:      query,
		}
	}

	searchOpts["query"] = query

	return c.execSearchQuery(span, indexName, searchOpts, deadline, retryStrategy)
}

func maybeGetSearchOptionQuery(options map[string]interface{}) interface{} {
	if value, ok := options["query"]; ok {
		return value
	}
	return ""
}

func (c *Cluster) execSearchQuery(
	span requestSpan,
	indexName string,
	options map[string]interface{},
	deadline time.Time,
	retryStrategy *retryStrategyWrapper,
) (*SearchResult, error) {
	provider, err := c.getSearchProvider()
	if err != nil {
		return nil, SearchError{
			InnerError: wrapError(err, "failed to get query provider"),
			Query:      maybeGetSearchOptionQuery(options),
		}
	}

	reqBytes, err := json.Marshal(options)
	if err != nil {
		return nil, SearchError{
			InnerError: wrapError(err, "failed to marshall query body"),
			Query:      maybeGetSearchOptionQuery(options),
		}
	}

	res, err := provider.SearchQuery(gocbcore.SearchQueryOptions{
		IndexName:     indexName,
		Payload:       reqBytes,
		RetryStrategy: retryStrategy,
		Deadline:      deadline,
		TraceContext:  span.Context(),
	})
	if err != nil {
		return nil, maybeEnhanceSearchError(err)
	}

	return newSearchResult(res), nil
}
