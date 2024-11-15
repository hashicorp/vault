// Copyright 2014-2019 Aerospike, Inc.
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

// HLLWriteFlags specifies the HLL write operation flags.

const (
	// HLLWriteFlagsDefault is Default. Allow create or update.
	HLLWriteFlagsDefault = 0

	// HLLWriteFlagsCreateOnly behaves like the following:
	// If the bin already exists, the operation will be denied.
	// If the bin does not exist, a new bin will be created.
	HLLWriteFlagsCreateOnly = 1

	// HLLWriteFlagsUpdateOnly behaves like the following:
	// If the bin already exists, the bin will be overwritten.
	// If the bin does not exist, the operation will be denied.
	HLLWriteFlagsUpdateOnly = 2

	// HLLWriteFlagsNoFail does not raise error if operation is denied.
	HLLWriteFlagsNoFail = 4

	// HLLWriteFlagsAllowFold allows the resulting set to be the minimum of provided index bits.
	// Also, allow the usage of less precise HLL algorithms when minHash bits
	// of all participating sets do not match.
	HLLWriteFlagsAllowFold = 8
)
