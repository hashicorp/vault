package search

import (
	"encoding/json"
)

// Query represents a search query.
type Query interface {
}

// MatchOperator defines how the individual match terms should be logically concatenated.
type MatchOperator string

const (
	// MatchOperatorOr specifies that individual match terms are concatenated with a logical OR - this is the default if not provided.
	MatchOperatorOr MatchOperator = "or"

	// MatchOperatorAnd specifies that individual match terms are concatenated with a logical AND.
	MatchOperatorAnd MatchOperator = "and"
)

// MatchQuery represents a search match query.
type MatchQuery struct {
	match        string
	field        *string
	analyzer     *string
	prefixLength *uint64
	fuzziness    *uint64
	boost        *float32
	operator     *MatchOperator
}

// marshal's query to JSON for use with search REST API.
func (q MatchQuery) MarshalJSON() ([]byte, error) {
	outStruct := &struct {
		Match        string   `json:"match"`
		Field        *string  `json:"field,omitempty"`
		Analyzer     *string  `json:"analyzer,omitempty"`
		PrefixLength *uint64  `json:"prefix_length,omitempty"`
		Fuzziness    *uint64  `json:"fuzziness,omitempty"`
		Boost        *float32 `json:"boost,omitempty"`
		Operator     *string  `json:"operator,omitempty"`
	}{
		Match:        q.match,
		Field:        q.field,
		Analyzer:     q.analyzer,
		PrefixLength: q.prefixLength,
		Fuzziness:    q.fuzziness,
		Boost:        q.boost,
	}
	if q.operator != nil {
		operator := string(*q.operator)
		outStruct.Operator = &operator
	}

	return json.Marshal(outStruct)
}

// NewMatchQuery creates a new MatchQuery.
func NewMatchQuery(match string) *MatchQuery {
	q := &MatchQuery{match: match}
	return q
}

// Field specifies the field for this query.
func (q *MatchQuery) Field(field string) *MatchQuery {
	q.field = &field
	return q
}

// Analyzer specifies the analyzer to use for this query.
func (q *MatchQuery) Analyzer(analyzer string) *MatchQuery {
	q.analyzer = &analyzer
	return q
}

// PrefixLength specifies the prefix length from this query.
func (q *MatchQuery) PrefixLength(length uint64) *MatchQuery {
	q.prefixLength = &length
	return q
}

// Fuzziness specifies the fuziness for this query.
func (q *MatchQuery) Fuzziness(fuzziness uint64) *MatchQuery {
	q.fuzziness = &fuzziness
	return q
}

// Boost specifies the boost for this query.
func (q *MatchQuery) Boost(boost float32) *MatchQuery {
	q.boost = &boost
	return q
}

// Operator defines how the individual match terms should be logically concatenated.
func (q *MatchQuery) Operator(operator MatchOperator) *MatchQuery {
	// q.options["operator"] = string(operator)
	q.operator = &operator
	return q
}

// MatchPhraseQuery represents a search match phrase query.
type MatchPhraseQuery struct {
	matchPhrase string
	field       *string
	analyzer    *string
	boost       *float32
}

// marshal's query to JSON for use with search REST API.
func (q MatchPhraseQuery) MarshalJSON() ([]byte, error) {
	outStruct := &struct {
		MatchPhrase string   `json:"match_phrase"`
		Field       *string  `json:"field,omitempty"`
		Analyzer    *string  `json:"analyzer,omitempty"`
		Boost       *float32 `json:"boost,omitempty"`
	}{
		MatchPhrase: q.matchPhrase,
		Field:       q.field,
		Analyzer:    q.analyzer,
		Boost:       q.boost,
	}

	return json.Marshal(outStruct)
}

// NewMatchPhraseQuery creates a new MatchPhraseQuery
func NewMatchPhraseQuery(phrase string) *MatchPhraseQuery {
	q := &MatchPhraseQuery{matchPhrase: phrase}
	return q
}

// Field specifies the field for this query.
func (q *MatchPhraseQuery) Field(field string) *MatchPhraseQuery {
	q.field = &field
	return q
}

// Analyzer specifies the analyzer to use for this query.
func (q *MatchPhraseQuery) Analyzer(analyzer string) *MatchPhraseQuery {
	q.analyzer = &analyzer
	return q
}

// Boost specifies the boost for this query.
func (q *MatchPhraseQuery) Boost(boost float32) *MatchPhraseQuery {
	q.boost = &boost
	return q
}

// RegexpQuery represents a search regular expression query.
type RegexpQuery struct {
	regexp string
	field  *string
	boost  *float32
}

// marshal's query to JSON for use with search REST API.
func (q RegexpQuery) MarshalJSON() ([]byte, error) {
	outStruct := &struct {
		Regexp string   `json:"regexp"`
		Field  *string  `json:"field,omitempty"`
		Boost  *float32 `json:"boost"`
	}{
		Regexp: q.regexp,
		Field:  q.field,
		Boost:  q.boost,
	}

	return json.Marshal(outStruct)
}

// NewRegexpQuery creates a new RegexpQuery.
func NewRegexpQuery(regexp string) *RegexpQuery {
	q := &RegexpQuery{regexp: regexp}
	return q
}

// Field specifies the field for this query.
func (q *RegexpQuery) Field(field string) *RegexpQuery {
	q.field = &field
	return q
}

// Boost specifies the boost for this query.
func (q *RegexpQuery) Boost(boost float32) *RegexpQuery {
	q.boost = &boost
	return q
}

// QueryStringQuery represents a search string query.
type QueryStringQuery struct {
	query string
	boost *float32
}

// marshal's query to JSON for use with search REST API.
func (q QueryStringQuery) MarshalJSON() ([]byte, error) {
	outStruct := &struct {
		Query string   `json:"query"`
		Boost *float32 `json:"boost"`
	}{
		Query: q.query,
		Boost: q.boost,
	}

	return json.Marshal(outStruct)

}

// NewQueryStringQuery creates a new StringQuery.
func NewQueryStringQuery(query string) *QueryStringQuery {
	q := &QueryStringQuery{query: query}
	return q
}

// Boost specifies the boost for this query.
func (q *QueryStringQuery) Boost(boost float32) *QueryStringQuery {
	q.boost = &boost
	return q
}

// NumericRangeQuery represents a search numeric range query.
type NumericRangeQuery struct {
	min          *float32
	inclusiveMin *bool
	max          *float32
	inclusiveMax *bool
	field        *string
	boost        *float32
}

// marshal's query to JSON for use with search REST API.
func (q NumericRangeQuery) MarshalJSON() ([]byte, error) {
	outStruct := &struct {
		Min          *float32 `json:"min,omitempty"`
		InclusiveMin *bool    `json:"inclusive_min,omitempty"`
		Max          *float32 `json:"max,omitempty"`
		InclusiveMax *bool    `json:"inclusive_max,omitempty"`
		Field        *string  `json:"field,omitempty"`
		Boost        *float32 `json:"boost,omitempty"`
	}{
		Min:          q.min,
		InclusiveMin: q.inclusiveMin,
		Max:          q.max,
		InclusiveMax: q.inclusiveMax,
		Field:        q.field,
		Boost:        q.boost,
	}

	return json.Marshal(outStruct)
}

// NewNumericRangeQuery creates a new NumericRangeQuery.
func NewNumericRangeQuery() *NumericRangeQuery {
	q := &NumericRangeQuery{}
	return q
}

// Min specifies the minimum value and inclusiveness for this range query.
func (q *NumericRangeQuery) Min(min float32, inclusive bool) *NumericRangeQuery {
	q.min = &min
	q.inclusiveMin = &inclusive
	return q
}

// Max specifies the maximum value and inclusiveness for this range query.
func (q *NumericRangeQuery) Max(max float32, inclusive bool) *NumericRangeQuery {
	q.max = &max
	q.inclusiveMax = &inclusive
	return q
}

// Field specifies the field for this query.
func (q *NumericRangeQuery) Field(field string) *NumericRangeQuery {
	q.field = &field
	return q
}

// Boost specifies the boost for this query.
func (q *NumericRangeQuery) Boost(boost float32) *NumericRangeQuery {
	q.boost = &boost
	return q
}

// DateRangeQuery represents a search date range query.
type DateRangeQuery struct {
	start          *string
	inclusiveStart *bool
	end            *string
	inclusiveEnd   *bool
	dateTimeParser *string
	field          *string
	boost          *float32
}

// marshal's query to JSON for use with search REST API.
func (q DateRangeQuery) MarshalJSON() ([]byte, error) {
	outStruct := &struct {
		Start          *string  `json:"start,omitempty"`
		InclusiveStart *bool    `json:"inclusive_start,omitempty"`
		End            *string  `json:"end,omitempty"`
		InclusiveEnd   *bool    `json:"inclusive_end,omitempty"`
		DateTimeParser *string  `json:"datetime_parser,omitempty"`
		Field          *string  `json:"field,omitempty"`
		Boost          *float32 `json:"boost,omitempty"`
	}{
		Start:          q.start,
		InclusiveStart: q.inclusiveStart,
		End:            q.end,
		InclusiveEnd:   q.inclusiveEnd,
		DateTimeParser: q.dateTimeParser,
		Field:          q.field,
		Boost:          q.boost,
	}

	return json.Marshal(outStruct)
}

// NewDateRangeQuery creates a new DateRangeQuery.
func NewDateRangeQuery() *DateRangeQuery {
	q := &DateRangeQuery{}
	return q
}

// Start specifies the start value and inclusiveness for this range query.
func (q *DateRangeQuery) Start(start string, inclusive bool) *DateRangeQuery {
	q.start = &start
	q.inclusiveStart = &inclusive
	return q
}

// End specifies the end value and inclusiveness for this range query.
func (q *DateRangeQuery) End(end string, inclusive bool) *DateRangeQuery {
	q.end = &end
	q.inclusiveEnd = &inclusive
	return q
}

// DateTimeParser specifies which date time string parser to use.
func (q *DateRangeQuery) DateTimeParser(parser string) *DateRangeQuery {
	q.dateTimeParser = &parser
	return q
}

// Field specifies the field for this query.
func (q *DateRangeQuery) Field(field string) *DateRangeQuery {
	q.field = &field
	return q
}

// Boost specifies the boost for this query.
func (q *DateRangeQuery) Boost(boost float32) *DateRangeQuery {
	q.boost = &boost
	return q
}

// ConjunctionQuery represents a search conjunction query.
type ConjunctionQuery struct {
	conjuncts []Query
	boost     *float32
}

// marshal's query to JSON for use with search REST API.
func (q ConjunctionQuery) MarshalJSON() ([]byte, error) {
	outStruct := &struct {
		Conjuncts []Query  `json:"conjuncts"`
		Boost     *float32 `json:"boost,omitempty"`
	}{
		Conjuncts: q.conjuncts,
		Boost:     q.boost,
	}

	return json.Marshal(outStruct)
}

// NewConjunctionQuery creates a new ConjunctionQuery.
func NewConjunctionQuery(queries ...Query) *ConjunctionQuery {
	q := &ConjunctionQuery{
		conjuncts: make([]Query, 0),
	}
	return q.And(queries...)
}

// And adds new predicate queries to this conjunction query.
func (q *ConjunctionQuery) And(queries ...Query) *ConjunctionQuery {
	q.conjuncts = append(q.conjuncts, queries...)
	return q
}

// Boost specifies the boost for this query.
func (q *ConjunctionQuery) Boost(boost float32) *ConjunctionQuery {
	q.boost = &boost
	return q
}

// DisjunctionQuery represents a search disjunction query.
type DisjunctionQuery struct {
	disjuncts []Query
	boost     *float32
	min       *uint32
}

// marshal's query to JSON for use with search REST API.
func (q DisjunctionQuery) MarshalJSON() ([]byte, error) {
	outStruct := &struct {
		Disjuncts []Query  `json:"disjuncts"`
		Boost     *float32 `json:"boost,omitempty"`
		Min       *uint32  `json:"min,omitempty"`
	}{
		Disjuncts: q.disjuncts,
		Boost:     q.boost,
		Min:       q.min,
	}

	return json.Marshal(outStruct)
}

// NewDisjunctionQuery creates a new DisjunctionQuery.
func NewDisjunctionQuery(queries ...Query) *DisjunctionQuery {
	q := &DisjunctionQuery{
		disjuncts: make([]Query, 0),
	}

	return q.Or(queries...)
}

// Or adds new predicate queries to this disjunction query.
func (q *DisjunctionQuery) Or(queries ...Query) *DisjunctionQuery {
	q.disjuncts = append(q.disjuncts, queries...)
	return q
}

// Boost specifies the boost for this query.
func (q *DisjunctionQuery) Boost(boost float32) *DisjunctionQuery {
	q.boost = &boost
	return q
}

// Min specifies the minimum number of queries that a document must satisfy.
func (q *DisjunctionQuery) Min(min uint32) *DisjunctionQuery {
	q.min = &min
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
	shouldMin uint32 // minimum value before the should query will boost.
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
func (q BooleanQuery) MarshalJSON() ([]byte, error) {
	if q.data.Should != nil {
		q.data.Should.min = &q.shouldMin
	}
	bytes, err := json.Marshal(q.data)
	if q.data.Should != nil {
		q.data.Should.min = nil
	}
	return bytes, err
}

// WildcardQuery represents a search wildcard query.
type WildcardQuery struct {
	wildcard string
	field    *string
	boost    *float32
}

// MarshalJSON marshal's this query to JSON for the search REST API.
func (q WildcardQuery) MarshalJSON() ([]byte, error) {
	outStruct := &struct {
		Wildcard string   `json:"wildcard"`
		Field    *string  `json:"field,omitempty"`
		Boost    *float32 `json:"boost,omitempty"`
	}{
		Wildcard: q.wildcard,
		Field:    q.field,
		Boost:    q.boost,
	}

	return json.Marshal(outStruct)
}

// NewWildcardQuery creates a new WildcardQuery.
func NewWildcardQuery(wildcard string) *WildcardQuery {
	q := &WildcardQuery{wildcard: wildcard}
	return q
}

// Field specifies the field for this query.
func (q *WildcardQuery) Field(field string) *WildcardQuery {
	q.field = &field
	return q
}

// Boost specifies the boost for this query.
func (q *WildcardQuery) Boost(boost float32) *WildcardQuery {
	q.boost = &boost
	return q
}

// DocIDQuery represents a search document id query.
type DocIDQuery struct {
	ids   []string
	field *string
	boost *float32
}

// MarshalJSON marshal's this query to JSON for the search REST API.
func (q DocIDQuery) MarshalJSON() ([]byte, error) {
	outStruct := &struct {
		Ids   []string `json:"ids"`
		Field *string  `json:"field,omitempty"`
		Boost *float32 `json:"boost,omitempty"`
	}{
		Ids:   q.ids,
		Field: q.field,
		Boost: q.boost,
	}

	return json.Marshal(outStruct)
}

// NewDocIDQuery creates a new DocIdQuery.
func NewDocIDQuery(ids ...string) *DocIDQuery {
	q := &DocIDQuery{ids: make([]string, 0)}

	return q.AddDocIds(ids...)
}

// AddDocIds adds addition document ids to this query.
func (q *DocIDQuery) AddDocIds(ids ...string) *DocIDQuery {
	q.ids = append(q.ids, ids...)
	return q
}

// Field specifies the field for this query.
func (q *DocIDQuery) Field(field string) *DocIDQuery {
	q.field = &field
	return q
}

// Boost specifies the boost for this query.
func (q *DocIDQuery) Boost(boost float32) *DocIDQuery {
	q.boost = &boost
	return q
}

// BooleanFieldQuery represents a search boolean field query.
type BooleanFieldQuery struct {
	val   bool     // required field
	field *string  // optional field
	boost *float32 // optional field
}

// MarshalJSON marshal's this query to JSON for the search REST API.
func (q BooleanFieldQuery) MarshalJSON() ([]byte, error) {
	outStruct := &struct {
		Val   bool     `json:"bool"`
		Field *string  `json:"field,omitempty"`
		Boost *float32 `json:"boost,omitempty"`
	}{
		Val:   q.val,
		Field: q.field,
		Boost: q.boost,
	}

	return json.Marshal(outStruct)
}

// NewBooleanFieldQuery creates a new BooleanFieldQuery.
func NewBooleanFieldQuery(val bool) *BooleanFieldQuery {
	q := &BooleanFieldQuery{val: val}
	return q
}

// Field specifies the field for this query.
func (q *BooleanFieldQuery) Field(field string) *BooleanFieldQuery {
	q.field = &field
	return q
}

// Boost specifies the boost for this query.
func (q *BooleanFieldQuery) Boost(boost float32) *BooleanFieldQuery {
	q.boost = &boost
	return q
}

// TermQuery represents a search term query.
type TermQuery struct {
	term         string // required field
	field        *string
	boost        *float32
	prefixLength *uint64
	fuzziness    *uint64
}

// MarshalJSON marshal's this query to JSON for the search REST API.
func (q TermQuery) MarshalJSON() ([]byte, error) {
	outStruct := &struct {
		Term         string   `json:"term"`
		Field        *string  `json:"field,omitempty"`
		Boost        *float32 `json:"boost,omitempty"`
		PrefixLength *uint64  `json:"prefix_length,omitempty"`
		Fuzziness    *uint64  `json:"fuzziness,omitempty"`
	}{
		Term:         q.term,
		Field:        q.field,
		Boost:        q.boost,
		PrefixLength: q.prefixLength,
		Fuzziness:    q.fuzziness,
	}

	return json.Marshal(outStruct)
}

// NewTermQuery creates a new TermQuery.
func NewTermQuery(term string) *TermQuery {
	q := &TermQuery{term: term}
	return q
}

// Field specifies the field for this query.
func (q *TermQuery) Field(field string) *TermQuery {
	q.field = &field
	return q
}

// PrefixLength specifies the prefix length from this query.
func (q *TermQuery) PrefixLength(length uint64) *TermQuery {
	q.prefixLength = &length
	return q
}

// Fuzziness specifies the fuziness for this query.
func (q *TermQuery) Fuzziness(fuzziness uint64) *TermQuery {
	q.fuzziness = &fuzziness
	return q
}

// Boost specifies the boost for this query.
func (q *TermQuery) Boost(boost float32) *TermQuery {
	q.boost = &boost
	return q
}

// PhraseQuery represents a search phrase query.
type PhraseQuery struct {
	terms []string
	field *string
	boost *float32
}

// MarshalJSON marshal's this query to JSON for the search REST API.
func (q PhraseQuery) MarshalJSON() ([]byte, error) {
	outStruct := &struct {
		Terms []string `json:"terms"`
		Field *string  `json:"field,omitempty"`
		Boost *float32 `json:"boost,omitempty"`
	}{
		Terms: q.terms,
		Field: q.field,
		Boost: q.boost,
	}

	return json.Marshal(outStruct)
}

// NewPhraseQuery creates a new PhraseQuery.
func NewPhraseQuery(terms ...string) *PhraseQuery {
	q := &PhraseQuery{terms: terms}
	return q
}

// Field specifies the field for this query.
func (q *PhraseQuery) Field(field string) *PhraseQuery {
	q.field = &field
	return q
}

// Boost specifies the boost for this query.
func (q *PhraseQuery) Boost(boost float32) *PhraseQuery {
	q.boost = &boost
	return q
}

// PrefixQuery represents a search prefix query.
type PrefixQuery struct {
	prefix string
	field  *string
	boost  *float32
}

// MarshalJSON marshal's this query to JSON for the search REST API.
func (q PrefixQuery) MarshalJSON() ([]byte, error) {
	outStruct := &struct {
		Prefix string   `json:"prefix"`
		Field  *string  `json:"field,omitempty"`
		Boost  *float32 `json:"boost,omitempty"`
	}{
		Prefix: q.prefix,
		Field:  q.field,
		Boost:  q.boost,
	}

	return json.Marshal(outStruct)
}

// NewPrefixQuery creates a new PrefixQuery.
func NewPrefixQuery(prefix string) *PrefixQuery {
	q := &PrefixQuery{prefix: prefix}
	return q
}

// Field specifies the field for this query.
func (q *PrefixQuery) Field(field string) *PrefixQuery {
	q.field = &field
	return q
}

// Boost specifies the boost for this query.
func (q *PrefixQuery) Boost(boost float32) *PrefixQuery {
	q.boost = &boost
	return q
}

// MatchAllQuery represents a search match all query.
type MatchAllQuery struct{}

// MarshalJSON marshal's this query to JSON for the search REST API.
func (q MatchAllQuery) MarshalJSON() ([]byte, error) {
	outStruct := &struct {
		MatchAll *struct{} `json:"match_all"`
	}{}

	return json.Marshal(outStruct)
}

// NewMatchAllQuery creates a new MatchAllQuery.
func NewMatchAllQuery() *MatchAllQuery {
	q := &MatchAllQuery{}
	return q
}

// MatchNoneQuery represents a search match none query.
type MatchNoneQuery struct{}

// MarshalJSON marshal's this query to JSON for the search REST API.
func (q MatchNoneQuery) MarshalJSON() ([]byte, error) {
	outStruct := &struct {
		MatchAll *struct{} `json:"match_none"`
	}{}

	return json.Marshal(outStruct)
}

// NewMatchNoneQuery creates a new MatchNoneQuery.
func NewMatchNoneQuery() *MatchNoneQuery {
	q := &MatchNoneQuery{}
	return q
}

// TermRangeQuery represents a search term range query.
type TermRangeQuery struct {
	term         string
	field        *string
	min          *string
	inclusiveMin *bool
	max          *string
	inclusiveMax *bool
	boost        *float32
}

// MarshalJSON marshal's this query to JSON for the search REST API.
func (q TermRangeQuery) MarshalJSON() ([]byte, error) {
	outStruct := &struct {
		Term         string   `json:"term"`
		Field        *string  `json:"field,omitempty"`
		Min          *string  `json:"min,omitempty"`
		InclusiveMin *bool    `json:"inclusive_min,omitempty"`
		Max          *string  `json:"max,omitempty"`
		InclusiveMax *bool    `json:"inclusive_max,omitempty"`
		Boost        *float32 `json:"boost,omitempty"`
	}{
		Term:         q.term,
		Field:        q.field,
		Min:          q.min,
		InclusiveMin: q.inclusiveMin,
		Max:          q.max,
		InclusiveMax: q.inclusiveMax,
		Boost:        q.boost,
	}

	return json.Marshal(outStruct)
}

// NewTermRangeQuery creates a new TermRangeQuery.
func NewTermRangeQuery(term string) *TermRangeQuery {
	q := &TermRangeQuery{term: term}
	return q
}

// Field specifies the field for this query.
func (q *TermRangeQuery) Field(field string) *TermRangeQuery {
	q.field = &field
	return q
}

// Min specifies the minimum value and inclusiveness for this range query.
func (q *TermRangeQuery) Min(min string, inclusive bool) *TermRangeQuery {
	q.min = &min
	q.inclusiveMin = &inclusive
	return q
}

// Max specifies the maximum value and inclusiveness for this range query.
func (q *TermRangeQuery) Max(max string, inclusive bool) *TermRangeQuery {
	q.max = &max
	q.inclusiveMax = &inclusive
	return q
}

// Boost specifies the boost for this query.
func (q *TermRangeQuery) Boost(boost float32) *TermRangeQuery {
	q.boost = &boost
	return q
}

// GeoDistanceQuery represents a search geographical distance query.
type GeoDistanceQuery struct {
	location []float64
	distance string
	field    *string
	boost    *float32
}

// MarshalJSON marshal's this query to JSON for the search REST API.
func (q GeoDistanceQuery) MarshalJSON() ([]byte, error) {
	outStruct := &struct {
		Location []float64 `json:"location"`
		Distance string    `json:"distance"`
		Field    *string   `json:"field,omitempty"`
		Boost    *float32  `json:"boost,omitempty"`
	}{
		Location: q.location,
		Distance: q.distance,
		Field:    q.field,
		Boost:    q.boost,
	}

	return json.Marshal(outStruct)
}

// NewGeoDistanceQuery creates a new GeoDistanceQuery.
func NewGeoDistanceQuery(lon, lat float64, distance string) *GeoDistanceQuery {
	q := &GeoDistanceQuery{location: []float64{lon, lat}, distance: distance}
	return q
}

// Field specifies the field for this query.
func (q *GeoDistanceQuery) Field(field string) *GeoDistanceQuery {
	q.field = &field
	return q
}

// Boost specifies the boost for this query.
func (q *GeoDistanceQuery) Boost(boost float32) *GeoDistanceQuery {
	q.boost = &boost
	return q
}

// GeoBoundingBoxQuery represents a search geographical bounding box query.
type GeoBoundingBoxQuery struct {
	topLeft     []float64
	bottomRight []float64
	field       *string
	boost       *float32
}

// MarshalJSON marshal's this query to JSON for the search REST API.
func (q GeoBoundingBoxQuery) MarshalJSON() ([]byte, error) {
	outStruct := &struct {
		TopLeft     []float64 `json:"top_left"`
		BottomRight []float64 `json:"bottom_right"`
		Field       *string   `json:"field,omitempty"`
		Boost       *float32  `json:"boost,omitempty"`
	}{
		TopLeft:     q.topLeft,
		BottomRight: q.bottomRight,
		Field:       q.field,
		Boost:       q.boost,
	}

	return json.Marshal(outStruct)
}

// NewGeoBoundingBoxQuery creates a new GeoBoundingBoxQuery.
func NewGeoBoundingBoxQuery(tlLon, tlLat, brLon, brLat float64) *GeoBoundingBoxQuery {
	q := &GeoBoundingBoxQuery{topLeft: []float64{tlLon, tlLat}, bottomRight: []float64{brLon, brLat}}
	return q
}

// Field specifies the field for this query.
func (q *GeoBoundingBoxQuery) Field(field string) *GeoBoundingBoxQuery {
	q.field = &field
	return q
}

// Boost specifies the boost for this query.
func (q *GeoBoundingBoxQuery) Boost(boost float32) *GeoBoundingBoxQuery {
	q.boost = &boost
	return q
}

// Coordinate is a tuple of a latitude and a longitude.
type Coordinate struct {
	Lon float64
	Lat float64
}

// GeoPolygonQuery represents a search query which allows to match inside a geo polygon.
type GeoPolygonQuery struct {
	polyPoints [][]float64
	field      *string
	boost      *float32
}

// MarshalJSON marshal's this query to JSON for the search REST API.
func (q GeoPolygonQuery) MarshalJSON() ([]byte, error) {
	outStruct := &struct {
		PolyPoints [][]float64 `json:"polygon_points"`
		Field      *string     `json:"field,omitempty"`
		Boost      *float32    `json:"boost,omitempty"`
	}{
		PolyPoints: q.polyPoints,
		Field:      q.field,
		Boost:      q.boost,
	}

	return json.Marshal(outStruct)
}

// NewGeoPolygonQuery creates a new GeoPolygonQuery.
func NewGeoPolygonQuery(coords []Coordinate) *GeoPolygonQuery {
	var polyPoints [][]float64
	for _, coord := range coords {
		polyPoints = append(polyPoints, []float64{coord.Lon, coord.Lat})
	}
	q := &GeoPolygonQuery{polyPoints: polyPoints}
	return q
}

// Field specifies the field for this query.
func (q *GeoPolygonQuery) Field(field string) *GeoPolygonQuery {
	q.field = &field
	return q
}

// Boost specifies the boost for this query.
func (q *GeoPolygonQuery) Boost(boost float32) *GeoPolygonQuery {
	q.boost = &boost
	return q
}
