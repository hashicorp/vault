package search

import (
	"encoding/json"
)

// SearchSort represents an search sorting for a search query.
type Sort interface {
}

type searchSortBase struct {
	options map[string]interface{}
}

func newSearchSortBase() searchSortBase {
	return searchSortBase{
		options: make(map[string]interface{}),
	}
}

// MarshalJSON marshal's this query to JSON for the search REST API.
func (q searchSortBase) MarshalJSON() ([]byte, error) {
	return json.Marshal(q.options)
}

// SearchSortScore represents a search score sort.
type SearchSortScore struct {
	searchSortBase
}

// NewSearchSortScore creates a new SearchSortScore.
func NewSearchSortScore() *SearchSortScore {
	q := &SearchSortScore{newSearchSortBase()}
	q.options["by"] = "score"
	return q
}

// Descending specifies the ordering of the results.
func (q *SearchSortScore) Descending(descending bool) *SearchSortScore {
	q.options["desc"] = descending
	return q
}

// SearchSortID represents a search Document ID sort.
type SearchSortID struct {
	searchSortBase
}

// NewSearchSortID creates a new SearchSortScore.
func NewSearchSortID() *SearchSortID {
	q := &SearchSortID{newSearchSortBase()}
	q.options["by"] = "id"
	return q
}

// Descending specifies the ordering of the results.
func (q *SearchSortID) Descending(descending bool) *SearchSortID {
	q.options["desc"] = descending
	return q
}

// SearchSortField represents a search field sort.
type SearchSortField struct {
	searchSortBase
}

// NewSearchSortField creates a new SearchSortField.
func NewSearchSortField(field string) *SearchSortField {
	q := &SearchSortField{newSearchSortBase()}
	q.options["by"] = "field"
	q.options["field"] = field
	return q
}

// Type allows you to specify the search field sort type.
func (q *SearchSortField) Type(value string) *SearchSortField {
	q.options["type"] = value
	return q
}

// Mode allows you to specify the search field sort mode.
func (q *SearchSortField) Mode(mode string) *SearchSortField {
	q.options["mode"] = mode
	return q
}

// Missing allows you to specify the search field sort missing behaviour.
func (q *SearchSortField) Missing(missing string) *SearchSortField {
	q.options["missing"] = missing
	return q
}

// Descending specifies the ordering of the results.
func (q *SearchSortField) Descending(descending bool) *SearchSortField {
	q.options["desc"] = descending
	return q
}

// SearchSortGeoDistance represents a search geo sort.
type SearchSortGeoDistance struct {
	searchSortBase
}

// NewSearchSortGeoDistance creates a new SearchSortGeoDistance.
func NewSearchSortGeoDistance(field string, lon, lat float64) *SearchSortGeoDistance {
	q := &SearchSortGeoDistance{newSearchSortBase()}
	q.options["by"] = "geo_distance"
	q.options["field"] = field
	q.options["location"] = []float64{lon, lat}
	return q
}

// Unit specifies the unit used for sorting
func (q *SearchSortGeoDistance) Unit(unit string) *SearchSortGeoDistance {
	q.options["unit"] = unit
	return q
}

// Descending specifies the ordering of the results.
func (q *SearchSortGeoDistance) Descending(descending bool) *SearchSortGeoDistance {
	q.options["desc"] = descending
	return q
}
