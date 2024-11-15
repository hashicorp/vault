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

// IndexType the type of the secondary index.
type IndexType string

const (
	// NUMERIC specifies an index on numeric values.
	NUMERIC IndexType = "NUMERIC"

	// STRING specifies an index on string values.
	STRING IndexType = "STRING"

	// GEO2DSPHERE specifies 2-dimensional spherical geospatial index.
	GEO2DSPHERE IndexType = "GEO2DSPHERE"
)
