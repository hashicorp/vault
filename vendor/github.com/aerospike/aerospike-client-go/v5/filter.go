// Copyright 2014-2021 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

import (
	"fmt"

	ParticleType "github.com/aerospike/aerospike-client-go/v5/internal/particle_type"
)

// Filter specifies a query filter definition.
type Filter struct {
	name              string
	idxType           IndexCollectionType
	valueParticleType int
	begin             Value
	end               Value
}

// NewEqualFilter creates a new equality filter instance for query.
func NewEqualFilter(binName string, value interface{}) *Filter {
	val := NewValue(value)
	return newFilter(binName, ICT_DEFAULT, val.GetType(), val, val)
}

// NewRangeFilter creates a range filter for query.
// Range arguments must be int64 values.
// String ranges are not supported.
func NewRangeFilter(binName string, begin int64, end int64) *Filter {
	vBegin, vEnd := NewValue(begin), NewValue(end)
	return newFilter(binName, ICT_DEFAULT, vBegin.GetType(), vBegin, vEnd)
}

// NewContainsFilter creates a contains filter for query on collection index.
func NewContainsFilter(binName string, indexCollectionType IndexCollectionType, value interface{}) *Filter {
	v := NewValue(value)
	return newFilter(binName, indexCollectionType, v.GetType(), v, v)
}

// NewContainsRangeFilter creates a contains filter for query on ranges of data in a collection index.
func NewContainsRangeFilter(binName string, indexCollectionType IndexCollectionType, begin, end int64) *Filter {
	vBegin, vEnd := NewValue(begin), NewValue(end)
	return newFilter(binName, indexCollectionType, vBegin.GetType(), vBegin, vEnd)
}

// NewGeoWithinRegionFilter creates a geospatial "within region" filter for query.
// Argument must be a valid GeoJSON region.
func NewGeoWithinRegionFilter(binName, region string) *Filter {
	v := NewStringValue(region)
	return newFilter(binName, ICT_DEFAULT, ParticleType.GEOJSON, v, v)
}

// NewGeoWithinRegionForCollectionFilter creates a geospatial "within region" filter for query on collection index.
// Argument must be a valid GeoJSON region.
func NewGeoWithinRegionForCollectionFilter(binName string, collectionType IndexCollectionType, region string) *Filter {
	v := NewStringValue(region)
	return newFilter(binName, collectionType, ParticleType.GEOJSON, v, v)
}

// NewGeoRegionsContainingPointFilter creates a geospatial "containing point" filter for query.
// Argument must be a valid GeoJSON point.
func NewGeoRegionsContainingPointFilter(binName, point string) *Filter {
	v := NewStringValue(point)
	return newFilter(binName, ICT_DEFAULT, ParticleType.GEOJSON, v, v)
}

// NewGeoRegionsContainingPointForCollectionFilter creates a geospatial "containing point" filter for query on collection index.
// Argument must be a valid GeoJSON point.
func NewGeoRegionsContainingPointForCollectionFilter(binName string, collectionType IndexCollectionType, point string) *Filter {
	v := NewStringValue(point)
	return newFilter(binName, collectionType, ParticleType.GEOJSON, v, v)
}

// NewGeoWithinRadiusFilter creates a geospatial "within radius" filter for query.
// Arguments must be valid longitude/latitude/radius (meters) values.
func NewGeoWithinRadiusFilter(binName string, lng, lat, radius float64) *Filter {
	rgnStr := fmt.Sprintf("{ \"type\": \"AeroCircle\", "+"\"coordinates\": [[%.8f, %.8f], %f] }", lng, lat, radius)
	return newFilter(binName, ICT_DEFAULT, ParticleType.GEOJSON, NewValue(rgnStr), NewValue(rgnStr))
}

// NewGeoWithinRadiusForCollectionFilter creates a geospatial "within radius" filter for query on collection index.
// Arguments must be valid longitude/latitude/radius (meters) values.
func NewGeoWithinRadiusForCollectionFilter(binName string, collectionType IndexCollectionType, lng, lat, radius float64) *Filter {
	rgnStr := fmt.Sprintf("{ \"type\": \"AeroCircle\", "+"\"coordinates\": [[%.8f, %.8f], %f] }", lng, lat, radius)
	return newFilter(binName, collectionType, ParticleType.GEOJSON, NewValue(rgnStr), NewValue(rgnStr))
}

// Create a filter for query.
// Range arguments must be longs or integers which can be cast to longs.
// String ranges are not supported.
func newFilter(name string, indexCollectionType IndexCollectionType, valueParticleType int, begin Value, end Value) *Filter {
	return &Filter{
		name:              name,
		idxType:           indexCollectionType,
		valueParticleType: valueParticleType,
		begin:             begin,
		end:               end,
	}
}

// IndexCollectionType returns filter's index type.
func (fltr *Filter) IndexCollectionType() IndexCollectionType {
	return fltr.idxType
}

// EstimateSize will estimate the size of the filter for wire protocol
func (fltr *Filter) EstimateSize() (int, Error) {
	// bin name size(1) + particle type size(1) + begin particle size(4) + end particle size(4) = 10
	szBegin, err := fltr.begin.EstimateSize()
	if err != nil {
		return szBegin, err
	}

	szEnd, err := fltr.end.EstimateSize()
	if err != nil {
		return szEnd, err
	}

	return len(fltr.name) + szBegin + szEnd + 10, nil
}

func (fltr *Filter) write(cmd *baseCommand) (int, Error) {
	size := 0

	// Write name length
	cmd.WriteByte(byte(len(fltr.name)))
	size++

	// Write Name
	n, err := cmd.WriteString(fltr.name)
	if err != nil {
		return size + n, err
	}
	size += n

	// Write particle type.
	cmd.WriteByte(byte(fltr.valueParticleType))
	size++

	// Write filter begin.
	esz, err := fltr.begin.EstimateSize()
	if err != nil {
		return size, err
	}

	n = cmd.WriteInt32(int32(esz))
	size += n

	n, err = fltr.begin.write(cmd)
	if err != nil {
		return size + n, err
	}
	size += n

	// Write filter end.
	esz, err = fltr.end.EstimateSize()
	if err != nil {
		return size, err
	}

	n = cmd.WriteInt32(int32(esz))
	size += n

	n, err = fltr.end.write(cmd)
	if err != nil {
		return size + n, err
	}
	size += n

	return size, nil
}
