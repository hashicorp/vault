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

package particleType

// Server particle types. Unsupported types are commented out.
const (
	//revive:disable
	NULL    = 0
	INTEGER = 1
	FLOAT   = 2
	STRING  = 3
	BLOB    = 4
	DIGEST  = 6
	BOOL    = 17
	HLL     = 18
	MAP     = 19
	LIST    = 20
	LDT     = 21
	GEOJSON = 23
	//revive:enable
)
