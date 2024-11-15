package search

import (
	"fmt"
	"math"

	"github.com/couchbase/goprotostellar/genproto/search_v1"
)

// Internal is used for internal functionality.
// Internal: This should never be used and is not supported.
type Internal struct {
}

func (i Internal) MapFacetsToPs(facets map[string]Facet) (map[string]*search_v1.Facet, error) {
	out := make(map[string]*search_v1.Facet)
	var err error
	for name, facet := range facets {
		if out[name], err = i.mapFacetToPS(facet); err != nil {
			return nil, err
		}
	}

	return out, nil
}

// converts uint64 to uint32, checking
// value won't be truncated by uint32.
func convertUint64ToUnit32(in uint64) (uint32, error) {
	if in > uint64(math.MaxUint32) {
		return 0, fmt.Errorf("value: %d overflows uint32", in)
	}

	return uint32(in), nil
}

func (i Internal) mapFacetToPS(facet Facet) (*search_v1.Facet, error) {
	psFacet := search_v1.Facet{}
	switch f := facet.(type) {
	case *TermFacet:
		size, err := convertUint64ToUnit32(f.data.Size)
		if err != nil {
			return nil, err
		}

		psFacet.Facet = &search_v1.Facet_TermFacet{
			TermFacet: &search_v1.TermFacet{
				Field: f.data.Field,
				Size:  size,
			},
		}
	case *DateFacet:
		size, err := convertUint64ToUnit32(f.data.Size)
		if err != nil {
			return nil, err
		}
		psFacet.Facet = &search_v1.Facet_DateRangeFacet{
			DateRangeFacet: &search_v1.DateRangeFacet{
				Field:      f.data.Field,
				Size:       size,
				DateRanges: i.mapDateRangeFacetToPs(f.data.DateRanges),
			},
		}
	case *NumericFacet:
		size, err := convertUint64ToUnit32(f.data.Size)
		if err != nil {
			return nil, err
		}

		numericRanges, err := i.mapNumericRangeFacetToPs(f.data.NumericRanges)
		if err != nil {
			return nil, err
		}

		psFacet.Facet = &search_v1.Facet_NumericRangeFacet{
			NumericRangeFacet: &search_v1.NumericRangeFacet{
				Field:         f.data.Field,
				Size:          size,
				NumericRanges: numericRanges,
			},
		}
	default:
		return nil, fmt.Errorf("invalid facet option specified")
	}

	return &psFacet, nil
}

func (i Internal) mapNumericRangeFacetToPs(numericRanges []numericFacetRange) ([]*search_v1.NumericRange, error) {
	out := make([]*search_v1.NumericRange, len(numericRanges))
	for i, numericRange := range numericRanges {
		min := float32(numericRange.Start) // TODO: float64 -> float32
		max := float32(numericRange.End)
		out[i] = &search_v1.NumericRange{
			Name: numericRange.Name,
			Min:  &min,
			Max:  &max,
		}
	}
	return out, nil
}

func (i Internal) mapDateRangeFacetToPs(dateRanges []dateFacetRange) []*search_v1.DateRange {
	out := make([]*search_v1.DateRange, len(dateRanges))

	for i := range dateRanges {
		out[i] = &search_v1.DateRange{
			Name:  dateRanges[i].Name,
			Start: &dateRanges[i].Start,
			End:   &dateRanges[i].End,
		}
	}

	return out
}

// helper function for handling conjunct/disjunct which consist of multiple queries.
func (i Internal) mapQueriesToPs(queries []Query) ([]*search_v1.Query, error) {
	out := make([]*search_v1.Query, len(queries))
	var err error
	for index, query := range queries {
		if out[index], err = i.MapQueryToPs(query); err != nil {
			return nil, err
		}
	}

	return out, nil
}

// convert query into the Protostellar format.
func (i Internal) MapQueryToPs(query Query) (*search_v1.Query, error) {
	switch q := query.(type) {
	case *BooleanFieldQuery:
		return &search_v1.Query{
			Query: &search_v1.Query_BooleanFieldQuery{
				BooleanFieldQuery: &search_v1.BooleanFieldQuery{
					Boost: q.boost,
					Field: q.field,
					Value: q.val,
				},
			},
		}, nil
	case *BooleanQuery:
		booleanQuery := &search_v1.BooleanQuery{
			Boost: &q.data.Boost,
		}
		if q.data.Must != nil {
			mustQueries, err := i.mapQueriesToPs(q.data.Must.conjuncts)
			if err != nil {
				return nil, err
			}

			booleanQuery.Must = &search_v1.ConjunctionQuery{
				Queries: mustQueries,
				Boost:   q.data.Must.boost,
			}
		}

		if q.data.MustNot != nil {
			mustNotQueries, err := i.mapQueriesToPs(q.data.MustNot.disjuncts)
			if err != nil {
				return nil, err
			}

			booleanQuery.MustNot = &search_v1.DisjunctionQuery{
				Queries: mustNotQueries,
				Boost:   q.data.MustNot.boost,
				Minimum: q.data.MustNot.min,
			}
		}

		if q.data.Should != nil {
			shouldQueries, err := i.mapQueriesToPs(q.data.Should.disjuncts)
			if err != nil {
				return nil, err
			}

			booleanQuery.Should = &search_v1.DisjunctionQuery{
				Queries: shouldQueries,
				Boost:   q.data.Should.boost,
				Minimum: q.data.Should.min,
			}
		}
		return &search_v1.Query{Query: &search_v1.Query_BooleanQuery{
			BooleanQuery: booleanQuery,
		}}, nil

	case *ConjunctionQuery:
		queries, err := i.mapQueriesToPs(q.conjuncts)
		if err != nil {
			return nil, err
		}
		return &search_v1.Query{Query: &search_v1.Query_ConjunctionQuery{
			ConjunctionQuery: &search_v1.ConjunctionQuery{
				Boost:   q.boost,
				Queries: queries,
			},
		}}, nil
	case *DateRangeQuery:
		return &search_v1.Query{Query: &search_v1.Query_DateRangeQuery{
			DateRangeQuery: &search_v1.DateRangeQuery{ // TODO: inclusive bool is missing.
				Boost:          q.boost,
				Field:          q.field,
				DateTimeParser: q.dateTimeParser,
				StartDate:      q.start,
				EndDate:        q.end,
			},
		},
		}, nil
	case *DisjunctionQuery:
		queries, err := i.mapQueriesToPs(q.disjuncts)
		if err != nil {
			return nil, err
		}
		return &search_v1.Query{Query: &search_v1.Query_DisjunctionQuery{
			DisjunctionQuery: &search_v1.DisjunctionQuery{
				Boost:   q.boost,
				Queries: queries,
				Minimum: q.min,
			},
		},
		}, nil
	case *DocIDQuery:
		return &search_v1.Query{Query: &search_v1.Query_DocIdQuery{
			DocIdQuery: &search_v1.DocIdQuery{
				Boost: q.boost,
				Ids:   q.ids,
			},
		},
		}, nil
	case *GeoBoundingBoxQuery:
		return &search_v1.Query{Query: &search_v1.Query_GeoBoundingBoxQuery{
			GeoBoundingBoxQuery: &search_v1.GeoBoundingBoxQuery{
				Boost: q.boost,
				Field: q.field,
				BottomRight: &search_v1.LatLng{
					Longitude: q.bottomRight[0],
					Latitude:  q.bottomRight[1],
				},
				TopLeft: &search_v1.LatLng{
					Longitude: q.topLeft[0],
					Latitude:  q.topLeft[1],
				},
			},
		},
		}, nil
	case *GeoDistanceQuery:
		return &search_v1.Query{Query: &search_v1.Query_GeoDistanceQuery{
			GeoDistanceQuery: &search_v1.GeoDistanceQuery{
				Boost:    q.boost,
				Field:    q.field,
				Distance: q.distance,
				Center: &search_v1.LatLng{
					Longitude: q.location[0],
					Latitude:  q.location[1],
				},
			},
		},
		}, nil
	case *GeoPolygonQuery:
		vertices := make([]*search_v1.LatLng, len(q.polyPoints))
		for i, lonLat := range q.polyPoints {
			vertices[i] = &search_v1.LatLng{
				Longitude: lonLat[0],
				Latitude:  lonLat[1],
			}
		}

		return &search_v1.Query{Query: &search_v1.Query_GeoPolygonQuery{
			GeoPolygonQuery: &search_v1.GeoPolygonQuery{
				Boost:    q.boost,
				Field:    q.field,
				Vertices: vertices,
			},
		},
		}, nil
	case *MatchAllQuery:
		return &search_v1.Query{Query: &search_v1.Query_MatchAllQuery{
			MatchAllQuery: &search_v1.MatchAllQuery{},
		},
		}, nil
	case *MatchNoneQuery:
		return &search_v1.Query{Query: &search_v1.Query_MatchNoneQuery{
			MatchNoneQuery: &search_v1.MatchNoneQuery{},
		},
		}, nil
	case *MatchPhraseQuery:
		return &search_v1.Query{Query: &search_v1.Query_MatchPhraseQuery{
			MatchPhraseQuery: &search_v1.MatchPhraseQuery{
				Boost:    q.boost,
				Analyzer: q.analyzer,
				Field:    q.field,
				Phrase:   q.matchPhrase,
			},
		},
		}, nil
	case *MatchQuery:
		var operator *search_v1.MatchQuery_Operator
		if q.operator != nil {
			switch *q.operator {
			case MatchOperatorAnd:
				operatorAnd := search_v1.MatchQuery_OPERATOR_AND
				operator = &operatorAnd
			case MatchOperatorOr:
				operatorOr := search_v1.MatchQuery_OPERATOR_OR
				operator = &operatorOr
			}
		}

		return &search_v1.Query{Query: &search_v1.Query_MatchQuery{
			MatchQuery: &search_v1.MatchQuery{
				Boost:        q.boost,
				Field:        q.field,
				Value:        q.match,
				Fuzziness:    q.fuzziness,
				Analyzer:     q.analyzer,
				Operator:     operator,
				PrefixLength: q.prefixLength,
			},
		},
		}, nil
	case *NumericRangeQuery:
		return &search_v1.Query{Query: &search_v1.Query_NumericRangeQuery{
			NumericRangeQuery: &search_v1.NumericRangeQuery{
				Boost:        q.boost,
				Field:        q.field,
				Min:          q.min,
				InclusiveMin: q.inclusiveMin,
				Max:          q.max,
				InclusiveMax: q.inclusiveMax,
			},
		}}, nil
	case *PhraseQuery:
		return &search_v1.Query{Query: &search_v1.Query_PhraseQuery{
			PhraseQuery: &search_v1.PhraseQuery{
				Boost: q.boost,
				Field: q.field,
				Terms: q.terms,
			},
		}}, nil
	case *PrefixQuery:
		return &search_v1.Query{Query: &search_v1.Query_PrefixQuery{
			PrefixQuery: &search_v1.PrefixQuery{
				Prefix: q.prefix,
				Field:  q.field,
				Boost:  q.boost,
			},
		}}, nil
	case *QueryStringQuery:
		return &search_v1.Query{Query: &search_v1.Query_QueryStringQuery{
			QueryStringQuery: &search_v1.QueryStringQuery{
				QueryString: q.query,
				Boost:       q.boost,
			},
		}}, nil
	case *RegexpQuery:
		return &search_v1.Query{Query: &search_v1.Query_RegexpQuery{
			RegexpQuery: &search_v1.RegexpQuery{
				Regexp: q.regexp,
				Field:  q.field,
				Boost:  q.boost,
			},
		}}, nil
	case *TermQuery:
		return &search_v1.Query{Query: &search_v1.Query_TermQuery{
			TermQuery: &search_v1.TermQuery{
				Term:         q.term,
				Field:        q.field,
				Fuzziness:    q.fuzziness,
				Boost:        q.boost,
				PrefixLength: q.prefixLength,
			},
		}}, nil
	case *TermRangeQuery:
		return &search_v1.Query{Query: &search_v1.Query_TermRangeQuery{
			TermRangeQuery: &search_v1.TermRangeQuery{
				Boost:        q.boost,
				Field:        q.field,
				Max:          q.max,
				InclusiveMax: q.inclusiveMax,
				Min:          q.min,
				InclusiveMin: q.inclusiveMax,
			},
		}}, nil
	case *WildcardQuery:
		return &search_v1.Query{Query: &search_v1.Query_WildcardQuery{
			WildcardQuery: &search_v1.WildcardQuery{
				Boost:    q.boost,
				Field:    q.field,
				Wildcard: q.wildcard,
			},
		}}, nil
	default:
		return nil, fmt.Errorf("invalid query option specified")
	}

}

func (i Internal) MapSortToPs(in []Sort) ([]*search_v1.Sorting, error) {
	out := make([]*search_v1.Sorting, len(in))

	for index, sorting := range in {
		switch s := sorting.(type) {
		case *SearchSortID:
			out[index] = &search_v1.Sorting{
				Sorting: &search_v1.Sorting_IdSorting{
					IdSorting: &search_v1.IdSorting{
						Descending: s.desc,
					},
				},
			}
		case *SearchSortField:
			out[index] = &search_v1.Sorting{
				Sorting: &search_v1.Sorting_FieldSorting{
					FieldSorting: &search_v1.FieldSorting{
						Field:      s.field,
						Descending: s.desc,
						Missing:    s.missing,
						Mode:       s.mode,
						Type:       s.sortType,
					},
				},
			}
		case *SearchSortScore:
			out[index] = &search_v1.Sorting{
				Sorting: &search_v1.Sorting_ScoreSorting{
					ScoreSorting: &search_v1.ScoreSorting{
						Descending: s.desc,
					},
				},
			}
		case *SearchSortGeoDistance:
			out[index] = &search_v1.Sorting{
				Sorting: &search_v1.Sorting_GeoDistanceSorting{
					GeoDistanceSorting: &search_v1.GeoDistanceSorting{
						Field:      s.field,
						Descending: s.desc,
						Center: &search_v1.LatLng{
							Longitude: s.location[0],
							Latitude:  s.location[1],
						},
						Unit: s.unit,
					},
				},
			}
		default:
			return nil, fmt.Errorf("invalid sort option specified")
		}
	}

	return out, nil
}
