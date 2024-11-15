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

// BitResizeFlags specifies the bitwise operation flags for resize.
type BitResizeFlags int

const (
	// BitResizeFlagsDefault specifies the defalt flag.
	BitResizeFlagsDefault BitResizeFlags = 0

	// BitResizeFlagsFromFront Adds/removes bytes from the beginning instead of the end.
	BitResizeFlagsFromFront BitResizeFlags = 1

	// BitResizeFlagsGrowOnly will only allow the []byte size to increase.
	BitResizeFlagsGrowOnly BitResizeFlags = 2

	// BitResizeFlagsShrinkOnly will only allow the []byte size to decrease.
	BitResizeFlagsShrinkOnly BitResizeFlags = 4
)
