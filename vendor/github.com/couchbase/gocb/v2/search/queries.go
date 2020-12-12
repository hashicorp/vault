package search

import "encoding/json"

// Query represents a search query.
type Query interface {
}

type searchQueryBase struct {
	options map[string]interface{}
}

func newSearchQueryBase() searchQueryBase {
	return searchQueryBase{
		options: make(map[string]interface{}),
	}
}

// MarshalJSON marshal's this query to JSON for the search REST API.
func (q searchQueryBase) MarshalJSON() ([]byte, error) {
	return json.Marshal(q.options)
}

// MatchQuery represents a search match query.
type MatchQuery struct {
	searchQueryBase
}

// NewMatchQuery creates a new MatchQuery.
func NewMatchQuery(match string) *MatchQuery {
	q := &MatchQuery{newSearchQueryBase()}
	q.options["match"] = match
	return q
}

// Field specifies the field for this query.
func (q *MatchQuery) Field(field string) *MatchQuery {
	q.options["field"] = field
	return q
}

// Analyzer specifies the analyzer to use for this query.
func (q *MatchQuery) Analyzer(analyzer string) *MatchQuery {
	q.options["analyzer"] = analyzer
	return q
}

// PrefixLength specifies the prefix length from this query.
func (q *MatchQuery) PrefixLength(length uint64) *MatchQuery {
	q.options["prefix_length"] = length
	return q
}

// Fuzziness specifies the fuziness for this query.
func (q *MatchQuery) Fuzziness(fuzziness uint64) *MatchQuery {
	q.options["fuzziness"] = fuzziness
	return q
}

// Boost specifies the boost for this query.
func (q *MatchQuery) Boost(boost float32) *MatchQuery {
	q.options["boost"] = boost
	return q
}

// MatchPhraseQuery represents a search match phrase query.
type MatchPhraseQuery struct {
	searchQueryBase
}

// NewMatchPhraseQuery creates a new MatchPhraseQuery
func NewMatchPhraseQuery(phrase string) *MatchPhraseQuery {
	q := &MatchPhraseQuery{newSearchQueryBase()}
	q.options["match_phrase"] = phrase
	return q
}

// Field specifies the field for this query.
func (q *MatchPhraseQuery) Field(field string) *MatchPhraseQuery {
	q.options["field"] = field
	return q
}

// Analyzer specifies the analyzer to use for this query.
func (q *MatchPhraseQuery) Analyzer(analyzer string) *MatchPhraseQuery {
	q.options["analyzer"] = analyzer
	return q
}

// Boost specifies the boost for this query.
func (q *MatchPhraseQuery) Boost(boost float32) *MatchPhraseQuery {
	q.options["boost"] = boost
	return q
}

// RegexpQuery represents a search regular expression query.
type RegexpQuery struct {
	searchQueryBase
}

// NewRegexpQuery creates a new RegexpQuery.
func NewRegexpQuery(regexp string) *RegexpQuery {
	q := &RegexpQuery{newSearchQueryBase()}
	q.options["regexp"] = regexp
	return q
}

// Field specifies the field for this query.
func (q *RegexpQuery) Field(field string) *RegexpQuery {
	q.options["field"] = field
	return q
}

// Boost specifies the boost for this query.
func (q *RegexpQuery) Boost(boost float32) *RegexpQuery {
	q.options["boost"] = boost
	return q
}

// QueryStringQuery represents a search string query.
type QueryStringQuery struct {
	searchQueryBase
}

// NewQueryStringQuery creates a new StringQuery.
func NewQueryStringQuery(query string) *QueryStringQuery {
	q := &QueryStringQuery{newSearchQueryBase()}
	q.options["query"] = query
	return q
}

// Boost specifies the boost for this query.
func (q *QueryStringQuery) Boost(boost float32) *QueryStringQuery {
	q.options["boost"] = boost
	return q
}

// NumericRangeQuery represents a search numeric range query.
type NumericRangeQuery struct {
	searchQueryBase
}

// NewNumericRangeQuery creates a new NumericRangeQuery.
func NewNumericRangeQuery() *NumericRangeQuery {
	q := &NumericRangeQuery{newSearchQueryBase()}
	return q
}

// Min specifies the minimum value and inclusiveness for this range query.
func (q *NumericRangeQuery) Min(min float32, inclusive bool) *NumericRangeQuery {
	q.options["min"] = min
	q.options["inclusive_min"] = inclusive
	return q
}

// Max specifies the maximum value and inclusiveness for this range query.
func (q *NumericRangeQuery) Max(max float32, inclusive bool) *NumericRangeQuery {
	q.options["max"] = max
	q.options["inclusive_max"] = inclusive
	return q
}

// Field specifies the field for this query.
func (q *NumericRangeQuery) Field(field string) *NumericRangeQuery {
	q.options["field"] = field
	return q
}

// Boost specifies the boost for this query.
func (q *NumericRangeQuery) Boost(boost float32) *NumericRangeQuery {
	q.options["boost"] = boost
	return q
}

// DateRangeQuery represents a search date range query.
type DateRangeQuery struct {
	searchQueryBase
}

// NewDateRangeQuery creates a new DateRangeQuery.
func NewDateRangeQuery() *DateRangeQuery {
	q := &DateRangeQuery{newSearchQueryBase()}
	return q
}

// Start specifies the start value and inclusiveness for this range query.
func (q *DateRangeQuery) Start(start string, inclusive bool) *DateRangeQuery {
	q.options["start"] = start
	q.options["inclusive_start"] = inclusive
	return q
}

// End specifies the end value and inclusiveness for this range query.
func (q *DateRangeQuery) End(end string, inclusive bool) *DateRangeQuery {
	q.options["end"] = end
	q.options["inclusive_end"] = inclusive
	return q
}

// DateTimeParser specifies which date time string parser to use.
func (q *DateRangeQuery) DateTimeParser(parser string) *DateRangeQuery {
	q.options["datetime_parser"] = parser
	return q
}

// Field specifies the field for this query.
func (q *DateRangeQuery) Field(field string) *DateRangeQuery {
	q.options["field"] = field
	return q
}

// Boost specifies the boost for this query.
func (q *DateRangeQuery) Boost(boost float32) *DateRangeQuery {
	q.options["boost"] = boost
	return q
}

// ConjunctionQuery represents a search conjunction query.
type ConjunctionQuery struct {
	searchQueryBase
}

// NewConjunctionQuery creates a new ConjunctionQuery.
func NewConjunctionQuery(queries ...Query) *ConjunctionQuery {
	q := &ConjunctionQuery{newSearchQueryBase()}
	q.options["conjuncts"] = []Query{}
	return q.And(queries...)
}

// And adds new predicate queries to this conjunction query.
func (q *ConjunctionQuery) And(queries ...Query) *ConjunctionQuery {
	q.options["conjuncts"] = append(q.options["conjuncts"].([]Query), queries...)
	return q
}

// Boost specifies the boost for this query.
func (q *ConjunctionQuery) Boost(boost float32) *ConjunctionQuery {
	q.options["boost"] = boost
	return q
}

// DisjunctionQuery represents a search disjunction query.
type DisjunctionQuery struct {
	searchQueryBase
}

// NewDisjunctionQuery creates a new DisjunctionQuery.
func NewDisjunctionQuery(queries ...Query) *DisjunctionQuery {
	q := &DisjunctionQuery{newSearchQueryBase()}
	q.options["disjuncts"] = []Query{}
	return q.Or(queries...)
}

// Or adds new predicate queries to this disjunction query.
func (q *DisjunctionQuery) Or(queries ...Query) *DisjunctionQuery {
	q.options["disjuncts"] = append(q.options["disjuncts"].([]Query), queries...)
	return q
}

// Boost specifies the boost for this query.
func (q *DisjunctionQuery) Boost(boost float32) *DisjunctionQuery {
	q.options["boost"] = boost
	return q
}

type booleanQueryData struct {
	Must    *ConjunctionQuery `json:"must,omitempty"`
	Should  *DisjunctionQuery `json:"should,omitempty"`
	MustNot *DisjunctionQuery `json:"must_not,omitempty"`
	Boost   float32           `json:"boost,omitempty"`
}

// BooleanQuery represents a search boolean query.
type BooleanQuery struct {
	data      booleanQueryData
	shouldMin uint32
}

// NewBooleanQuery creates a new BooleanQuery.
func NewBooleanQuery() *BooleanQuery {
	q := &BooleanQuery{}
	return q
}

// Must specifies a query which must match.
func (q *BooleanQuery) Must(query Query) *BooleanQuery {
	switch val := query.(type) {
	case ConjunctionQuery:
		q.data.Must = &val
	case *ConjunctionQuery:
		q.data.Must = val
	default:
		q.data.Must = NewConjunctionQuery(val)
	}
	return q
}

// Should specifies a query which should match.
func (q *BooleanQuery) Should(query Query) *BooleanQuery {
	switch val := query.(type) {
	case DisjunctionQuery:
		q.data.Should = &val
	case *DisjunctionQuery:
		q.data.Should = val
	default:
		q.data.Should = NewDisjunctionQuery(val)
	}
	return q
}

// MustNot specifies a query which must not match.
func (q *BooleanQuery) MustNot(query Query) *BooleanQuery {
	switch val := query.(type) {
	case DisjunctionQuery:
		q.data.MustNot = &val
	case *DisjunctionQuery:
		q.data.MustNot = val
	default:
		q.data.MustNot = NewDisjunctionQuery(val)
	}
	return q
}

// ShouldMin specifies the minimum value before the should query will boost.
func (q *BooleanQuery) ShouldMin(min uint32) *BooleanQuery {
	q.shouldMin = min
	return q
}

// Boost specifies the boost for this query.
func (q *BooleanQuery) Boost(boost float32) *BooleanQuery {
	q.data.Boost = boost
	return q
}

// MarshalJSON marshal's this query to JSON for the search REST API.
func (q *BooleanQuery) MarshalJSON() ([]byte, error) {
	if q.data.Should != nil {
		q.data.Should.options["min"] = q.shouldMin
	}
	bytes, err := json.Marshal(q.data)
	if q.data.Should != nil {
		delete(q.data.Should.options, "min")
	}
	return bytes, err
}

// WildcardQuery represents a search wildcard query.
type WildcardQuery struct {
	searchQueryBase
}

// NewWildcardQuery creates a new WildcardQuery.
func NewWildcardQuery(wildcard string) *WildcardQuery {
	q := &WildcardQuery{newSearchQueryBase()}
	q.options["wildcard"] = wildcard
	return q
}

// Field specifies the field for this query.
func (q *WildcardQuery) Field(field string) *WildcardQuery {
	q.options["field"] = field
	return q
}

// Boost specifies the boost for this query.
func (q *WildcardQuery) Boost(boost float32) *WildcardQuery {
	q.options["boost"] = boost
	return q
}

// DocIDQuery represents a search document id query.
type DocIDQuery struct {
	searchQueryBase
}

// NewDocIDQuery creates a new DocIdQuery.
func NewDocIDQuery(ids ...string) *DocIDQuery {
	q := &DocIDQuery{newSearchQueryBase()}
	q.options["ids"] = []string{}
	return q.AddDocIds(ids...)
}

// AddDocIds adds addition document ids to this query.
func (q *DocIDQuery) AddDocIds(ids ...string) *DocIDQuery {
	q.options["ids"] = append(q.options["ids"].([]string), ids...)
	return q
}

// Field specifies the field for this query.
func (q *DocIDQuery) Field(field string) *DocIDQuery {
	q.options["field"] = field
	return q
}

// Boost specifies the boost for this query.
func (q *DocIDQuery) Boost(boost float32) *DocIDQuery {
	q.options["boost"] = boost
	return q
}

// BooleanFieldQuery represents a search boolean field query.
type BooleanFieldQuery struct {
	searchQueryBase
}

// NewBooleanFieldQuery creates a new BooleanFieldQuery.
func NewBooleanFieldQuery(val bool) *BooleanFieldQuery {
	q := &BooleanFieldQuery{newSearchQueryBase()}
	q.options["bool"] = val
	return q
}

// Field specifies the field for this query.
func (q *BooleanFieldQuery) Field(field string) *BooleanFieldQuery {
	q.options["field"] = field
	return q
}

// Boost specifies the boost for this query.
func (q *BooleanFieldQuery) Boost(boost float32) *BooleanFieldQuery {
	q.options["boost"] = boost
	return q
}

// TermQuery represents a search term query.
type TermQuery struct {
	searchQueryBase
}

// NewTermQuery creates a new TermQuery.
func NewTermQuery(term string) *TermQuery {
	q := &TermQuery{newSearchQueryBase()}
	q.options["term"] = term
	return q
}

// Field specifies the field for this query.
func (q *TermQuery) Field(field string) *TermQuery {
	q.options["field"] = field
	return q
}

// PrefixLength specifies the prefix length from this query.
func (q *TermQuery) PrefixLength(length uint64) *TermQuery {
	q.options["prefix_length"] = length
	return q
}

// Fuzziness specifies the fuziness for this query.
func (q *TermQuery) Fuzziness(fuzziness uint64) *TermQuery {
	q.options["fuzziness"] = fuzziness
	return q
}

// Boost specifies the boost for this query.
func (q *TermQuery) Boost(boost float32) *TermQuery {
	q.options["boost"] = boost
	return q
}

// PhraseQuery represents a search phrase query.
type PhraseQuery struct {
	searchQueryBase
}

// NewPhraseQuery creates a new PhraseQuery.
func NewPhraseQuery(terms ...string) *PhraseQuery {
	q := &PhraseQuery{newSearchQueryBase()}
	q.options["terms"] = terms
	return q
}

// Field specifies the field for this query.
func (q *PhraseQuery) Field(field string) *PhraseQuery {
	q.options["field"] = field
	return q
}

// Boost specifies the boost for this query.
func (q *PhraseQuery) Boost(boost float32) *PhraseQuery {
	q.options["boost"] = boost
	return q
}

// PrefixQuery represents a search prefix query.
type PrefixQuery struct {
	searchQueryBase
}

// NewPrefixQuery creates a new PrefixQuery.
func NewPrefixQuery(prefix string) *PrefixQuery {
	q := &PrefixQuery{newSearchQueryBase()}
	q.options["prefix"] = prefix
	return q
}

// Field specifies the field for this query.
func (q *PrefixQuery) Field(field string) *PrefixQuery {
	q.options["field"] = field
	return q
}

// Boost specifies the boost for this query.
func (q *PrefixQuery) Boost(boost float32) *PrefixQuery {
	q.options["boost"] = boost
	return q
}

// MatchAllQuery represents a search match all query.
type MatchAllQuery struct {
	searchQueryBase
}

// NewMatchAllQuery creates a new MatchAllQuery.
func NewMatchAllQuery() *MatchAllQuery {
	q := &MatchAllQuery{newSearchQueryBase()}
	q.options["match_all"] = nil
	return q
}

// MatchNoneQuery represents a search match none query.
type MatchNoneQuery struct {
	searchQueryBase
}

// NewMatchNoneQuery creates a new MatchNoneQuery.
func NewMatchNoneQuery() *MatchNoneQuery {
	q := &MatchNoneQuery{newSearchQueryBase()}
	q.options["match_none"] = nil
	return q
}

// TermRangeQuery represents a search term range query.
type TermRangeQuery struct {
	searchQueryBase
}

// NewTermRangeQuery creates a new TermRangeQuery.
func NewTermRangeQuery(term string) *TermRangeQuery {
	q := &TermRangeQuery{newSearchQueryBase()}
	q.options["term"] = term
	return q
}

// Field specifies the field for this query.
func (q *TermRangeQuery) Field(field string) *TermRangeQuery {
	q.options["field"] = field
	return q
}

// Min specifies the minimum value and inclusiveness for this range query.
func (q *TermRangeQuery) Min(min string, inclusive bool) *TermRangeQuery {
	q.options["min"] = min
	q.options["inclusive_min"] = inclusive
	return q
}

// Max specifies the maximum value and inclusiveness for this range query.
func (q *TermRangeQuery) Max(max string, inclusive bool) *TermRangeQuery {
	q.options["max"] = max
	q.options["inclusive_max"] = inclusive
	return q
}

// Boost specifies the boost for this query.
func (q *TermRangeQuery) Boost(boost float32) *TermRangeQuery {
	q.options["boost"] = boost
	return q
}

// GeoDistanceQuery represents a search geographical distance query.
type GeoDistanceQuery struct {
	searchQueryBase
}

// NewGeoDistanceQuery creates a new GeoDistanceQuery.
func NewGeoDistanceQuery(lon, lat float64, distance string) *GeoDistanceQuery {
	q := &GeoDistanceQuery{newSearchQueryBase()}
	q.options["location"] = []float64{lon, lat}
	q.options["distance"] = distance
	return q
}

// Field specifies the field for this query.
func (q *GeoDistanceQuery) Field(field string) *GeoDistanceQuery {
	q.options["field"] = field
	return q
}

// Boost specifies the boost for this query.
func (q *GeoDistanceQuery) Boost(boost float32) *GeoDistanceQuery {
	q.options["boost"] = boost
	return q
}

// GeoBoundingBoxQuery represents a search geographical bounding box query.
type GeoBoundingBoxQuery struct {
	searchQueryBase
}

// NewGeoBoundingBoxQuery creates a new GeoBoundingBoxQuery.
func NewGeoBoundingBoxQuery(tlLon, tlLat, brLon, brLat float64) *GeoBoundingBoxQuery {
	q := &GeoBoundingBoxQuery{newSearchQueryBase()}
	q.options["top_left"] = []float64{tlLon, tlLat}
	q.options["bottom_right"] = []float64{brLon, brLat}
	return q
}

// Field specifies the field for this query.
func (q *GeoBoundingBoxQuery) Field(field string) *GeoBoundingBoxQuery {
	q.options["field"] = field
	return q
}

// Boost specifies the boost for this query.
func (q *GeoBoundingBoxQuery) Boost(boost float32) *GeoBoundingBoxQuery {
	q.options["boost"] = boost
	return q
}
