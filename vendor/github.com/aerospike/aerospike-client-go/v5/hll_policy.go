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

// HLLPolicy determines the HyperLogLog operation policy.
type HLLPolicy struct {
	flags int
}

// DefaultHLLPolicy uses the default policy when performing HLL operations.
func DefaultHLLPolicy() *HLLPolicy {
	return &HLLPolicy{HLLWriteFlagsDefault}
}

// NewHLLPolicy uses specified HLLWriteFlags when performing HLL operations.
func NewHLLPolicy(flags int) *HLLPolicy {
	return &HLLPolicy{flags}
}
