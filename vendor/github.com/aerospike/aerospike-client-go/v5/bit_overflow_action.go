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

// BitOverflowAction specifies the action to take when bitwise add/subtract results in overflow/underflow.
type BitOverflowAction int

const (
	// BitOverflowActionFail specifies to fail operation with error.
	BitOverflowActionFail BitOverflowAction = 0

	// BitOverflowActionSaturate specifies that in add/subtract overflows/underflows, set to max/min value.
	// Example: MAXINT + 1 = MAXINT
	BitOverflowActionSaturate BitOverflowAction = 2

	// BitOverflowActionWrap specifies that in add/subtract overflows/underflows, wrap the value.
	// Example: MAXINT + 1 = -1
	BitOverflowActionWrap BitOverflowAction = 4
)
