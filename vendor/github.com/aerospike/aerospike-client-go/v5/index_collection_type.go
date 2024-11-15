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

import "fmt"

// IndexCollectionType is the secondary index collection type.
type IndexCollectionType int

const (

	// ICT_DEFAULT is the Normal scalar index.
	ICT_DEFAULT IndexCollectionType = iota

	// ICT_LIST is Index list elements.
	ICT_LIST

	// ICT_MAPKEYS is Index map keys.
	ICT_MAPKEYS

	// ICT_MAPVALUES is Index map values.
	ICT_MAPVALUES
)

// ictToString converts IndexCollectionType to string representations
func ictToString(ict IndexCollectionType) string {
	switch ict {

	case ICT_LIST:
		return "LIST"

	case ICT_MAPKEYS:
		return "MAPKEYS"

	case ICT_MAPVALUES:
		return "MAPVALUES"

	default:
		panic(fmt.Sprintf("Unknown IndexCollectionType value %v", ict))
	}
}
