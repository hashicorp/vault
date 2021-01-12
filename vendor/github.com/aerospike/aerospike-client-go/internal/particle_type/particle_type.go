// Copyright 2013-2020 Aerospike, Inc.
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

const (
	// Server particle types. Unsupported types are commented out.
	NULL    = 0
	INTEGER = 1
	FLOAT   = 2
	STRING  = 3
	BLOB    = 4
	// TIMESTAMP       = 5
	DIGEST = 6
	// JBLOB  = 7
	// CSHARP_BLOB     = 8
	// PYTHON_BLOB     = 9
	// RUBY_BLOB       = 10
	// PHP_BLOB        = 11
	// ERLANG_BLOB     = 12
	// SEGMENT_POINTER = 13
	// RTA_LIST        = 14
	// RTA_DICT        = 15
	// RTA_APPEND_DICT = 16
	// RTA_APPEND_LIST = 17
	HLL     = 18
	MAP     = 19
	LIST    = 20
	LDT     = 21
	GEOJSON = 23
)
