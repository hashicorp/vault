package gocb

import (
	"encoding/json"
	"errors"
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

type jsonSearchTermFacet struct {
	Term  string `json:"term,omitempty"`
	Count int    `json:"count,omitempty"`
}

type jsonSearchNumericFacet struct {
	Name  string  `json:"name,omitempty"`
	Min   float64 `json:"min,omitempty"`
	Max   float64 `json:"max,omitempty"`
	Count int     `json:"count,omitempty"`
}

type jsonSearchDateFacet struct {
	Name  string `json:"name,omitempty"`
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
	Count int    `json:"count,omitempty"`
}

type jsonSearchFacet struct {
	Name          string                   `json:"name"`
	Field         string                   `json:"field"`
	Total         uint64                   `json:"total"`
	Missing       uint64                   `json:"missing"`
	Other         uint64                   `json:"other"`
	Terms         []jsonSearchTermFacet    `json:"terms"`
	NumericRanges []jsonSearchNumericFacet `json:"numeric_ranges"`
	DateRanges    []jsonSearchDateFacet    `json:"date_ranges"`
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

type searchRowReader interface {
	NextRow() []byte
	Err() error
	MetaData() ([]byte, error)
	Close() error
}

// SearchResultRaw provides raw access to search data.
// VOLATILE: This API is subject to change at any time.
type SearchResultRaw struct {
	reader searchRowReader
}

// NextBytes returns the next row as bytes.
func (srr *SearchResultRaw) NextBytes() []byte {
	return srr.reader.NextRow()
}

// Err returns any errors that have occurred on the stream
func (srr *SearchResultRaw) Err() error {
	return srr.reader.Err()
}

// Close marks the results as closed, returning any errors that occurred during reading the results.
func (srr *SearchResultRaw) Close() error {
	return srr.reader.Close()
}

// MetaData returns any meta-data that was available from this query as bytes.
func (srr *SearchResultRaw) MetaData() ([]byte, error) {
	return srr.reader.MetaData()
}

// SearchResult allows access to the results of a search query.
type SearchResult struct {
	reader searchRowReader

	currentRow SearchRow
	jsonErr    error
}

func newSearchResult(reader searchRowReader) *SearchResult {
	return &SearchResult{
		reader: reader,
	}
}

// Raw returns a SearchResultRaw which can be used to access the raw byte data from search queries.
// Calling this function invalidates the underlying SearchResult which will no longer be able to be used.
// VOLATILE: This API is subject to change at any time.
func (r *SearchResult) Raw() *SearchResultRaw {
	vr := &SearchResultRaw{
		reader: r.reader,
	}

	r.reader = nil
	return vr
}

// Next assigns the next result from the results into the value pointer, returning whether the read was successful.
func (r *SearchResult) Next() bool {
	if r.reader == nil {
		return false
	}

	rowBytes := r.reader.NextRow()
	if rowBytes == nil {
		return false
	}

	r.currentRow = SearchRow{}

	var rowData jsonSearchRow
	if err := json.Unmarshal(rowBytes, &rowData); err != nil {
		// This should never happen but if it does then lets store it in a best efforts basis and maybe the next
		// row will be ok. We can then return this from .Err().
		r.jsonErr = err
		return true
	}

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

	return true
}

// Row returns the contents of the current row.
func (r *SearchResult) Row() SearchRow {
	if r.reader == nil {
		return SearchRow{}
	}

	return r.currentRow
}

// Err returns any errors that have occurred on the stream
func (r *SearchResult) Err() error {
	if r.reader == nil {
		return errors.New("result object is no longer valid")
	}

	err := r.reader.Err()
	if err != nil {
		return err
	}
	return r.jsonErr
}

// Close marks the results as closed, returning any errors that occurred during reading the results.
func (r *SearchResult) Close() error {
	if r.reader == nil {
		return r.Err()
	}

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
	if r.reader == nil {
		return nil, r.Err()
	}

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
