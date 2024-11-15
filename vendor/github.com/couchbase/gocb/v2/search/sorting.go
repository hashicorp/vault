package search

import (
	"encoding/json"
)

// SearchSort represents an search sorting for a search query.
type Sort interface {
}

// SearchSortScore represents a search score sort.
type SearchSortScore struct {
	by   string
	desc bool
}

// MarshalJSON marshal's this query to JSON for the search REST API.
func (q SearchSortScore) MarshalJSON() ([]byte, error) {

	return json.Marshal(struct {
		By   string `json:"by"`
		Desc bool   `json:"desc,omitempty"`
	}{
		By:   q.by,
		Desc: q.desc,
	})
}

// NewSearchSortScore creates a new SearchSortScore.
func NewSearchSortScore() *SearchSortScore {
	q := &SearchSortScore{by: "score"}
	return q
}

// Descending specifies the ordering of the results.
func (q *SearchSortScore) Descending(descending bool) *SearchSortScore {
	q.desc = descending
	return q
}

// SearchSortID represents a search Document ID sort.
type SearchSortID struct {
	by   string
	desc bool
}

// MarshalJSON marshal's this query to JSON for the search REST API.
func (q SearchSortID) MarshalJSON() ([]byte, error) {

	return json.Marshal(struct {
		By   string `json:"by"`
		Desc bool   `json:"desc,omitempty"`
	}{
		By:   q.by,
		Desc: q.desc,
	})
}

// NewSearchSortID creates a new SearchSortScore.
func NewSearchSortID() *SearchSortID {
	q := &SearchSortID{by: "id"}
	return q
}

// Descending specifies the ordering of the results.
func (q *SearchSortID) Descending(descending bool) *SearchSortID {
	q.desc = descending
	return q
}

// SearchSortField represents a search field sort.
type SearchSortField struct {
	by       string
	field    string
	sortType string
	mode     string
	missing  string
	desc     bool
}

// MarshalJSON marshal's this query to JSON for the search REST API.
func (q SearchSortField) MarshalJSON() ([]byte, error) {

	return json.Marshal(struct {
		By       string `json:"by"`
		Field    string `json:"field"`
		Desc     bool   `json:"desc,omitempty"`
		SortType string `json:"type,omitempty"`
		Mode     string `json:"mode,omitempty"`
		Missing  string `json:"missing,omitempty"`
	}{
		By:       q.by,
		Field:    q.field,
		Desc:     q.desc,
		SortType: q.sortType,
		Mode:     q.mode,
		Missing:  q.missing,
	})
}

// NewSearchSortField creates a new SearchSortField.
func NewSearchSortField(field string) *SearchSortField {
	q := &SearchSortField{by: "field", field: field}
	return q
}

// Type allows you to specify the search field sort type.
func (q *SearchSortField) Type(value string) *SearchSortField {
	q.sortType = value
	return q
}

// Mode allows you to specify the search field sort mode.
func (q *SearchSortField) Mode(mode string) *SearchSortField {
	q.mode = mode
	return q
}

// Missing allows you to specify the search field sort missing behaviour.
func (q *SearchSortField) Missing(missing string) *SearchSortField {
	q.missing = missing
	return q
}

// Descending specifies the ordering of the results.
func (q *SearchSortField) Descending(descending bool) *SearchSortField {
	q.desc = descending
	return q
}

// SearchSortGeoDistanceUnits represents the set of unit values available for use with SearchSortGeoDistance.
type SearchSortGeoDistanceUnits string

const (
	SearchSortGeoDistanceUnitsMeters        SearchSortGeoDistanceUnits = "meters"
	SearchSortGeoDistanceUnitsCentimeters   SearchSortGeoDistanceUnits = "centimeters"
	SearchSortGeoDistanceUnitsFeet          SearchSortGeoDistanceUnits = "feet"
	SearchSortGeoDistanceUnitsInches        SearchSortGeoDistanceUnits = "inch"
	SearchSortGeoDistanceUnitsKilometers    SearchSortGeoDistanceUnits = "kilometers"
	SearchSortGeoDistanceUnitsMiles         SearchSortGeoDistanceUnits = "miles"
	SearchSortGeoDistanceUnitsMilliMeters   SearchSortGeoDistanceUnits = "millimeters"
	SearchSortGeoDistanceUnitsNauticalMiles SearchSortGeoDistanceUnits = "nauticalmiles"
	SearchSortGeoDistanceUnitsYards         SearchSortGeoDistanceUnits = "yards"
)

// SearchSortGeoDistance represents a search geo sort.
type SearchSortGeoDistance struct {
	by       string
	field    string
	location []float64
	// See SearchSortGeoDistanceUnits for the set of values available to use with unit.
	unit string
	desc bool
}

// MarshalJSON marshal's this query to JSON for the search REST API.
func (q SearchSortGeoDistance) MarshalJSON() ([]byte, error) {

	return json.Marshal(struct {
		By       string    `json:"by"`
		Field    string    `json:"field"`
		Location []float64 `json:"location"`
		Unit     string    `json:"unit,omitempty"`
		Desc     bool      `json:"desc,omitempty"`
	}{
		By:       q.by,
		Field:    q.field,
		Location: q.location,
		Unit:     q.unit,
		Desc:     q.desc,
	})
}

// NewSearchSortGeoDistance creates a new SearchSortGeoDistance.
func NewSearchSortGeoDistance(field string, lon, lat float64) *SearchSortGeoDistance {
	q := &SearchSortGeoDistance{by: "geo_distance", field: field, location: []float64{lon, lat}}
	return q
}

// Unit specifies the unit used for sorting
func (q *SearchSortGeoDistance) Unit(unit string) *SearchSortGeoDistance {
	q.unit = unit
	return q
}

// Descending specifies the ordering of the results.
func (q *SearchSortGeoDistance) Descending(descending bool) *SearchSortGeoDistance {
	q.desc = descending
	return q
}
