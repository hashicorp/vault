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

// BitWriteFlags specify bitwise operation policy write flags.

const (
	// BitWriteFlagsDefault allows create or update.
	BitWriteFlagsDefault = 0

	// BitWriteFlagsCreateOnly specifies that:
	// If the bin already exists, the operation will be denied.
	// If the bin does not exist, a new bin will be created.
	BitWriteFlagsCreateOnly = 1

	// BitWriteFlagsUpdateOnly specifies that:
	// If the bin already exists, the bin will be overwritten.
	// If the bin does not exist, the operation will be denied.
	BitWriteFlagsUpdateOnly = 2

	// BitWriteFlagsNoFail specifies not to raise error if operation is denied.
	BitWriteFlagsNoFail = 4

	// BitWriteFlagsPartial allows other valid operations to be committed if this operations is
	// denied due to flag constraints.
	BitWriteFlagsPartial = 8
)
