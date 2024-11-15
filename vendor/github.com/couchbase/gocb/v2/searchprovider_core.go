package gocb

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/couchbase/gocb/v2/vector"
	"time"

	cbsearch "github.com/couchbase/gocb/v2/search"
	"github.com/couchbase/gocbcore/v10"
)

var defaultVectorQueryNumCandidates = uint32(3)

type searchProviderWrapper struct {
	agent *gocbcore.AgentGroup
}

func (search *searchProviderWrapper) SearchQuery(ctx context.Context, opts gocbcore.SearchQueryOptions) (searchRowReader, error) {

	opm := newAsyncOpManager(ctx)
	var errOut error
	var sOut *gocbcore.SearchRowReader

	err := opm.Wait(search.agent.SearchQuery(opts, func(reader *gocbcore.SearchRowReader, err error) {
		if err != nil {
			errOut = err
			opm.Reject()
			return
		}

		sOut = reader
		opm.Resolve()
	}))
	if err != nil {
		errOut = err
	}

	if errOut != nil {
		err = errOut
	}
	return sOut, err
}

// used to allow mocking for testing
type searchProviderCoreProvider interface {
	SearchQuery(ctx context.Context, opts gocbcore.SearchQueryOptions) (searchRowReader, error)
}

type searchCapabilityVerifier interface {
	SearchCapabilityStatus(cap gocbcore.SearchCapability) gocbcore.CapabilityStatus
}

type searchProviderCore struct {
	// agent *gocbcore.AgentGroup
	provider searchProviderCoreProvider

	retryStrategyWrapper *coreRetryStrategyWrapper
	transcoder           Transcoder
	timeouts             TimeoutsConfig
	tracer               RequestTracer
	meter                *meterWrapper
}

func (search *searchProviderCore) Search(scope *Scope, indexName string, request SearchRequest, opts *SearchOptions) (*SearchResult, error) {
	searchQuery := request.SearchQuery
	if searchQuery == nil {
		// See MB-60312.
		searchQuery = cbsearch.NewMatchNoneQuery()
	}
	return search.search(scope, indexName, searchQuery, request.VectorSearch, false, opts)
}

func (search *searchProviderCore) SearchQuery(indexName string, query cbsearch.Query, opts *SearchOptions) (*SearchResult, error) {
	return search.search(nil, indexName, query, nil, true, opts)
}

func (search *searchProviderCore) search(scope *Scope, indexName string, sQuery cbsearch.Query, vSearch *vector.Search, showRequest bool, opts *SearchOptions) (*SearchResult, error) {
	if sQuery == nil && vSearch == nil {
		return nil, makeInvalidArgumentsError("must specify either a search query or a vector search")
	}

	start := time.Now()
	defer search.meter.ValueRecord(meterValueServiceSearch, "search", start)

	span := createSpan(search.tracer, opts.ParentSpan, "search", "search")
	span.SetAttribute("db.operation", indexName)
	if scope != nil {
		span.SetAttribute("db.name", scope.BucketName())
		span.SetAttribute("db.couchbase.scope", scope.Name())
	}
	defer span.End()

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = search.timeouts.SearchTimeout
	}
	deadline := time.Now().Add(timeout)

	retryStrategy := search.retryStrategyWrapper
	if opts.RetryStrategy != nil {
		retryStrategy = newCoreRetryStrategyWrapper(opts.RetryStrategy)
	}

	searchOpts, err := opts.toMap(indexName)
	if err != nil {
		return nil, &SearchError{
			InnerError: wrapError(err, "failed to generate query options"),
		}
	}
	if !showRequest {
		searchOpts["showrequest"] = false
	}
	if sQuery != nil {
		searchOpts["query"] = sQuery
	}
	if vSearch != nil {
		internalVSearch := vSearch.Internal()

		if err := internalVSearch.Validate(); err != nil {
			return nil, makeInvalidArgumentsError(err.Error())
		}

		queries := make([]vector.InternalQuery, len(internalVSearch.Queries))
		for i, query := range internalVSearch.Queries {
			if query.NumCandidates == nil {
				query.NumCandidates = &defaultVectorQueryNumCandidates
			}
			queries[i] = query
		}

		searchOpts["knn"] = queries
		if internalVSearch.VectorQueryCombination != vector.VectorQueryCombinationNotSet {
			searchOpts["knn_operator"] = string(internalVSearch.VectorQueryCombination)
		}
	}

	return search.execSearchQuery(opts.Context, span, scope, indexName, searchOpts, deadline, retryStrategy, opts.Internal.User)

}

func (search *searchProviderCore) execSearchQuery(
	ctx context.Context,
	span RequestSpan,
	scope *Scope,
	indexName string,
	options map[string]interface{},
	deadline time.Time,
	retryStrategy *coreRetryStrategyWrapper,
	user string,
) (*SearchResult, error) {

	eSpan := createSpan(search.tracer, span, "request_encoding", "")
	reqBytes, err := json.Marshal(options)
	eSpan.End()
	if err != nil {
		return nil, &SearchError{
			InnerError: wrapError(err, "failed to marshall query body"),
			Query:      maybeGetSearchOptionQuery(options),
		}
	}

	coreOpts := gocbcore.SearchQueryOptions{
		IndexName:     indexName,
		Payload:       reqBytes,
		RetryStrategy: retryStrategy,
		Deadline:      deadline,
		TraceContext:  span.Context(),
		User:          user,
	}

	if scope != nil {
		coreOpts.BucketName = scope.bucket.bucketName
		coreOpts.ScopeName = scope.scopeName
	}

	res, err := search.provider.SearchQuery(ctx, coreOpts)
	if err != nil {
		return nil, maybeEnhanceSearchError(err)
	}

	return newSearchResult(res), nil
}

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

type jsonSearchResponseStatus struct {
	Errors     map[string]string `json:"errors"`
	Failed     uint64            `json:"failed"`
	Successful uint64            `json:"successful"`
}

type jsonSearchResponse struct {
	Status    jsonSearchResponseStatus   `json:"status,omitempty"`
	TotalHits uint64                     `json:"total_hits"`
	MaxScore  float64                    `json:"max_score"`
	Took      uint64                     `json:"took"`
	Facets    map[string]jsonSearchFacet `json:"facets"`
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
	err := srr.reader.Err()
	if err != nil {
		return maybeEnhanceSearchError(err)
	}

	return nil
}

// Close marks the results as closed, returning any errors that occurred during reading the results.
func (srr *SearchResultRaw) Close() error {
	err := srr.reader.Close()
	if err != nil {
		return maybeEnhanceSearchError(err)
	}

	return nil
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
		return maybeEnhanceSearchError(err)
	}
	// This is an error from json unmarshal so no point in trying to enhance it.
	return r.jsonErr
}

// Close marks the results as closed, returning any errors that occurred during reading the results.
func (r *SearchResult) Close() error {
	if r.reader == nil {
		return r.Err()
	}

	err := r.reader.Close()
	if err != nil {
		return maybeEnhanceSearchError(err)
	}

	return nil
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
