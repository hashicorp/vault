package search

import (
	"encoding/json"
)

// Facet represents a facet for a search query.
type Facet interface {
}

type termFacetData struct {
	Field string `json:"field,omitempty"`
	Size  uint64 `json:"size,omitempty"`
}

// TermFacet is an search term facet.
type TermFacet struct {
	data termFacetData
}

// MarshalJSON marshal's this facet to JSON for the search REST API.
func (f TermFacet) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.data)
}

// NewTermFacet creates a new TermFacet
func NewTermFacet(field string, size uint64) *TermFacet {
	mq := &TermFacet{}
	mq.data.Field = field
	mq.data.Size = size
	return mq
}

type numericFacetRange struct {
	Name  string  `json:"name,omitempty"`
	Start float64 `json:"min,omitempty"`
	End   float64 `json:"max,omitempty"`
}
type numericFacetData struct {
	Field         string              `json:"field,omitempty"`
	Size          uint64              `json:"size,omitempty"`
	NumericRanges []numericFacetRange `json:"numeric_ranges,omitempty"`
}

// NumericFacet is an search numeric range facet.
type NumericFacet struct {
	data numericFacetData
}

// MarshalJSON marshal's this facet to JSON for the search REST API.
func (f NumericFacet) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.data)
}

// AddRange adds a new range to this numeric range facet.
func (f *NumericFacet) AddRange(name string, start, end float64) *NumericFacet {
	f.data.NumericRanges = append(f.data.NumericRanges, numericFacetRange{
		Name:  name,
		Start: start,
		End:   end,
	})
	return f
}

// NewNumericFacet creates a new numeric range facet.
func NewNumericFacet(field string, size uint64) *NumericFacet {
	mq := &NumericFacet{}
	mq.data.Field = field
	mq.data.Size = size
	return mq
}

type dateFacetRange struct {
	Name  string `json:"name,omitempty"`
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}
type dateFacetData struct {
	Field      string           `json:"field,omitempty"`
	Size       uint64           `json:"size,omitempty"`
	DateRanges []dateFacetRange `json:"date_ranges,omitempty"`
}

// DateFacet is an search date range facet.
type DateFacet struct {
	data dateFacetData
}

// MarshalJSON marshal's this facet to JSON for the search REST API.
func (f DateFacet) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.data)
}

// AddRange adds a new range to this date range facet.
func (f *DateFacet) AddRange(name string, start, end string) *DateFacet {
	f.data.DateRanges = append(f.data.DateRanges, dateFacetRange{
		Name:  name,
		Start: start,
		End:   end,
	})
	return f
}

// NewDateFacet creates a new date range facet.
func NewDateFacet(field string, size uint64) *DateFacet {
	mq := &DateFacet{}
	mq.data.Field = field
	mq.data.Size = size
	return mq
}
