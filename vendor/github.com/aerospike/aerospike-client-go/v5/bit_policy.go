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

// BitPolicy determines the Bit operation policy.
type BitPolicy struct {
	flags int
}

// DefaultBitPolicy will return the default BitPolicy
func DefaultBitPolicy() *BitPolicy {
	return &BitPolicy{BitWriteFlagsDefault}
}

// NewBitPolicy will return a BitPolicy will provided flags.
func NewBitPolicy(flags int) *BitPolicy {
	return &BitPolicy{flags}
}
